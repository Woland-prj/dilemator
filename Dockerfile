# Step 1: Modules caching
FROM golang:1.24.6-alpine3.21 AS modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN apk add --no-cache git=2.47.3-r0 ca-certificates=20250911-r0

ARG GITHUB_TOKEN
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"
RUN go env -w GOPRIVATE=github.com/Giashka/* \
    && go mod download

# Step 2: Builder
FROM golang:1.24.6-alpine3.21 AS builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch

COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/docs /docs

CMD ["/app"]