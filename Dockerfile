FROM golang:alpine as builder
RUN apk add --no-cache git make gcc musl-dev
WORKDIR /work
COPY . .
RUN go mod download
RUN make build

FROM alpine:latest
WORKDIR /
COPY --from=builder /work/bin/sql-executor .
CMD [ "/sql-executor" ]
