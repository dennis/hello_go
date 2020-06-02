FROM golang:1.14.3-alpine as builder

RUN mkdir -p /src
WORKDIR /src
COPY . .
RUN go build -o /main main.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /main /app/
WORKDIR /app
EXPOSE 8080
CMD ["./main"]
