FROM golang:1.24.4-alpine

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api

CMD ["/app/api"]