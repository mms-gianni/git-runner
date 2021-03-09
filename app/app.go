package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

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

func Run(c *clif.Command, out clif.Output) {
	client := login(c)
	o = out

	_, repo := GetGitdir()
	if repo != nil {
		repodetails := getRepodetails(repo)
		downloads, _, _ := client.Actions.ListRunnerApplicationDownloads(ctx, repodetails.owner, repodetails.name)

		repoToke, _, _ := client.Actions.CreateRegistrationToken(ctx, repodetails.owner, repodetails.name)
		url := "https://github.com/" + repodetails.owner + "/" + repodetails.name
		for _, download := range downloads {
			if download.GetOS() == "osx" && download.GetArchitecture() == "x64" && runtime.GOOS == "darwin" && runtime.GOARCH == "amd64" {
				runOnDarwin(out, download, url, repoToke)
			}
		}
	}
}

func runOnDarwin(out clif.Output, download *github.RunnerApplicationDownload, url string, token *github.RegistrationToken) {
	var err error
	var msg string

	fmt.Println(download.GetArchitecture(), download.GetOS(), runtime.GOOS, runtime.GOARCH)

	msg = "Dwonload runner binaries"
	out.Printf("    run: %s", msg)
	err = downloadRunner(download.GetDownloadURL(), "/tmp/runner.osx.tar.gz")
	if err == nil {
		out.Printf("\r <ok> run: %s\n", msg)
	} else {
		out.Printf("\r <err> run: %s => <error>%s\n", msg, err)
		return
	}

	msg = "rm -rf /tmp/runner.osx"
	out.Printf("    run: %s", msg)
	err = exec.Command("rm", "-rf", "/tmp/runner.osx").Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", msg)
	} else {
		out.Printf("\r <err> run: %s => <error>%s\n", msg, err)
		return
	}

	msg = "mkdir /tmp/runner.osx"
	out.Printf("    run: %s", msg)
	err = exec.Command("mkdir", "/tmp/runner.osx").Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", msg)
	} else {
		out.Printf("\r <err> run: %s => <error>%s\n", msg, err)
		return
	}

	msg = "tar -xzf /tmp/runner.osx.tar.gz -C /tmp/runner.osx"
	out.Printf("    run: %s", msg)
	err = exec.Command("tar", "-xzf", "/tmp/runner.osx.tar.gz", "-C", "/tmp/runner.osx").Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", msg)
	} else {
		out.Printf("\r <err> run: %s => <error>%s\n", msg, err)
		return
	}

	msg = "rm -f /tmp/runner.osx.tar.gz"
	out.Printf("    run: %s", msg)
	err = exec.Command("rm", "-f", "/tmp/runner.osx.tar.gz").Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", msg)
	} else {
		out.Printf("\r <err> run: %s => <error>%s\n", msg, err)
		return
	}

	msg = "/tmp/runner.osx/config.sh --url " + url + " --token " + token.GetToken()
	out.Printf("    run: %s", msg)
	err = exec.Command("/tmp/runner.osx/config.sh", "--url", url, "--token", token.GetToken()).Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", msg)
	} else {
		out.Printf("\r <err> run: %s => <error>%s\n", msg, err)
		return
	}

	msg = "/tmp/runner.osx/run.sh"
	out.Printf("    run: %s", msg)
	err = exec.Command("/tmp/runner.osx/run.sh").Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", msg)
	} else {
		out.Printf("\r <err> run: %s => %s\n", msg, err)
		return
	}
}

func downloadRunner(url string, filepath string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func Clean(c *clif.Command, out clif.Output) {
	client := login(c)
	o = out

	_, repo := GetGitdir()

	if repo != nil {
		repodetails := getRepodetails(repo)
		repoRunners, _, _ := client.Actions.ListRunners(ctx, repodetails.owner, repodetails.name, nil)
		title := "(" + repodetails.owner + ") " + repodetails.name
		o.Printf("<important>%s<reset>\n", title)
		for _, repoRunner := range repoRunners.Runners {
			if repoRunner.GetStatus() == "offline" {
				client.Actions.RemoveRunner(ctx, repodetails.owner, repodetails.name, repoRunner.GetID())
				o.Printf("   <offline> %s  => %s \n", repoRunner.GetName(), "removed")
			}
		}
	}

	organisations := c.Option("organisations").String()
	for _, organisation := range strings.Split(organisations, ",") {
		organisationRunners, _, _ := client.Actions.ListOrganizationRunners(ctx, organisation, nil)
		if organisationRunners != nil {
			title := "(" + organisation + ") "
			o.Printf("<important>%s<reset>\n", title)
			for _, orgRunner := range organisationRunners.Runners {
				if orgRunner.GetStatus() == "offline" {
					client.Actions.RemoveOrganizationRunner(ctx, organisation, orgRunner.GetID())
					o.Printf("   <offline> %s  => %s \n", orgRunner.GetName(), "removed")
				}
			}
		}
	}

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

	organisations := c.Option("organisations").String()
	for _, organisation := range strings.Split(organisations, ",") {
		organisationRunners, _, _ := client.Actions.ListOrganizationRunners(ctx, organisation, nil)
		printRunners(organisationRunners, organisation)
	}
}

func printRunners(runners *github.Runners, title string) {

	if runners == nil {
		return
	}
	o.Printf("  <important>%s<reset>\n", title)
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
	o.Printf("\n")
}
