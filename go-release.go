package main

import (
	"github.com/gigurra/go-release/internal/cliUtil"
	"github.com/gigurra/go-release/internal/config"
	"github.com/gigurra/go-release/internal/fileutil"
	"github.com/gigurra/go-release/internal/shell"
	"github.com/gigurra/go-release/internal/stringutil"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	app := cli.NewApp()

	app.HideHelpCommand = true
	app.Usage = "Create a simple release (git tag, push)"
	app.Flags = cliUtil.FindAllFlags(&config.CliFlags)
	app.Version = Version
	app.Action = func(c *cli.Context) error {
		appConfig := config.GetDefaultAppConfig()
		appConfig.Version = c.String(config.CliFlags.Version.Name)
		appConfig.IgnoreUncommittedChanges = c.Bool(config.CliFlags.IgnoreUncommittedChanges.Name)
		complementConfig(&appConfig)
		appConfig.Validate()
		run(appConfig)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(appConfig config.AppConfig) {

	if !appConfig.IgnoreUncommittedChanges && hasUncommittedChanges() {
		log.Fatalf("Cannot release because repo has uncommitted changes\n")
	}

	log.Printf("Building %+v...\n", appConfig)
}

func complementConfig(appConfig *config.AppConfig) {
	appConfig.Module = getCurrentModuleName()
	if appConfig.Version == "" {
		appConfig.Version = getCurrentModuleVersion(appConfig.Module)
	}
}

func getCurrentModuleName() string {
	firstLineOfGoMod := stringutil.SplitLines(strings.TrimSpace(fileutil.File2String("go.mod")))[0]
	parts := strings.Split(firstLineOfGoMod, "/")
	return parts[len(parts)-1]
}

func getCurrentModuleVersion(module string) string {
	splitter := regexp.MustCompile(`\s+`)
	commandResult := shell.RunCommand("./"+module, "--version")
	parts := splitter.Split(commandResult, -1)
	return parts[len(parts)-1]
}

func hasUncommittedChanges() bool {
	commandResult := shell.RunCommand("git", "status", "--porcelain")
	return commandResult != ""
}
