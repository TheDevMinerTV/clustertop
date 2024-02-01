FROM golang:1.21 AS builder
WORKDIR /src

COPY ./go.sum ./go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /clustertop -ldflags="-w -s"

FROM alpine:3.19 AS runner
RUN chmod +x /entrypoint.sh
EXPOSE 80

COPY --from=builder /clustertop /clustertop

ENTRYPOINT ["/entrypoint.sh"]