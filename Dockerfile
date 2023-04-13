FROM golang:1.20.3-alpine3.17
WORKDIR /app
ADD . /app
RUN cd /app
RUN go mod download
CMD ["go", "run", "app/main.go"]