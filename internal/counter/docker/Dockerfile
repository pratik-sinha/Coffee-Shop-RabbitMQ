FROM golang:1.21.3 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download



FROM golang:1.21.3 as builder
COPY --from=modules /go/pkg /go/pkg
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/main ./cmd/counter/main.go





FROM alpine:3.17 

WORKDIR /app

COPY --from=builder /bin/main /app/main
RUN chmod a+x /app/main 

COPY --from=builder /app/config /app/config/.

EXPOSE 5001

CMD [ "/app/main" ]

