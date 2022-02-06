# syntax=docker/dockerfile:1

FROM golang:alpine
WORKDIR /app

# copy go module files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy and build go source code
COPY *.go ./
COPY controllers ./controllers
RUN go build -o blind-chess

# copy react build
COPY frontend/blind-chess/build ./frontend/blind-chess/build

ENV PORT 80
EXPOSE 80
CMD [ "./blind-chess" ]