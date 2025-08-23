# thape - casting container images to gzipped tarballs.
# Copyright(c) 2025 Star Inc. and its contributors.
# The software is licensed under BSD-3-Clause.

FROM golang:alpine AS builder
COPY . /factory
RUN apk add git make \
    && cd /factory \
    && make \
    && go clean -cache

FROM alpine:latest
ENV GIN_MODE release
COPY --from=builder /factory/LICENSE /app/LICENSE
COPY --from=builder /factory/.env.sample /app/.env
COPY --from=builder /factory/build/thape /app/thape
WORKDIR /app
ENTRYPOINT /app/thape
EXPOSE 6000
