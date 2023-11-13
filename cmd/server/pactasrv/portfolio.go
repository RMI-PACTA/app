package pactasrv

import (
	"context"

	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/task"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Server) CreatePortfolioAsset(ctx context.Context, req api.CreatePortfolioAssetRequestObject) (api.CreatePortfolioAssetResponseObject, error) {
	id := uuid.NewString()
	uri := blob.Join(s.Blob.Scheme(), s.PorfolioUploadURI, id)
	signed, err := s.Blob.SignedUploadURL(ctx, uri)
	if err != nil {
		return nil, oapierr.Internal("failed to sign blob URI", zap.String("uri", uri), zap.Error(err))
	}
	return api.CreatePortfolioAsset200JSONResponse{
		UploadUrl: signed,
		AssetId:   id,
	}, nil
}

func (s *Server) ProcessPortfolio(ctx context.Context, req api.ProcessPortfolioRequestObject) (api.ProcessPortfolioResponseObject, error) {
	taskID, runnerID, err := s.TaskRunner.ProcessPortfolio(ctx, &task.ProcessPortfolioRequest{
		AssetIDs: req.Body.AssetIds,
		// PortfolioID: req.Body.PortfolioID,
	})
	if err != nil {
		return nil, oapierr.Internal("failed to start task", zap.Error(err))
	}
	s.Logger.Info("triggered process portfolio task",
		zap.String("task_id", string(taskID)),
		zap.String("task_runner_id", string(runnerID)))
	return api.ProcessPortfolio200JSONResponse{
		TaskId: string(taskID),
	}, nil
}
