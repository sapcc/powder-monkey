FROM golang:1.15-alpine AS builder
RUN apk add --no-cache make gcc

COPY . /src
RUN make -C /src

################################################################################

FROM alpine:3.12
RUN apk add --no-cache ca-certificates curl
LABEL source_repository="https://github.com/sapcc/powder-monkey"

COPY --from=builder /src/bin/linux/amd64/powder-monkey /usr/bin/powder-monkey
ENTRYPOINT ["/usr/bin/powder-monkey"]
