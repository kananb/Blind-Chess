# Dockerfile

FROM node:17-alpine AS builder
WORKDIR /app

# copy and build react code
COPY client ./
RUN npm install --production --no-audit
RUN npm run build

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

# copy built react code
COPY --from=builder /app/build ./client/build

CMD [ "./blind-chess" ]