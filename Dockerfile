FROM golang:1.15.5-alpine3.12 as builder

RUN apk add make
WORKDIR /go/src/github.com/cycloidio/infrapolicy-resource
COPY . ./
RUN make

FROM cycloid/cycloid-toolkit:develop
COPY --from=builder /go/src/github.com/cycloidio/infrapolicy-resource/resource/ /opt/resource
