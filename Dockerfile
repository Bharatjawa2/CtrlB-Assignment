FROM golang:1.24.3-alpine

# Install dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Enable CGO
ENV CGO_ENABLED=1 \
    GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p storage
RUN go build -o server ./cmd/CTRLB

EXPOSE 8082
CMD ["./server"]
