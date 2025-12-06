# echo -------------------------------------------------------------------------
# echo Description: build go binary with flags 
# echo -------------------------------------------------------------------------

VERSION=1.0

# compile binary for production 
prod() {

    CGO_ENABLED=0 go build ./cmd/agent/agent.go -ldflags="-s -w -X main.version=$VERSION" -o agent

}

# compile binary in debug mode
debug() {

     go build -race -o debug ./cmd/agent/agent.go 
}


if [ "$1" = "debug" ]; then 
    debug && exit
elif [ "$1" = "prod" ]; then
    prod && exit
fi

echo "Usage: [debug|prod]" && exit 1
