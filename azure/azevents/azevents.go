// Package azevents handles incoming webhook events from Azure Event Grid. For
// more info on the verification/validation logic, see
// https://learn.microsoft.com/en-us/azure/event-grid/webhook-event-delivery#validation-details
package azevents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/task"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Config struct {
	Logger *zap.Logger
	DB     DB
	Now    func() time.Time

	Subscription  string
	ResourceGroup string

	// AllowedAuthSecrets is a list of shared random secrets that will be accepted
	// for incoming webhooks from the event subscription to this receiver, to
	// prevent random unauthenticated internet requests from triggering webhooks.
	AllowedAuthSecrets []string

	ParsedPortfolioTopicName string
}

type DB interface {
	Transactional(context.Context, func(tx db.Tx) error) error

	CreateBlob(tx db.Tx, b *pacta.Blob) (pacta.BlobID, error)
	CreatePortfolio(tx db.Tx, p *pacta.Portfolio) (pacta.PortfolioID, error)

	IncompleteUploads(tx db.Tx, ids []pacta.IncompleteUploadID) (map[pacta.IncompleteUploadID]*pacta.IncompleteUpload, error)
	UpdateIncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID, mutations ...db.UpdateIncompleteUploadFn) error
}

const parsedPortfolioPath = "/events/parsed_portfolio"

func (c *Config) validate() error {
	if c.Logger == nil {
		return errors.New("no logger was given")
	}
	if c.Subscription == "" {
		return errors.New("no subscription given")
	}
	if c.ResourceGroup == "" {
		return errors.New("no resource group given")
	}
	if len(c.AllowedAuthSecrets) == 0 {
		return errors.New("no auth secrets were given")
	}
	if c.ParsedPortfolioTopicName == "" {
		return errors.New("no parsed portfolio topic name given")
	}
	if c.DB == nil {
		return errors.New("no DB was given")
	}
	if c.Now == nil {
		return errors.New("no Now function was given")
	}
	return nil
}

// Server handles both validating the Event Grid subscription and handling incoming events.
type Server struct {
	logger *zap.Logger

	allowedAuthSecrets []string

	subscription  string
	resourceGroup string
	pathToTopic   map[string]string
	db            DB
	now           func() time.Time
}

func NewServer(cfg *Config) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config given: %w", err)
	}

	return &Server{
		logger:             cfg.Logger,
		allowedAuthSecrets: cfg.AllowedAuthSecrets,
		subscription:       cfg.Subscription,
		resourceGroup:      cfg.ResourceGroup,
		db:                 cfg.DB,
		now:                cfg.Now,
		pathToTopic: map[string]string{
			parsedPortfolioPath: cfg.ParsedPortfolioTopicName,
		},
	}, nil
}

func (s *Server) findValidAuthSecret(auth string) (int, bool) {
	for idx, allowed := range s.allowedAuthSecrets {
		if auth == allowed {
			return idx, true
		}
	}
	return 0, false
}

