FROM golang AS builder

RUN mkdir /builder
COPY go.mod /builder
COPY middleware.go /builder
RUN cd /builder && go build -o middleware
RUN chmod +x /usr/sbin

FROM ubuntu

COPY --from=builder /builder/middleware /usr/sbin

ENTRYPOINT [ "middleware" ]