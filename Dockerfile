FROM golang:1.21 as builder
LABEL maintainer="Yui Qi Tang <yqtang1222@gmail.com>"
WORKDIR /app
COPY . .
RUN go mod download
RUN make build

FROM alpine:3.19
LABEL maintainer="Yui Qi Tang <yqtang1222@gmail.com>"
WORKDIR /app
COPY --from=builder /app/bin/app /app/glookbs
CMD ["./glookbs", "runserver"]