FROM alpine

RUN apk add --no-cache libc6-compat ca-certificates

COPY git-runner /usr/local/bin/git-runner
RUN chmod +x /usr/local/bin/git-runner

#ENTRYPOINT [ "/usr/local/bin/git-runner" ]
