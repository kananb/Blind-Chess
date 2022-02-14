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

# copy and build react code
FROM node:alpine
COPY client ./client
RUN npm --prefix client install
RUN npm --prefix client run build

CMD [ "./blind-chess" ]