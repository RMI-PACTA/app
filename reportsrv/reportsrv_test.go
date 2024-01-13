package reportsrv

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap/zaptest"
)

func TestServeReport(t *testing.T) {
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"analysis_id"},
			Values: []string{"analysis.id1"},
		},
	})

	srv, env := setup(t)

	env.db.analysisArtifacts = []*pacta.AnalysisArtifact{
		&pacta.AnalysisArtifact{
			ID:         "analysisartifact.id1",
			AnalysisID: "analysis.id1",
			Blob:       &pacta.Blob{ID: "blob.id1"},
		},
		&pacta.AnalysisArtifact{
			ID:         "analysisartifact.id2",
			AnalysisID: "analysis.id1",
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

	serveReportAndCheckResponse := func(path, contentType, responseBody string) {
		t.Run(path, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, path, nil).WithContext(ctx)
			w := httptest.NewRecorder()

			srv.serveReport(w, r)

			res := w.Result()
			// Successful response...
			if got, want := w.Code, http.StatusOK; got != want {
				t.Errorf("got status code %d, want %d", got, want)
			}
			// ...with the correct content-type...
			if got, want := res.Header.Get("Content-Type"), contentType; got != want {
				t.Errorf("Content-Type = %q, want %q", got, want)
			}
			// ...and the correct response body.
			if got, want := w.Body.String(), responseBody; got != want {
				t.Errorf("Response body = %q, want %q", got, want)
			}
		})
	}

	// First, request the root in a few different ways. All of then should return the /index.html file contents.
	serveReportAndCheckResponse("/report/analysis.id1", "text/html", htmlContent)
	serveReportAndCheckResponse("/report/analysis.id1/", "text/html", htmlContent)
	serveReportAndCheckResponse("/report/analysis.id1/index.html", "text/html", htmlContent)

	// Now, try loading an asset from the report.
	serveReportAndCheckResponse("/report/analysis.id1/lib/some/package.js", "text/javascript", jsContent)
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

	// Hardcoded outputs
	analysisArtifacts []*pacta.AnalysisArtifact
	blobs             map[pacta.BlobID]*pacta.Blob
}

func (tdb *testDB) Begin(ctx context.Context) (db.Tx, error) {
	return testTx{}, nil
}

func (tdb *testDB) NoTxn(ctx context.Context) db.Tx {
	return testTx{}
}

func (tdb *testDB) Transactional(ctx context.Context, fn func(tx db.Tx) error) error {
	return fn(testTx{})
}

func (tdb *testDB) RunOrContinueTransaction(tx db.Tx, fn func(tx db.Tx) error) error {
	return fn(testTx{})
}

func (tdb *testDB) AnalysisArtifactsForAnalysis(tx db.Tx, id pacta.AnalysisID) ([]*pacta.AnalysisArtifact, error) {
	tdb.gotAnalysisIDs = append(tdb.gotAnalysisIDs, id)
	return tdb.analysisArtifacts, nil
}

func (tdb *testDB) Blobs(tx db.Tx, ids []pacta.BlobID) (map[pacta.BlobID]*pacta.Blob, error) {
	tdb.gotBlobIDs = append(tdb.gotBlobIDs, ids)
	return tdb.blobs, nil
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
