package main

import (
	_ "embed"
	"os"

	"github.com/mms-gianni/git-runner/commands"
	"gopkg.in/ukautz/clif.v1"
)

func addDefaultOptions(cli *clif.Cli) {
	githubtoken := clif.NewOption("githubtoken", "t", "Private Github Token", "", true, false).
		SetEnv("GITHUB_TOKEN")

	githubusername := clif.NewOption("username", "u", "Github username", "", false, false).
		SetEnv("GITHUB_USERNAME")

	githubOrganisations := clif.NewOption("organisations", "o", "Github organisations (comma separated)", "", false, false).
		SetEnv("GITHUB_ORGANISATIONS")
	cli.AddDefaultOptions(githubtoken, githubusername, githubOrganisations)
}

//go:embed VERSION
var version string

func main() {
	cli := clif.New("git-runner", version, "Manage your github runners with git cli")

	var OwnStyles = map[string]string{
		"error":     "\033[31;1m",
		"warn":      "\033[33m",
		"info":      "\033[0;97m",
		"success":   "\033[32m",
		"debug":     "\033[30;1m",
		"headline":  "\033[4;1m",
		"subline":   "\033[4m",
		"important": "\033[47;30;1m",
		"query":     "\033[36m",
		"reset":     "\033[0m",
		"online":    "\U00002705",
		"offline":   "\U0001F480",
		"busy":      "\U0001F525",
		"ok":        "\U00002705",
		"err":       "\U000026D4",
	}

	cli.SetOutput(clif.NewColorOutput(os.Stdout).SetFormatter(clif.NewDefaultFormatter(OwnStyles)))

	addDefaultOptions(cli)

	for _, cb := range commands.Commands {
		cli.Add(cb())
	}

	cli.Run()
}