func (s *Server) verifyWebhook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		topic, ok := s.pathToTopic[r.URL.Path]
		if !ok {
			s.logger.Error("no topic found for path", zap.String("path", r.URL.Path))
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if r.Header.Get("aeg-event-type") != "SubscriptionValidation" {
			// Azure doesn't have a native way to validate where a webhook request
			// came from, so we implement a shared secret approach as described here:
			// https://learn.microsoft.com/en-us/azure/event-grid/security-authentication#using-client-secret-as-a-query-parameter
			auth := r.Header.Get("authorization")
			if auth == "" {
				s.logger.Error("webhook request was missing 'authorization' header")
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			validAuthIdx, valid := s.findValidAuthSecret(auth)
			if !valid {
				s.logger.Error("webhook request had invalid 'authorization' header", zap.String("invalid_auth", auth))
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			s.logger.Info("received webhook request with valid 'authorization' header", zap.Int("auth_secret_index", validAuthIdx), zap.String("path", r.URL.Path))
			next.ServeHTTP(w, r)
			return
		}

		// If we're here, we're validating to Azure that we own this endpoint and want
		// to accept webhook calls here. Only calls to valid webhook endpoints will
		// trigger this middleware (the rest will just get 404s), so if we've made it
		// this far, validate that yes, we'll take webhook invocations.
		var reqs []struct {
			Id              string    `json:"id"`
			Topic           string    `json:"topic"`
			Subject         string    `json:"subject"`
			EventType       string    `json:"eventType"`
			EventTime       time.Time `json:"eventTime"`
			MetadataVersion string    `json:"metadataVersion"`
			DataVersion     string    `json:"dataVersion"`
			Data            *struct {
				ValidationCode string `json:"validationCode"`
				ValidationUrl  string `json:"validationUrl"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
			s.logger.Error("failed to decode subscription validation request", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if len(reqs) != 1 {
			s.logger.Error("unexpected number of validation requests", zap.Any("reqs", reqs))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		req := reqs[0]
		if req.Data == nil {
			s.logger.Error("no data provided in validation request", zap.Any("req", req))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s.logger.Info("received SubscriptionValidation request", zap.Any("req", req))

		// Validate the request event type and topic
		if got, want := req.EventType, "Microsoft.EventGrid.SubscriptionValidationEvent"; got != want {
			s.logger.Error("invalid topic given for path", zap.String("got_event_type", got), zap.String("expected_event_type", want))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fullTopic := path.Join("/subscriptions", s.subscription, "resourceGroups", s.resourceGroup, "providers/Microsoft.EventGrid/topics", topic)
		// We lowercase these because *sometimes* the request comes from Azure with
		// "microsoft.eventgrid" instead of "Microsoft.EventGrid". This is exceptionally
		// annoying.
		if strings.ToLower(req.Topic) != strings.ToLower(fullTopic) {
			s.logger.Error("invalid topic given for path", zap.String("got_topic", req.Topic), zap.String("expected_topic", fullTopic))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s.logger.Info("validated SubscriptionValidation, responding success", zap.String("request_id", req.Id))

		resp := struct {
			ValidationResponse string `json:"validationResponse"`
		}{req.Data.ValidationCode}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.logger.Error("failed to encode JSON validation response", zap.Error(err))
		}
	})
}

func (s *Server) RegisterHandlers(r chi.Router) {
	r.Use(s.verifyWebhook)
	r.Post(parsedPortfolioPath, func(w http.ResponseWriter, r *http.Request) {
		var reqs []struct {
			Data            *task.ParsePortfolioResponse `json:"data"`
			EventType       string                       `json:"eventType"`
			ID              string                       `json:"id"`
			Subject         string                       `json:"subject"`
			DataVersion     string                       `json:"dataVersion"`
			MetadataVersion string                       `json:"metadataVersion"`
			EventTime       time.Time                    `json:"eventTime"`
			Topic           string                       `json:"topic"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
			s.logger.Error("failed to parse webhook request body", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if len(reqs) != 1 {
			s.logger.Error("webhook response had unexpected number of events", zap.Int("event_count", len(reqs)))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		req := reqs[0]

		if req.Data == nil {
			s.logger.Error("webhook response had no payload", zap.String("event_grid_id", req.ID))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if len(req.Data.Outputs) == 0 {
			s.logger.Error("webhook response had no processed portfolios", zap.String("event_grid_id", req.ID))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		portfolioIDs := []pacta.PortfolioID{}
		var ranAt time.Time
		now := s.now()
		err := s.db.Transactional(r.Context(), func(tx db.Tx) error {
			incompleteUploads, err := s.db.IncompleteUploads(tx, req.Data.Request.IncompleteUploadIDs)
			if err != nil {
				return fmt.Errorf("reading incomplete uploads: %w", err)
			}
			if len(incompleteUploads) == 0 {
				return fmt.Errorf("no incomplete uploads found for ids: %v", req.Data.Request.IncompleteUploadIDs)
			}
			var holdingsDate *pacta.HoldingsDate
			var ownerID pacta.OwnerID
			for _, iu := range incompleteUploads {
				if ownerID == "" {
					ownerID = iu.Owner.ID
				} else if ownerID != iu.Owner.ID {
					return fmt.Errorf("multiple owners found for incomplete uploads: %+v", incompleteUploads)
				}
				if iu.HoldingsDate == nil {
					return fmt.Errorf("incomplete upload %s had no holdings date", iu.ID)
				}
				if holdingsDate == nil {
					holdingsDate = iu.HoldingsDate
				} else if *holdingsDate != *iu.HoldingsDate {
					return fmt.Errorf("multiple holdings dates found for incomplete uploads: %+v", incompleteUploads)
				}
				ranAt = iu.RanAt
			}
			for i, output := range req.Data.Outputs {
				blobID, err := s.db.CreateBlob(tx, &output.Blob)
				if err != nil {
					return fmt.Errorf("creating blob %d: %w", i, err)
				}
				portfolioID, err := s.db.CreatePortfolio(tx, &pacta.Portfolio{
					Owner:        &pacta.Owner{ID: ownerID},
					Name:         output.Blob.FileName,
					NumberOfRows: output.LineCount,
					Blob:         &pacta.Blob{ID: blobID},
					HoldingsDate: holdingsDate,
				})
				if err != nil {
					return fmt.Errorf("creating portfolio %d: %w", i, err)
				}
				portfolioIDs = append(portfolioIDs, portfolioID)
			}
			for i, iu := range incompleteUploads {
				err := s.db.UpdateIncompleteUpload(
					tx,
					iu.ID,
					db.SetIncompleteUploadCompletedAt(now),
					db.SetIncompleteUploadFailureMessage(""),
					db.SetIncompleteUploadFailureCode(""))
				if err != nil {
					return fmt.Errorf("updating incomplete upload %d: %w", i, err)
				}
			}
			return nil
		})
		if err != nil {
			s.logger.Error("failed to save response to database", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		s.logger.Info("parsed portfolio",
			zap.String("task_id", string(req.Data.TaskID)),
			zap.Duration("run_time", now.Sub(ranAt)),
			zap.Strings("incomplete_upload_ids", asStrs(req.Data.Request.IncompleteUploadIDs)),
			zap.Int("incomplete_upload_count", len(req.Data.Request.IncompleteUploadIDs)),
			zap.Strings("portfolio_ids", asStrs(portfolioIDs)),
			zap.Int("portfolio_count", len(portfolioIDs)))
	})
}

func asStrs[T ~string](ts []T) []string {
	ss := make([]string, len(ts))
	for i, t := range ts {
		ss[i] = string(t)
	}
	return ss
}
