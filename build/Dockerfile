FROM golang:1.25.6-alpine AS builder

COPY go.mod go.sum /github.com/totorialman/go-test-ac/
WORKDIR /github.com/totorialman/go-test-ac

RUN go mod download
COPY . /github.com/totorialman/go-test-ac

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main.go

FROM scratch AS runner

WORKDIR /build_v1/

COPY --from=builder /github.com/totorialman/go-test-ac/.bin .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8080

ENTRYPOINT ["./.bin"]