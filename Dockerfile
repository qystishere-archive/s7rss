FROM golang:1.15-alpine

RUN apk add --no-cache --update \
	git

# Set enviroment variables.
ENV WORKSPACE=github.com/qystishere/s7rss \
    GOPATH=/go \
	PATH="/go/bin:$PATH" \
	CGO_ENABLED=0

# Copy the local package files to the container's workspace. Add to GOPATH.
ADD . /go/src/${WORKSPACE}

# Set workdir.
WORKDIR /go/src/${WORKSPACE}

# Build.
RUN go install ${WORKSPACE}/cmd/s7rss

# Run the compiled bin by default when the container start.
CMD /go/bin/s7rss

# Service listens on port 80.
EXPOSE 9000
