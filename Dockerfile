FROM golang:1.16 as builder
COPY go.mod go.sum cmd/main.go ./
RUN export GOPATH="" && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /goserver .

RUN mkdir /cert && mkdir /static
RUN useradd scratchuser
RUN chown -R scratchuser:scratchuser /cert && chown -R scratchuser:scratchuser /static

FROM alpine
COPY --from=builder /goserver /
COPY --from=builder /etc/passwd /etc/passwd
USER scratchuser
CMD ["/goserver"]