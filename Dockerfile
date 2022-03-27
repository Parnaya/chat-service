FROM golang:1.16-alpine

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o /app main.go

EXPOSE 80
CMD ["./main"]

# RUN CGO_ENABLED=0 go build -o /app/bin/go-lunch main.go
