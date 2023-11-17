package pactasrv

import (
	"context"
	"fmt"

	api "github.com/RMI/pacta/openapi/pacta"
)

// Starts the process of uploading one or more portfolio files
// (POST /portfolio-upload)
func (s *Server) StartPortfolioUpload(ctx context.Context, request api.StartPortfolioUploadRequestObject) (api.StartPortfolioUploadResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// Called after uploads of portfolios to cloud storage are complete.
// (POST /portfolio-upload:complete)
func (s *Server) CompletePortfolioUpload(ctx context.Context, request api.CompletePortfolioUploadRequestObject) (api.CompletePortfolioUploadResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// (GET /incomplete-uploads)
func (s *Server) ListIncompleteUploads(ctx context.Context, request api.ListIncompleteUploadsRequestObject) (api.ListIncompleteUploadsResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
