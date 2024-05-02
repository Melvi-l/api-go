FROM golang:1.22.2-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o myapp .

FROM alpine
COPY --from=builder /app/myapp /myapp
EXPOSE 8080
CMD ["/myapp"]
