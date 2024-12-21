package async

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/RMI/pacta/task"
)

func LoadParsePortfolioRequestFromEnv() (*task.ParsePortfolioRequest, error) {
	return loadFromEnv[task.ParsePortfolioRequest]("PARSE_PORTFOLIO_REQUEST", "ParsePortfolioRequest")
}

func LoadCreateDashboardRequestFromEnv() (*task.CreateDashboardRequest, error) {
	return loadFromEnv[task.CreateDashboardRequest]("CREATE_DASHBOARD_REQUEST", "CreateDashboardRequest")
}

func LoadCreateAuditRequestFromEnv() (*task.CreateAuditRequest, error) {
	return loadFromEnv[task.CreateAuditRequest]("CREATE_AUDIT_REQUEST", "CreateAuditRequest")
}

func LoadCreateReportRequestFromEnv() (*task.CreateReportRequest, error) {
	return loadFromEnv[task.CreateReportRequest]("CREATE_REPORT_REQUEST", "CreateReportRequest")
}

func loadFromEnv[T any](envVar string, entityName string) (*T, error) {
	envStr := os.Getenv(envVar)
	if envStr == "" {
		return nil, errors.New("no CREATE_REPORT_REQUEST was given")
	}
	var task T
	if err := json.NewDecoder(strings.NewReader(envStr)).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to load %q: %w", entityName, err)
	}
	return &task, nil
}
