package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
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

func Kill(c *clif.Command, out clif.Output) {
	o = out

	cmd := exec.Command("pgrep", "Runner.Listener")
	pid, err := cmd.Output()
	if err != nil {
		out.Printf(" <err> run: %s => <error>%s<reset>\n", cmd.String(), err)
		return
	}

	cmd = exec.Command("kill", string(pid[:len(pid)-1]))
	out.Printf("    run: %s", cmd.String())
	err = cmd.Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", cmd.String())
	} else {
		out.Printf("\r <err> run: %s => <error>%s<reset>\n", cmd.String(), err)
		return
	}
}

func Run(c *clif.Command, in clif.Input, out clif.Output) {
	client := login(c)
	o = out
	var runnerDestination map[string]string
	runnerDestination = make(map[string]string)
	var orgs []string

	detached := c.Option("detached").IsFlag().Bool()

	if c.Option("organisations").String() != "" {
		organisations := c.Option("organisations").String()
		orgs = strings.Split(organisations, ",")
		for i := 0; i < len(orgs); i++ {
			runnerDestination[strconv.Itoa(i+1)] = "Orga: " + orgs[i]
		}
	}

	_, repo := GetGitdir()

	if repo != nil {
		repodetails := getRepodetails(repo)
		runnerDestination[strconv.Itoa(len(runnerDestination)+1)] = "Repo: " + repodetails.name + "/" + repodetails.name
	}

	selectedDestination := "1"

	if len(runnerDestination) > 1 {
		selectedDestination = in.Choose("where do you want to start the runner: ", runnerDestination)
	}

	selectedDestI, _ := strconv.Atoi(selectedDestination)
	if repo != nil && selectedDestI > len(orgs) {
		repodetails := getRepodetails(repo)
		downloads, _, _ := client.Actions.ListRunnerApplicationDownloads(ctx, repodetails.owner, repodetails.name)

		token, _, _ := client.Actions.CreateRegistrationToken(ctx, repodetails.owner, repodetails.name)
		url := "https://github.com/" + repodetails.owner + "/" + repodetails.name
		runOnOs(out, downloads, url, token, detached, c.Option("labels").String())
	}

	if selectedDestI <= len(orgs) {
		owner := orgs[selectedDestI-1]
		downloads, _, _ := client.Actions.ListOrganizationRunnerApplicationDownloads(ctx, owner)

		token, _, _ := client.Actions.CreateOrganizationRegistrationToken(ctx, owner)
		url := "https://github.com/" + owner
		runOnOs(out, downloads, url, token, detached, c.Option("labels").String())
	}
}

func runOnOs(out clif.Output, downloads []*github.RunnerApplicationDownload, url string, token *github.RegistrationToken, detached bool, labels string) {
	for _, download := range downloads {
		// https://stackoverflow.com/questions/20728767/all-possible-goos-value
		// possible GetOS values: osx, linux, win
		// possible GetArchitecture values: x64, arm, arm64
		if download.GetOS() == "osx" && download.GetArchitecture() == "x64" && runtime.GOOS == "darwin" && runtime.GOARCH == "amd64" {
			out.Printf("Start installation for OS:%s Arch:%s\n\n", download.GetOS(), download.GetArchitecture())
			runOnPosix(out, download, url, token, detached, labels)
			break
		}
		if download.GetOS() == "linux" && download.GetArchitecture() == "x64" && runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
			out.Printf("Start installation for OS:%s Arch:%s\n\n", download.GetOS(), download.GetArchitecture())
			runOnPosix(out, download, url, token, detached, labels)
			break
		}
		if download.GetOS() == "linux" && download.GetArchitecture() == "arm" && runtime.GOOS == "linux" && runtime.GOARCH == "arn" {
			out.Printf("Start installation for OS:%s Arch:%s\n\n", download.GetOS(), download.GetArchitecture())
			runOnPosix(out, download, url, token, detached, labels)
			break
		}
		if download.GetOS() == "linux" && download.GetArchitecture() == "arm64" && runtime.GOOS == "linux" && runtime.GOARCH == "arm64" {
			out.Printf("Start installation for OS:%s Arch:%s\n\n", download.GetOS(), download.GetArchitecture())
			runOnPosix(out, download, url, token, detached, labels)
			break
		}
		if download.GetOS() == "win" && download.GetArchitecture() == "x64" && runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
			out.Printf("Start installation for OS:%s Arch:%s\n\n", download.GetOS(), download.GetArchitecture())
			out.Printf("===> Sorry not implemented yet")
			break
		}
	}
}

