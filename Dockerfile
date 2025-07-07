FROM golang:1.23

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o go-insight cmd/main.go

EXPOSE 8001
CMD ["./go-insight"]