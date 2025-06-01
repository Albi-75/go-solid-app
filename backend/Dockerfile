FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o my-app ./cmd/app

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/my-app .
COPY db/migrations ./db/migrations
EXPOSE 8080
CMD ["./my-app"]