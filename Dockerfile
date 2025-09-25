FROM golang:1.20.1-buster
WORKDIR /app
COPY . .
RUN apt-get update && apt-get upgrade -y && apt-get -y install sqlite3
CMD go build -buildvcs=false -o ./bin/entities ./cmd/entities && ./bin/entities
