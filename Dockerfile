FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc musl-dev

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . ./

RUN go build -o ./bin/app cmd/main.go

FROM alpine 

COPY --from=builder /usr/local/src/.env /
COPY --from=builder /usr/local/src/migrations /migrations
COPY --from=builder /usr/local/src/docs /docs
COPY --from=builder /usr/local/src/bin/app /
CMD ["/app"]