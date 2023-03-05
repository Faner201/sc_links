FROM golang:1.20

WORKDIR /usr/src/sc_links

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o sc_links ./cmd/main.go

CMD ["sc_links"]