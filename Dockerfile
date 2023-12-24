FROM golang:1.21-alpine3.18 as builder
RUN apk update && apk add git

WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o crane_ws

FROM alpine:3.18.5 as runner
COPY --from=builder /go/src/crane_ws /app/
ENTRYPOINT ["/app/crane_ws"]
