package commands

import (
	githubcommands "github.com/mms-gianni/git-runner/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdClean() *clif.Command {
	cb := func(c *clif.Command, out clif.Output) {
		githubcommands.Clean(c, out)
	}

	return clif.NewCommand("clean", "Remove offline runners", cb)
}

func init() {
	Commands = append(Commands, cmdClean)
}
