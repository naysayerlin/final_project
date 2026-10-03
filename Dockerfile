FROM golang:1.26-alpine AS builder

WORKDIR /usr/src/app/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /usr/src/app/scheduler .
COPY --from=builder /usr/src/app/web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD="12345"

EXPOSE 7540

CMD ["./scheduler"]