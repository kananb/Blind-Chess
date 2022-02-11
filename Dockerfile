# syntax=docker/dockerfile:1

FROM golang:alpine
WORKDIR /app

# copy go module files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy and build go source code
COPY *.go ./
COPY game ./game
COPY data ./data
RUN go build -o blind-chess

# copy react build
COPY client/build ./client/build

ENV PORT 80
EXPOSE 80
CMD [ "./blind-chess" ]