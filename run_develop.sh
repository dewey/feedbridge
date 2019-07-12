source develop.env

function cleanup() {
    rm -f feedbridge
}
trap cleanup EXIT

# Compile Go
GO111MODULE=on GOGC=off go build -mod=vendor -v -o feedbridge ./cmd/api/
./feedbridge
