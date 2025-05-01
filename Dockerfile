FROM golang:1.23.4-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ls -la /app

WORKDIR /app/cmd

RUN go build -o /app/app .

EXPOSE 8181

CMD ["/app/app"]
