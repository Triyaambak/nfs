FROM golang:alpine

ARG API_PORT

WORKDIR /app/api

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build ./main.go

EXPOSE ${API_PORT}

CMD ["./main"]
