FROM golang:1.13 as builder
WORKDIR /build
ADD . .
RUN cd cmd/qrcoded && CGO_ENABLED=0 go build

FROM alpine:3.10.2
WORKDIR /app
COPY --from=builder /build/cmd/qrcoded/qrcoded /app/qrcoded
EXPOSE 80
CMD ["./qrcoded"]