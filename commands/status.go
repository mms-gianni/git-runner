package commands

import (
	githubcommands "github.com/mms-gianni/git-runner/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdStatus() *clif.Command {
	cb := func(c *clif.Command, out clif.Output) {
		githubcommands.GetStatus(c, out)
	}

	return clif.NewCommand("status", "List runners", cb)
}

func init() {
	Commands = append(Commands, cmdStatus)
}
