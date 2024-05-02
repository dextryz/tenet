# --------------------------------------------------------------------

FROM golang:1.22-alpine3.19 as builder

RUN apk add --no-cache gcc musl-dev sqlite-dev ca-certificates

# Set a temporary work directory
WORKDIR /app

# Add necessary go files
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Generate Go codes from template files
RUN go run -mod=mod github.com/a-h/templ/cmd/templ@latest generate

RUN CGO_ENABLED=1 GOOS=linux go build -tags 'osusergo netgo static_build' -ldflags '-extldflags "-static"' -o server ./cmd/server/main.go

# --------------------------------------------------------------------

FROM ubuntu:latest

COPY --from=builder /app/server ./
COPY --from=builder /app/static ./static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["/server"]

# --------------------------------------------------------------------
