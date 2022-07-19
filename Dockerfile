#
# MIT License
#
# (C) Copyright 2018-2022 Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.
#
# Dockerfile for building HMS bss.

### build-base stage ###
# Build base just has the packages installed we need.
FROM artifactory.algol60.net/docker.io/library/golang:1.18-alpine AS build-base

RUN set -ex \
    && apk -U upgrade \
    && apk add build-base

### base stage ###
# Base copies in the files we need to test/build.
FROM build-base AS base

RUN go env -w GO111MODULE=auto

# Copy all the necessary files to the image.
COPY vendor $GOPATH/src/github.com/Cray-HPE/cray-nls/vendor

COPY api $GOPATH/src/github.com/Cray-HPE/cray-nls/api
COPY bootstrap $GOPATH/src/github.com/Cray-HPE/cray-nls/bootstrap
COPY docs $GOPATH/src/github.com/Cray-HPE/cray-nls/docs
COPY utils $GOPATH/src/github.com/Cray-HPE/cray-nls/utils
COPY main.go $GOPATH/src/github.com/Cray-HPE/cray-nls/main.go
COPY .version $GOPATH/src/github.com/Cray-HPE/hms-bss/.version

### Build Stage ###
FROM base AS builder

RUN set -ex && go build -v -i -o /usr/local/bin/ncn-lifecycle-service github.com/Cray-HPE/cray-nls

### Final Stage ###
FROM gcr.io/distroless/static
LABEL maintainer="Hewlett Packard Enterprise"
EXPOSE 5000
STOPSIGNAL SIGTERM

# Get the boot-script-service from the builder stage.
COPY --from=builder /usr/local/bin/ncn-lifecycle-service /ncn-lifecycle-service

COPY .version /
COPY .env.example .env
# noboby 65534:65534
USER 65534:65534
# Setup environment variables.
ENV ENV=production
# Set up the command to start the service, the run the init script.
CMD /ncn-lifecycle-service
