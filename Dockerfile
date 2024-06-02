FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod go.sum .env prometheus.yml ./
COPY backend ./backend
COPY frontend ./frontend
COPY cmd ./cmd
RUN go mod tidy
RUN go build -o /setup ./cmd/setup
RUN go build -o /server ./cmd/server

CMD ["/server"]
