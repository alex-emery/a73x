FROM golang:1.22

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /usr/local/bin/app ./cmd/serve/

FROM scratch 

COPY --from=0 /usr/local/bin/app /app

CMD ["/app"]
