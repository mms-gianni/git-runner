FROM alpine

RUN apk add --no-cache libc6-compat

COPY cmd/git-runner/git-runner.linux.64bit /usr/local/bin/git-runner
RUN chmod +x /usr/local/bin/git-runner

#ENTRYPOINT [ "/usr/local/bin/git-runner" ]
