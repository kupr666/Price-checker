FROM golang:1.26-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["make", "service-run"]