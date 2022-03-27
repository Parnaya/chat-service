FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go get .
RUN go mod download

COPY configuration ./configuration
COPY database ./database
COPY operations ./operations
COPY public ./public
COPY main.go .
COPY application.yaml .

RUN go build -o /app main.go

EXPOSE 80
CMD ["./main"]

# RUN CGO_ENABLED=0 go build -o /app/bin/go-lunch main.go
