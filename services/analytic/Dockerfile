# First stage: build the application
FROM golang:1.22.3-alpine AS builder

WORKDIR /tenk-backend

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o analytic ./services/analytic

# Second stage: build the final image
FROM alpine:latest

WORKDIR /tenk-backend

COPY --from=builder /tenk-backend/analytic analytic
COPY --from=builder /tenk-backend/pkg pkg
COPY --from=builder /tenk-backend/proto proto

EXPOSE 8082
EXPOSE 50052

CMD ["./analytic"]
