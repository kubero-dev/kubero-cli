package version

import _ "embed"

const gitModelUrl = "https://github.com/kubero-dev/kubero"

//go:embed VERSION
var embedVersion string

var (
	version = embedVersion
)

const (
	fixedInternalVersion = "0.0.1"
)
