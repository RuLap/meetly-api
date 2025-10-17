FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o meetly-api cmd/api/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

COPY --from=build /app/meetly-api .
COPY --from=build /go/bin/goose /usr/local/bin/goose
COPY internal/pkg/config/config.yaml ./config.yaml
COPY migrations/ ./migrations/

EXPOSE 8080

CMD ["./meetly-api"]
