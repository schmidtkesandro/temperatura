FROM golang:latest as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tempcloudrun


FROM scratch
WORKDIR /app
COPY --from=build /app/tempcloudrun .
ENTRYPOINT ["./tempcloudrun"]
