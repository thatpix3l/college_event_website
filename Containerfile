FROM golang:1.21-alpine as go-build

WORKDIR /workdir

COPY . .

RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -o /app .

FROM scratch

COPY --from=go-build /app .

CMD ["/app"]