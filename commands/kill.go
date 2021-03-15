package commands

import (
	githubcommands "github.com/mms-gianni/git-runner/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdKill() *clif.Command {
	cb := func(c *clif.Command, in clif.Input, out clif.Output) {
		githubcommands.Kill(c, out)
	}

	return clif.NewCommand("kill", "Kill local runner if running", cb)
}

func init() {
	Commands = append(Commands, cmdKill)
}
