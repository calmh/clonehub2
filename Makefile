all: clonehub-linux-amd64 clonehub-linux-arm64

clonehub-linux-amd64:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o clonehub-linux-amd64
clonehub-linux-arm64:
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -v -o clonehub-linux-arm64
