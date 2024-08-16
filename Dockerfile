FROM golang:1.22-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go install github.com/swaggo/swag/cmd/swag@latest
# Copy all files
COPY . .
# Init swagger
WORKDIR /app/app
RUN swag init
# Copy docs files
WORKDIR /app
COPY /app/docs /app/docs
# Build app
WORKDIR /app/app
RUN go build -o /app/main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]
EXPOSE 3000