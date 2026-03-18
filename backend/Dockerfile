FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY bin/api ./api
COPY migrations ./migrations

EXPOSE 8080

CMD ["./api"]
