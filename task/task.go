// Package task holds domain types for asynchronous PACTA work, like analyzing profiles
package task

import (
	"bytes"

	"github.com/RMI/pacta/pacta"
)

// ID uniquely identifies a task being processed.
type ID string

// RunnerID also uniquely identifies a task being processed, but is specific to
// the underlying task runner (e.g. docker or AZ Container App Jobs)
type RunnerID string

type Type string

const (
	ProcessPortfolio = Type("process_portfolio")
	CreateReport     = Type("create_report")
)

type ProcessPortfolioRequest struct {
	// Note: This is temporary just to test the full end-to-end flow. We'll likely
	// want to reference assets by the portfolio (group?) they were uploaded to.
	AssetIDs []string
}

type ProcessPortfolioResponse struct {
	TaskID   ID
	AssetIDs []string
	Outputs  []string
}

type CreateReportRequest struct {
	PortfolioID pacta.PortfolioID
}

type EnvVar struct {
	Key   string
	Value string
}

type BaseImage struct {
	// Like 'rmipacta.azurecr.io'
	Registry string
	// Like 'runner'
	Name string
}

type Image struct {
	Base BaseImage
	// Like 'latest'
	Tag string
}

func (i *BaseImage) WithTag(tag string) string {
	var buf bytes.Buffer
	// <registry>/<name>:<tag>
	buf.WriteString(i.Registry)
	buf.WriteRune('/')
	buf.WriteString(i.Name)
	buf.WriteRune(':')
	buf.WriteString(tag)
	return buf.String()
}

func (i *Image) String() string {
	return i.Base.WithTag(i.Tag)
}

type Config struct {
	Env     []EnvVar
	Image   *Image
	Flags   []string
	Command []string
}
