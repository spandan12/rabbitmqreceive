FROM golang:alpine

WORKDIR /app

# Copy go source code and install dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

CMD ["go", "run", "main.go"]