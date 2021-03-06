package app

import (
	"context"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"gopkg.in/ukautz/clif.v1"
)

var ctx = context.Background()

var o clif.Output
var i clif.Input

func login(c *clif.Command) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Option("githubtoken").String()},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func GetStatus(c *clif.Command, out clif.Output) {
	client := login(c)
	o = out

	_, repo := GetGitdir()

	if repo != nil {
		repodetails := getRepodetails(repo)
		repoRunners, _, _ := client.Actions.ListRunners(ctx, repodetails.owner, repodetails.name, nil)
		title := "(" + repodetails.owner + ") " + repodetails.name
		printRunners(repoRunners, title)
	}

	//client.Actions.ListRunners(ctx , "mmds")
	organisationRunners, _, _ := client.Actions.ListOrganizationRunners(ctx, "mmz-srf", nil)
	printRunners(organisationRunners, "mmz-srf")
}

func printRunners(runners *github.Runners, title string) {

	if runners == nil {
		return
	}
	o.Printf("<important>%s<reset>\n", title)
	statusicon := ""
	for _, runner := range runners.Runners {

		if runner.GetStatus() == "offline" {
			statusicon = "offline"
		} else {
			if runner.GetBusy() {
				statusicon = "busy"
			} else {
				statusicon = "online"
			}
		}
		o.Printf("   <%s> %s\n", statusicon, runner.GetName())
	}
}
