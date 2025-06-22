# First stage: build the application
FROM golang:1.22.3-alpine AS builder

WORKDIR /tenk-backend

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o notification ./services/notification

# Second stage: build the final image
FROM alpine:latest

WORKDIR /tenk-backend

COPY --from=builder /tenk-backend/notification notification
COPY --from=builder /tenk-backend/pkg pkg
COPY --from=builder /tenk-backend/proto proto

EXPOSE 8084
EXPOSE 50054

CMD ["./notification"]
