FROM golang:latest AS build

WORKDIR /go/src/generatedashboard

COPY . .

# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 go install ./...

# get certificates, will be copied below
FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates

# actual Docker image
FROM scratch

WORKDIR /

COPY --from=build /go/bin/generatedashboard generatedashboard
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["./generatedashboard"]
