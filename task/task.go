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
	ParsePortfolio = Type("parse_portfolio")
	CreateReport   = Type("create_report")
	CreateAudit    = Type("create_audit")
)

type ParsePortfolioRequest struct {
	IncompleteUploadIDs []pacta.IncompleteUploadID
	BlobURIs            []pacta.BlobURI
}

type ParsePortfolioResponseItem struct {
	Blob      pacta.Blob
	LineCount int
}

type ParsePortfolioResponse struct {
	TaskID  ID
	Request *ParsePortfolioRequest
	Outputs []*ParsePortfolioResponseItem
}

type CreateAuditRequest struct {
	AnalysisID pacta.AnalysisID
	BlobURIs   []pacta.BlobURI
}

type CreateAuditResponse struct {
	TaskID  ID
	Request *CreateAuditRequest
}

type CreateReportRequest struct {
	AnalysisID pacta.AnalysisID
	BlobURIs   []pacta.BlobURI
}

type CreateReportResponse struct {
	TaskID  ID
	Request *CreateReportRequest
}

type EnvVar struct {
	Key   string
	Value string
}

type BaseImage struct {
	// Like 'rmisa.azurecr.io'
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
