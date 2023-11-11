# Build App
FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /go/src
ENV GOPROXY=https://goproxy.io,direct
COPY . .
RUN go mod tidy -x
RUN GOAMD64=v2 CGO_ENABLED=1 go build -ldflags="-s -w" -o main .

# Reduce Size Without Golang Image
FROM cgr.dev/chainguard/wolfi-base:latest AS intermediate
RUN adduser capstone -D && mkdir -p /home/capstone/bin && chown -Rf capstone:capstone /home/capstone/bin
USER capstone
WORKDIR /home/capstone/bin
COPY --from=builder /go/src/main .

FROM scratch AS runner
COPY --from=intermediate / /
CMD ["/home/capstone/bin/main", "server"]
