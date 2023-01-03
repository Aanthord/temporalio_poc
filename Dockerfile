#FROM golang:alpine AS builder
FROM golang:alpine
# Install git.
# Git is required for fetching the dependencies.
#
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/aanthord/temporalio_poc/create_wallet/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
#RUN go build -o /go/bin/hello
############################
# STEP 2 build a small image
############################
#FROM scratch
# Copy our static executable.
#COPY --from=builder /go/bin/hello /go/bin/hello
# Run the hello binary.
ENTRYPOINT ["go run create_wallet/worker/main.go && go run create_wallet/starter/main.go"]