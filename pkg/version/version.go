// Package version shows build information
package version

import (
	// init embed package
	_ "embed"
	"io"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
)

// we'll embed the version info into the file and make it available in the build
//
//go:embed version.txt
var hello string

var (
	// BuildVersion is a build version
	BuildVersion string
	// BuildDate is build time
	BuildDate string
	// BuildAuthor is the author of the build
	BuildAuthor string
	// BuildCommit is the commit hash
	BuildCommit string
)

type buildInfo struct {
	Version string
	Date    string
	Author  string
	Commit  string
}

// WriteBuildInfo write build info
func WriteBuildInfo(w io.Writer) {
	info := buildInfo{
		Version: "0.0.1", // N/A - default
		Date:    "2025/06/02",
		Author:  "Painkiller675",
		Commit:  "12345678",
	}

	if BuildVersion != "" {
		info.Version = BuildVersion
	}
	if BuildDate != "" {
		info.Date = BuildDate
	}
	if BuildCommit != "" {
		info.Commit = BuildCommit
	}
	// generate version info
	tmpl := template.Must(template.New("version").Parse(hello))
	if err := tmpl.Execute(w, info); err != nil {
		log.Warn().Err(err).Msg("Failed to print build info")
	}
}

// Info returns the string with build information
func Info() string {
	builder := &strings.Builder{}
	WriteBuildInfo(builder)
	return builder.String()
}
