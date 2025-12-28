# Stage 0: Cache Go modules
FROM golang:1.25.1-alpine AS go-modules

WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache ca-certificates git && \
  CGO_ENABLED=0 go mod download


# Stage 1: Build static assets (CSS/JS) with Node.js
FROM node:22-slim AS js-builder

WORKDIR /app
COPY web/package.json ./web/
# RUN npm --prefix web ci --omit=dev
RUN npm --prefix web install --omit=dev

COPY web/ ./web/
COPY internal/view/ ./internal/view/
RUN npm --prefix web run build:css && \
  npm --prefix web run build:js -- --minify


# Stage 2: Generate templ files and build Go binary
FROM golang:1.25.1-alpine AS go-builder
# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy Go modules cache
COPY --from=go-modules /go/pkg /go/pkg
COPY go.mod go.sum ./

# Copy source (without built assets yet)
COPY . .

# Copy prebuild assets
COPY --from=js-builder /app/web/public/assets /web/public/assets/

# Generate templ Go files
RUN templ generate

# Build Go binary (CGO disabled for scratch compatibility)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o /app/bin/app ./cmd/app

# Stage 3: Final minimal image
FROM scratch

# Copy binary
COPY --from=go-builder /app/bin/app /app

# Copy CA certificates (for HTTPS requests)
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Optional: copy static assets if served by Go app
# (adjust path based on your app's expectations)
COPY --from=js-builder /app/web/public/assets /web/public/assets/

CMD ["/app"]
