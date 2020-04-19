FROM golang as builder
COPY go.mod go.sum main.go ./
RUN export GOPATH="" && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /goserver .

RUN mkdir /cert && mkdir /static
RUN useradd scratchuser
RUN chown -R scratchuser:scratchuser /cert && chown -R scratchuser:scratchuser /static

FROM scratch
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /goserver /
COPY --from=builder /cert /cert
COPY --from=builder /static /static
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/passwd /etc/passwd
USER scratchuser
CMD ["/goserver"]