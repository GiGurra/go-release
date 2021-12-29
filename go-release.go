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
		run(&appConfig)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(appConfig *config.AppConfig) {
	checkUncommittedChanges(appConfig)
	buildModule(appConfig)
	figureOutModuleName(appConfig)
	figureOutVersion(appConfig)
	tagInGitAndPush(appConfig)
	log.Printf("----------------------------------------")
	log.Printf("go-release succeeded! config was:\n%+v", appConfig)
}

func tagInGitAndPush(appConfig *config.AppConfig) {

	log.Printf("Creating git tag %s...", appConfig.Version)
	shell.RunCommand("git", "tag", appConfig.Version)
	log.Printf("ok")

	log.Printf("Pushing module=%s version=%s to remote origin...", appConfig.Module, appConfig.Version)
	shell.RunCommand("git", "push", "origin", appConfig.Version)
	log.Printf("ok")
}

func checkUncommittedChanges(appConfig *config.AppConfig) {
	if !appConfig.IgnoreUncommittedChanges {
		log.Printf("Checking for uncommitted changes...")
		if hasUncommittedChanges() {
			log.Fatalf("Cannot release because repo has uncommitted changes\n")
		}
		log.Printf("ok")
	}
}

func buildModule(appConfig *config.AppConfig) {
	log.Printf("Building module...")
	shell.RunCommand("go", "build", ".")
	log.Printf("ok")
}

func figureOutModuleName(appConfig *config.AppConfig) {
	log.Printf("Finding module name...")
	appConfig.Module = getCurrentModuleName()
	log.Printf("ok: %s", appConfig.Module)
}

func figureOutVersion(appConfig *config.AppConfig) {
	if appConfig.Version == "" {
		log.Printf("Finding version...")
		appConfig.Version = getCurrentModuleVersion(appConfig.Module)
		log.Printf("ok: %s", appConfig.Version)
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
