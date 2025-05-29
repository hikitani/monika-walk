FROM golang:alpine AS builder
WORKDIR /build
ENV CGO_ENABLED=0
RUN apk update --no-cache && apk add --no-cache tzdata
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -o sigma .

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates
WORKDIR /build
COPY --from=builder /build/ /build/
CMD ["./sigma"]
