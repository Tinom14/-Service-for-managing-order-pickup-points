FROM golang:1.23.6-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --no-cache make

RUN go build -o server ./main.go

FROM alpine AS runner

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=build /build/server ./server
COPY config/config.yml ./config/config.yml
COPY --from=build /build/pkg/postgres_connect/migrations /app/pkg/postgres_connect/migrations

CMD ["/app/server", "--config=./config/config.yml"]