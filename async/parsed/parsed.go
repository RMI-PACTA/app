// Package parsed just holds the domain types for dealing with the output of the
// ParsePortfolio async task.
package parsed

type SourceFile struct {
	InputFilename      string      `json:"input_filename"`
	InputMD5           string      `json:"input_md5"`
	SystemInfo         SystemInfo  `json:"system_info"`
	InputEntries       int         `json:"input_entries"`
	GroupCols          []string    `json:"group_cols"`
	SubportfoliosCount int         `json:"subportfolios_count"`
	Portfolios         []Portfolio `json:"portfolios"`
	Errors             [][]string  `json:"errors"`
}

type SystemInfo struct {
	Timestamp      string       `json:"timestamp"`
	Package        string       `json:"package"`
	PackageVersion string       `json:"packageVersion"`
	RVersion       string       `json:"RVersion"`
	Dependencies   []Dependency `json:"dependencies"`
}

type Dependency struct {
	Package string `json:"package"`
	Version string `json:"version"`
}

type Portfolio struct {
	OutputMD5      string `json:"output_md5"`
	OutputFilename string `json:"output_filename"`
	OutputRows     int    `json:"output_rows"`
	PortfolioName  string `json:"portfolio_name"`
	InvestorName   string `json:"investor_name"`
}
