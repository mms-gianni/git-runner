package commands

import (
	githubcommands "github.com/mms-gianni/git-runner/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdRun() *clif.Command {
	cb := func(c *clif.Command, in clif.Input, out clif.Output) {
		githubcommands.Run(c, in, out)
	}

	return clif.NewCommand("run", "Start a new runner", cb).
		NewFlag("detached", "d", "Start a runner and detach (needs a manual kill)", false)
}

func init() {
	Commands = append(Commands, cmdRun)
}
