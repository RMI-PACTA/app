package reportsrv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/session"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap/zaptest"
)

func TestServeReport(t *testing.T) {

	srv, env := setup(t)
	router := chi.NewRouter()
	srv.RegisterHandlers(router)

	aIDStr := "analysis.id1"
	analysisID := pacta.AnalysisID(aIDStr)
	ownerID := pacta.OwnerID("owner.id1")
	userID := pacta.UserID("user.id1")
	otherUserID := pacta.UserID("user.id2")
	adminUserID := pacta.UserID("user.id3")

	env.db.users = []*pacta.User{{
		ID: userID,
	}, {
		ID: otherUserID,
	}, {
		ID:    adminUserID,
		Admin: true,
	}}

	env.db.userToOwner = map[pacta.UserID]pacta.OwnerID{
		userID:      ownerID,
		otherUserID: "other.owner.id",
		adminUserID: "admin.owner.id",
	}

	env.db.analyses = []*pacta.Analysis{
		&pacta.Analysis{
			ID:    analysisID,
			Owner: &pacta.Owner{ID: ownerID},
		},
	}

	env.db.analysisArtifacts = []*pacta.AnalysisArtifact{
		&pacta.AnalysisArtifact{
			ID:                "analysisartifact.id1",
			AnalysisID:        analysisID,
			Blob:              &pacta.Blob{ID: "blob.id1"},
			AdminDebugEnabled: true,
		},
		&pacta.AnalysisArtifact{
			ID:         "analysisartifact.id2",
			AnalysisID: analysisID,
			Blob:       &pacta.Blob{ID: "blob.id2"},
		},
	}

	env.db.blobs = map[pacta.BlobID]*pacta.Blob{
		"blob.id1": &pacta.Blob{
			ID:       "blob.id1",
			BlobURI:  "test://reports/1111-2222-3333-4444/index.html",
			FileType: pacta.FileType_HTML,
			FileName: "index.html",
		},
		"blob.id2": &pacta.Blob{
			ID:       "blob.id2",
			BlobURI:  "test://reports/1111-2222-3333-4444/lib/some/package.js",
			FileType: pacta.FileType_JS,
			FileName: "package.js",
		},
	}

	htmlContent := "<html>this is the index</html>"
	jsContent := "function() { return 'some js' }"
	env.blob.blobContents = map[string]string{
		"test://reports/1111-2222-3333-4444/index.html":          htmlContent,
		"test://reports/1111-2222-3333-4444/lib/some/package.js": jsContent,
	}

	standardPath := "/report/" + aIDStr + "/"
	cases := []struct {
		asUser          pacta.UserID
		path            string
		wantErr         int
		wantContentType string
		wantRespContent string
	}{{
		asUser:  userID,
		path:    "/report/" + aIDStr,
		wantErr: http.StatusTemporaryRedirect,
	}, {
		asUser:          userID,
		path:            standardPath,
		wantContentType: "text/html",
		wantRespContent: htmlContent,
	}, {
		asUser:          userID,
		path:            "/report/" + aIDStr + "/index.html",
		wantContentType: "text/html",
		wantRespContent: htmlContent,
	}, {
		asUser:          userID,
		path:            "/report/" + aIDStr + "/lib/some/package.js",
		wantContentType: "text/javascript",
		wantRespContent: jsContent,
	}, {
		asUser:          "",
		path:            standardPath,
		wantContentType: "text/html",
		wantErr:         http.StatusUnauthorized,
	}, {
		asUser:          otherUserID,
		path:            standardPath,
		wantContentType: "text/html",
		wantErr:         http.StatusUnauthorized,
	}, {
		asUser:          adminUserID,
		path:            standardPath,
		wantContentType: "text/html",
		wantRespContent: htmlContent,
	}, {
		asUser:          userID,
		path:            "/report/a-nonsense-report/",
		wantContentType: "text/html",
		wantErr:         http.StatusNotFound,
	}}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, &chi.Context{
				URLParams: chi.RouteParams{
					Keys:   []string{"analysis_id"},
					Values: []string{"analysis.id1"},
				},
			})
			if c.asUser != "" {
				ctx = session.WithUserID(ctx, c.asUser)
			}

			r := httptest.NewRequest(http.MethodGet, c.path, nil).WithContext(ctx)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			res := w.Result()
			// Successful response...
			wantCode := http.StatusOK
			if c.wantErr != 0 {
				wantCode = c.wantErr
			}
			gotCode := res.StatusCode
			if gotCode != wantCode {
				t.Errorf("got status code %d, want %d", gotCode, wantCode)
			}
			if c.wantErr != 0 {
				return
			}
			// ...with the correct content-type...
			if got, want := res.Header.Get("Content-Type"), c.wantContentType; got != want {
				t.Errorf("Content-Type = %q, want %q", got, want)
			}
			// ...and the correct response body.
			if got, want := w.Body.String(), c.wantRespContent; got != want {
				t.Errorf("Response body = %q, want %q", got, want)
			}
		})
	}
}

