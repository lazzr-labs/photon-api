FROM golang:1.27.0

WORKDIR /app

COPY api-v1/ ./

RUN curl -sSf https://atlasgo.sh | sh

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-go
