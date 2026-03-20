# ---- Build Stage ----
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o grade-api .

# ---- Run Stage ----
FROM alpine:3.19

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/grade-api .

USER appuser

EXPOSE 8080

CMD ["./grade-api"]