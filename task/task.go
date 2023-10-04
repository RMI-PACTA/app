// Package task holds domain types for asynchronous PACTA work, like analyzing profiles
package task

import "github.com/RMI/pacta/pacta"

// TaskID uniquely identifies a task being processed.
type ID string

type StartRunRequest struct {
	// This is just an example, this will likely change over time.
	PortfolioID pacta.PortfolioID
}
