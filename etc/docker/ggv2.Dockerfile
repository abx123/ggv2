# Name this as builder stage
FROM golang:1.16-alpine AS builder
# Move to working directory /getground
WORKDIR /go/src/ggv2
# Copy the code into the container
COPY . .
# ...
RUN go build -o server .

# Runtime stage
FROM golang:1.16-alpine
WORKDIR /app
ENV PORT=1323
ENV DSN=getground:password@tcp(mysql.c8ajbiky1mzj.ap-southeast-1.rds.amazonaws.com:3306)/getground
# Copy binary from builder stage
COPY --from=builder /go/src/ggv2/server .
EXPOSE 1323
# RUN mkdir logs
CMD ["./server"]