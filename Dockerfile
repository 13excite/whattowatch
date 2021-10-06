FROM golang:alpine as builder

RUN apk update && apk add curl \
                          bash \
                          make  && \
     rm -rf /var/cache/apk/*

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy and download dependency using go mod
COPY . .
RUN  make deps

RUN make build

WORKDIR /dist

RUN cp /build/whattowatchcmd .

# Build a small image
FROM scratch

COPY --from=builder /dist/ /

CMD ["/whattowatchcmd", "api"]

