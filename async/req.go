package async

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/RMI/pacta/task"
)

func ParsePortfolioReq() (*task.ParsePortfolioRequest, error) {
	taskStr := os.Getenv("PARSE_PORTFOLIO_REQUEST")
	if taskStr == "" {
		return nil, errors.New("no PARSE_PORTFOLIO_REQUEST given")
	}
	var task task.ParsePortfolioRequest
	if err := json.NewDecoder(strings.NewReader(taskStr)).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to load ParsePortfolioRequest: %w", err)
	}
	return &task, nil
}

func CreateAuditReq() (*task.CreateAuditRequest, error) {
	car := os.Getenv("CREATE_AUDIT_REQUEST")
	if car == "" {
		return nil, errors.New("no CREATE_AUDIT_REQUEST was given")
	}
	var task task.CreateAuditRequest
	if err := json.NewDecoder(strings.NewReader(car)).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to load CreateAuditRequest: %w", err)
	}
	return &task, nil
}

func CreateReportReq() (*task.CreateReportRequest, error) {
	crr := os.Getenv("CREATE_REPORT_REQUEST")
	if crr == "" {
		return nil, errors.New("no CREATE_REPORT_REQUEST was given")
	}
	var task task.CreateReportRequest
	if err := json.NewDecoder(strings.NewReader(crr)).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to load CreateReportRequest: %w", err)
	}
	return &task, nil
}
