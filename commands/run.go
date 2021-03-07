package commands

import (
	githubcommands "github.com/mms-gianni/git-runner/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdRun() *clif.Command {
	cb := func(c *clif.Command, out clif.Output) {
		githubcommands.Run(c, out)
	}

	return clif.NewCommand("run", "Start a new runner", cb)
}

func init() {
	Commands = append(Commands, cmdRun)
}
