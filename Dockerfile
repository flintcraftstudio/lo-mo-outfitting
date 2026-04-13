FROM golang:1.25-alpine AS build

WORKDIR /src

RUN apk add --no-cache curl

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install Tailwind standalone CLI
RUN curl -sL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 -o /usr/local/bin/tailwindcss \
    && chmod +x /usr/local/bin/tailwindcss

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate templ
RUN templ generate

# Build Tailwind CSS
RUN tailwindcss -c tailwind/tailwind.config.js -i tailwind/input.css -o web/static/css/site.css --minify

# Build binary
RUN CGO_ENABLED=0 go build -o /app ./cmd/server

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=build /app /app
COPY --from=build /src/web/static /web/static

ENV PORT=8080
EXPOSE 8080

CMD ["/app"]
