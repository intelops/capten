FROM golang AS builder

WORKDIR /go/src/app

COPY . /go/src/app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o capten ./cmd/main.go

FROM alpine

WORKDIR /app


COPY --from=builder /go/src/app/capten /app/

CMD ["./capten"]
