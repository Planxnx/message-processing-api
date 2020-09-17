FROM 1.15.2-alpine AS builder
WORKDIR /src/build 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:3.12.0
COPY --from=builder /src/build/main /src/app/
WORKDIR /src/app
CMD ["./main"]
