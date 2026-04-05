FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-marketer ./...

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/ai-marketer .
COPY --from=builder /app/prompts/ ./prompts/
COPY --from=builder /app/config.yaml .
ENTRYPOINT ["./ai-marketer"]
CMD ["run"]
