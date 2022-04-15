FROM golang:1.16-alpine AS build_go

ARG DATE
ARG COMMIT
ARG VERSION
ARG OS=$(go env GOOS)
ARG ARCH=$(go env GOARCH)
WORKDIR /tmp/crud-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v ./...

# Build the Go app
RUN CGO_ENABLED=0 go build -ldflags \
	"-X github.com/srimaln91/go-make.version=$VERSION \
	-X github.com/srimaln91/go-make.date=$DATE \
	-X github.com/srimaln91/go-make.gitCommit=$COMMIT \
	-X github.com/srimaln91/go-make.osArch=$OS/$ARCH" -o crud-app-linux-amd64 cmd/api/main.go

FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_go /tmp/crud-app/crud-app-linux-amd64 /app/
COPY config.yaml /app/config/

WORKDIR /app

EXPOSE 8080

CMD ["/app/crud-app-linux-amd64"]