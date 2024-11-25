FROM golang:1.19 AS build_go

ARG DATE
ARG COMMIT
ARG VERSION
ARG OS="linux"
ARG ARCH="amd64"
WORKDIR /tmp/crud-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -ldflags \
	"-X github.com/srimaln91/go-make.version=$VERSION \
	-X github.com/srimaln91/go-make.date=$DATE \
	-X github.com/srimaln91/go-make.gitCommit=$COMMIT \
	-X github.com/srimaln91/go-make.osArch=$OS/$ARCH" -o crud-app-linux-amd64 cmd/api/main.go


FROM alpine:3.20
RUN apk add ca-certificates

COPY --from=build_go /tmp/crud-app/crud-app-linux-amd64 /app/
# COPY config.yaml /app/

WORKDIR /app

EXPOSE 8080

RUN pwd
RUN ls -al
CMD ["/app/crud-app-linux-amd64"]