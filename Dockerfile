FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY configuration ./configuration
COPY database ./database
COPY operations ./operations
COPY public ./public
COPY main.go .
COPY application.yaml .

RUN go mod download
RUN go mod vendor
RUN go build -o /app main.go

EXPOSE 80
CMD ["./main"]

# RUN CGO_ENABLED=0 go build -o /app/bin/go-lunch main.go
