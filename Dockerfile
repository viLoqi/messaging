############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache 'git=~2'

# Install dependencies
ENV GO111MODULE=on
WORKDIR $GOPATH/src/packages/goginapp/

COPY go.mod .

# Fetch dependencies.
# Using go mod download.
RUN go mod download

COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/main .

############################
# STEP 2 build a small image
############################
FROM alpine:3

WORKDIR /

# Copy our static executable.
COPY --from=builder /go/main /go/main
# COPY public /go/public

ENV PORT 8080
ENV GIN_MODE release
EXPOSE 8080

WORKDIR /go

COPY credentials.json .

# Run the Go Gin binary.
ENTRYPOINT ["/go/main"]
