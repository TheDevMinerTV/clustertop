FROM golang:1.25 AS builder
WORKDIR /src

COPY ./go.sum ./go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /clustertop -ldflags="-w -s"

FROM alpine:3.23 AS runner
RUN adduser -D -u 1000 app
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 80

COPY --from=builder /clustertop /clustertop

ENTRYPOINT ["/entrypoint.sh"]