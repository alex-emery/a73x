FROM golang:1.22 as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /bin/site .

FROM scratch 
COPY --from=builder /bin/site /app

ENTRYPOINT ["/app"]

