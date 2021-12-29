package config

import (
	"github.com/urfave/cli/v2"
)

var CliFlags = struct {
	Version cli.StringFlag
}{
	Version: cli.StringFlag{
		Name:    "output-version",
		Aliases: []string{"o"},
		Usage:   "to use for tagging. If not set, uses result of go build . && ./module-name -v",
	},
}
