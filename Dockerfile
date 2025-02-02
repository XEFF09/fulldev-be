FROM golang:1.23.5

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8088

CMD ["air", "-c", "/app/.air.toml"]