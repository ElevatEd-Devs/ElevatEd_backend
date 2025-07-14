FROM golang:1.24.5-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.* ./

RUN go mod download

COPY . .

EXPOSE 3000

CMD [ "air", "-c", ".air.toml" ]