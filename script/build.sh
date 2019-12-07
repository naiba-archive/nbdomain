export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -ldflags "-X main.licenseDomain=$1 -X main.licenseIP=$2" -o cmd/installer/nbdomain/nbdomain-linux-amd64 cmd/panel/main.go