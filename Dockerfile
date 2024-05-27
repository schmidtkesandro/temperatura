FROM golang:latest as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tempcloudrun

#FROM scratch
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/tempcloudrun .
ENTRYPOINT ["./tempcloudrun"]
