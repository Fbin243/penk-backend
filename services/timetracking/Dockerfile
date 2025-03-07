# First stage: build the application
FROM golang:1.22.3-alpine AS builder

WORKDIR /tenk-backend

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o timetracking ./services/timetracking

# Second stage: build the final image
FROM alpine:latest

WORKDIR /tenk-backend

COPY --from=builder /tenk-backend/timetracking timetracking
COPY --from=builder /tenk-backend/pkg pkg
COPY --from=builder /tenk-backend/proto proto

EXPOSE 8083
EXPOSE 50053

CMD ["./timetracking"]