func runOnPosix(out clif.Output, download *github.RunnerApplicationDownload, url string, token *github.RegistrationToken, detached bool, labels string) {
	var err error
	var cmd *exec.Cmd

	//fmt.Println(download.GetArchitecture(), download.GetOS(), runtime.GOOS, runtime.GOARCH)

	msg := "curl -O -L " + download.GetDownloadURL()
	out.Printf("    run: %s", msg)
	err = downloadRunner(download.GetDownloadURL(), "/tmp/github_runner.tar.gz")
	if err == nil {
		out.Printf("\r <ok> run: %s\n", msg)
	} else {
		out.Printf("\r <err> run: %s => <error>%s<reset>\n", msg, err)
		return
	}

	cmd = exec.Command("rm", "-rf", "/tmp/github_runner")
	out.Printf("    run: %s", cmd.String())
	err = cmd.Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", cmd.String())
	} else {
		out.Printf("\r <err> run: %s => <error>%s<reset>\n", cmd.String(), err)
		return
	}

	cmd = exec.Command("mkdir", "/tmp/github_runner")
	out.Printf("    run: %s", cmd.String())
	err = cmd.Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", cmd.String())
	} else {
		out.Printf("\r <err> run: %s => <error>%s<reset>\n", cmd.String(), err)
		return
	}

	cmd = exec.Command("tar", "-xzf", "/tmp/github_runner.tar.gz", "-C", "/tmp/github_runner")
	out.Printf("    run: %s", cmd.String())
	err = cmd.Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", cmd.String())
	} else {
		out.Printf("\r <err> run: %s => <error>%s<reset>\n", cmd.String(), err)
		return
	}

	cmd = exec.Command("rm", "-f", "/tmp/github_runner.tar.gz")
	out.Printf("    run: %s", cmd.String())
	err = cmd.Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", cmd.String())
	} else {
		out.Printf("\r <err> run: %s => <error>%s<reset>\n", cmd.String(), err)
		return
	}

	labelsoption := ""
	if labels != "" {
		labelsoption = "--labels"
	}

	cmd = exec.Command("/tmp/github_runner/config.sh", "--unattended", "--replace", "--url", url, "--token", token.GetToken(), labelsoption, labels)
	out.Printf("    run: %s", cmd.String())
	err = cmd.Run()
	if err == nil {
		out.Printf("\r <ok> run: %s\n", cmd.String())
	} else {
		out.Printf("\r <err> run: %s => <error>%s<reset>\n", cmd.String(), err)
		return
	}

	cmd = exec.Command("/tmp/github_runner/run.sh")
	out.Printf("    run: %s", cmd.String())

	stdout, _ := cmd.StdoutPipe()
	err = cmd.Start()

	if err == nil {
		out.Printf("\r <ok> run: %s\n", cmd.String())
	} else {
		out.Printf("\r <err> run: %s => %s\n", cmd.String(), err)
		return
	}

	if !detached {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			out.Printf(string(buf[0:n]))
		}

		if err := cmd.Wait(); err != nil {
			out.Printf("\r <err> run: %s => %s\n", cmd.String(), err)
		}
	} else {
		out.Printf("run 'kill $(pgrep Runner.Listener)' to remove the runner")
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
		o.Printf("   <%s> %s   <debug>", statusicon, runner.GetName())
		for _, label := range runner.Labels {
			o.Printf("[%s] ", label.GetName())
		}
		o.Printf("<reset>\n")
	}
	o.Printf("\n")
}
