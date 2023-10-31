// Package azevents handles incoming webhook events from Azure Event Grid. For
// more info on the verification/validation logic, see
// https://learn.microsoft.com/en-us/azure/event-grid/webhook-event-delivery#validation-details
package azevents

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Config struct {
	Logger *zap.Logger

	Subscription  string
	ResourceGroup string
}

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
	return nil
}

// Server handles both validating the Event Grid subscription and handling incoming events.
type Server struct {
	logger *zap.Logger

	subscription  string
	resourceGroup string
}

func NewServer(cfg *Config) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config given: %w", err)
	}

	return &Server{
		logger:        cfg.Logger,
		subscription:  cfg.Subscription,
		resourceGroup: cfg.ResourceGroup,
	}, nil
}

var pathToTopic = map[string]string{
	"/events/processed_portfolio": "processed-portfolios",
}

func (s *Server) verifyWebhook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		topic, ok := pathToTopic[r.URL.Path]
		if !ok {
			s.logger.Error("no topic found for path", zap.String("path", r.URL.Path))
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if r.Header.Get("aeg-event-type") != "SubscriptionValidation" {
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
		if req.Topic != fullTopic {
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
	r.Post("/events/processed_portfolio", func(w http.ResponseWriter, r *http.Request) {
		dat, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error("failed to read webhook request body", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		s.logger.Info("processed porfolio", zap.String("portfolio_data", string(dat)))
	})
}