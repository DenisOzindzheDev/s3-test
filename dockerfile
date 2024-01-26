FROM golang:1.19


ENV CGO_ENABLED="0"
ENV GO111MODULE="on"
ENV GOFLAGS=-mod=vendor

LABEL Author = "denis.ozindzhe@tages.ru"

WORKDIR /app
COPY . . 

RUN go mod download
RUN go build -o s3-test cmd/main.go

EXPOSE 8080

CMD ["s3-test"]