type testEnv struct {
	db   *testDB
	blob *testBlob
}

func setup(t *testing.T) (*Server, *testEnv) {
	env := &testEnv{
		db:   &testDB{},
		blob: &testBlob{},
	}

	return &Server{
		db:     env.db,
		blob:   env.blob,
		logger: zaptest.NewLogger(t),
	}, env
}

type testTx struct{}

func (testTx) Commit() error   { return nil }
func (testTx) Rollback() error { return nil }

type testDB struct {
	// Recording inputs
	gotAnalysisIDs []pacta.AnalysisID
	gotBlobIDs     [][]pacta.BlobID
	gotAuditLogs   []pacta.AuditLog

	// Hardcoded outputs
	analyses          []*pacta.Analysis
	analysisArtifacts []*pacta.AnalysisArtifact
	blobs             map[pacta.BlobID]*pacta.Blob
	userToOwner       map[pacta.UserID]pacta.OwnerID
	users             []*pacta.User
}

func (tdb *testDB) NoTxn(ctx context.Context) db.Tx {
	return testTx{}
}

func (tdb *testDB) AnalysisArtifactsForAnalysis(tx db.Tx, id pacta.AnalysisID) ([]*pacta.AnalysisArtifact, error) {
	tdb.gotAnalysisIDs = append(tdb.gotAnalysisIDs, id)
	return tdb.analysisArtifacts, nil
}

func (tdb *testDB) Blobs(tx db.Tx, ids []pacta.BlobID) (map[pacta.BlobID]*pacta.Blob, error) {
	tdb.gotBlobIDs = append(tdb.gotBlobIDs, ids)
	return tdb.blobs, nil
}

func (tdb *testDB) GetOwnerForUser(tx db.Tx, userID pacta.UserID) (pacta.OwnerID, error) {
	ownerID, ok := tdb.userToOwner[userID]
	if !ok {
		return "", fmt.Errorf("user %q not found", userID)
	}
	return ownerID, nil
}

func (tdb *testDB) CreateAuditLog(tx db.Tx, log *pacta.AuditLog) (pacta.AuditLogID, error) {
	tdb.gotAuditLogs = append(tdb.gotAuditLogs, *log)
	return pacta.AuditLogID("auditlog.id1"), nil
}

func (tdb *testDB) User(tx db.Tx, id pacta.UserID) (*pacta.User, error) {
	for _, u := range tdb.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user %q not found", id)
}

func (tdb *testDB) Analysis(tx db.Tx, id pacta.AnalysisID) (*pacta.Analysis, error) {
	for _, a := range tdb.analyses {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, fmt.Errorf("analysis %q not found", id)
}

type testBlob struct {
	// Recording input
	gotURIs []string

	// Hardcoded output, mapping of uri -> contents
	blobContents map[string]string
}

func (tb *testBlob) Scheme() blob.Scheme {
	return blob.Scheme("test")
}

type noopCloser struct {
	io.Reader
}

func (noopCloser) Close() error { return nil }

func (tb *testBlob) ReadBlob(ctx context.Context, uri string) (io.ReadCloser, error) {
	tb.gotURIs = append(tb.gotURIs, uri)
	return noopCloser{strings.NewReader(tb.blobContents[uri])}, nil
}
