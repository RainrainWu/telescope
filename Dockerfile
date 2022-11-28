FROM golang:1.19.3-alpine as builder

COPY . /src
WORKDIR /src
RUN apk add --update --no-cache ca-certificates && \
    CGO_ENABLED=0 go build -o /src/bin/telescope .

FROM scratch

ENV PATH=/
COPY --from=builder /src/bin/ /
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

CMD [ "telescope", "-h" ]