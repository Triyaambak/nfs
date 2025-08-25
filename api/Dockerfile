FROM golang:alpine

WORKDIR /app

COPY ./go.mod ./go.sum /app/

RUN go mod download

COPY ./ /app/

RUN go build ./main.go

CMD ["sh","./main"]
