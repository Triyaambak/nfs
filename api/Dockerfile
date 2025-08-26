FROM golang:alpine

WORKDIR /app/api

COPY ./go.mod ./go.sum /

RUN go mod download

COPY ./ /

RUN go build ./main.go

CMD ["sh","./main"]
