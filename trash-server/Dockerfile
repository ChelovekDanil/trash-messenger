FROM golang:1.22.7 AS builder
WORKDIR /src
COPY . .
RUN go get ./... && CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM scratch
WORKDIR /app
COPY --from=builder /src/server .
EXPOSE 8080
CMD [ "./server" ]