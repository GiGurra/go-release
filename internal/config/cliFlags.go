package config

import (
	"github.com/urfave/cli/v2"
)

var CliFlags = struct {
	Version                  cli.StringFlag
	IgnoreUncommittedChanges cli.BoolFlag
}{
	Version: cli.StringFlag{
		Name:    "output-version",
		Aliases: []string{"o"},
		Usage:   "to use for tagging. If not set, uses result of go build . && ./module-name -v",
	},
	IgnoreUncommittedChanges: cli.BoolFlag{
		Name:  "ignore-uncommitted-changes",
		Usage: "git status --porcelain check is ignored",
		Value: false,
	},
}
