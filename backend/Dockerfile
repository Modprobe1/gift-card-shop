FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd
RUN go mod tidy && go build -o /app/app .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
