# Build stage
FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /app

# Install curl and tar in a single layer
RUN apk add --no-cache curl tar

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# this will basically copy the rest of the code
COPY . .

RUN go build -o main main.go
RUN apk add curl
# Download golang-migrate binary
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz \
    | tar xvz 


# Run stage
FROM alpine:3.22
WORKDIR /app

# copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./

# Copy env file and migration files
COPY app.env .
COPY compose/start.sh .
COPY compose/wait-for.sh .


RUN chmod +x /app/start.sh /app/wait-for.sh

COPY db/migration ./migration

# Exposed port
EXPOSE 8080

# run the app 
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]


