FROM golang:latest as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tempcloudrun
RUN cp /etc/ssl/certs/ca-certificates.crt /app/ca-certificates.crt

FROM scratch
WORKDIR /app
COPY --from=build /app/tempcloudrun .
COPY --from=build /app/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./tempcloudrun"]
