# git-runner
Manage your github runner with your git cli

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mms-gianni/git-runner)
![GitHub top language](https://img.shields.io/github/languages/top/mms-gianni/git-runner)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/mms-gianni/git-runner/Upload%20Release%20Asset)
![GitHub MIT license](https://img.shields.io/github/license/mms-gianni/git-runner)
![Swiss made](https://img.shields.io/badge/swiss%20made-100%25-red)
## Why
- use your local workstation as a runner
- simplify runner startup
- simplify runner maintenance
- cli to automate runner management

![Screenshot](docs/img/screenshot.png?raw=true "Screenshot")

## Installation
Generate a token here : https://github.com/settings/tokens (You need to be loged in)

To export the Github username and organisation is optional. 
### Mac
```
echo 'export GITHUB_TOKEN="XXXXXXXXXXXXXXXXXXXXXXX"' >> ~/.zshrc
echo 'export GITHUB_USERNAME="change-me-to-your-username"' >> ~/.zshrc
echo 'export GITHUB_ORGANISATIONS="klustair,kubernetes"' >> ~/.zshrc
curl https://raw.githubusercontent.com/mms-gianni/git-runner/master/cmd/git-runner/git-runner.mac.64bit -o /usr/local/bin/git-runner
chmod +x /usr/local/bin/git-runner
```

### Linux 
```
echo 'export GITHUB_TOKEN="XXXXXXXXXXXXXXXXXXXXXXX"' >> ~/.bashrc
echo 'export GITHUB_USERNAME="change-me-to-your-username"' >> ~/.bashrc
echo 'export GITHUB_ORGANISATIONS="klustair,kubernetes"' >> ~/.bashrc
curl https://raw.githubusercontent.com/mms-gianni/git-runner/master/cmd/git-runner/git-runner.linux.64bit -o /usr/local/bin/git-runner
chmod +x /usr/local/bin/git-runner
```

### Windows
Windows is not implemented yet. But I'm working on it. Pullrequests wellcome. 

## Quick start

### Show status of runners
Display all runners attached to this repository.
```
cd /path/to/your/repo
git runner status
```

### Remove dead runners
```
cd /path/to/your/repo
git project clean 
```

### Start a new runner
```
cd /path/to/your/repo
git project run
```
