
all: build

linux:
	@echo " Building for Linux"
	GOOS=linux GOARCH=amd64 go build -o ./cmd/git-runner/git-runner.linux.64bit -ldflags "-s -X main.Version=$(version)" main/main.go

mac:
	@echo " Building for Mac"
	GOOS=darwin GOARCH=amd64 go build -o ./cmd/git-runner/git-runner.mac.64bit -ldflags "-s -X main.Version=$(version)" main/main.go

windows:
	@echo " Building for Windows"
	GOOS=windows GOARCH=amd64 go build -o ./cmd/git-runner/git-runner.windows.64bit -ldflags "-s -X main.Version=$(version)" main/main.go

build_pre:
	@echo "Building all"

build: build_pre linux mac windows
	@echo "Done"

clean:
	@echo "Cleaning up"
	@rm -f ./cmd/git-runner/git-runner.linux.64bit
	@rm -f ./cmd/git-runner/git-runner.mac.64bit
	@rm -f ./cmd/git-runner/git-runner.windows.64bit
	@echo "Done"

