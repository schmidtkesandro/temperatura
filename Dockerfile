FROM golang:1.21 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tempcloudrun

#FROM scratch
FROM alpine:latest  
WORKDIR /app
COPY --from=build /app/tempcloudrun .
EXPOSE 8080
ENTRYPOINT ["./tempcloudrun"]
