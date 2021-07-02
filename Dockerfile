FROM golang:1.15.5-alpine3.12 as builder

RUN apk add make
WORKDIR /go/src/github.com/cycloidio/cycloid-resource
COPY . ./
RUN make

FROM alpine:3.12
COPY --from=builder /go/src/github.com/cycloidio/cycloid-resource/resource/ /opt/resource

RUN set -e; \
	apk add --no-cache --virtual .build-deps \
		curl \
	; \
	curl https://raw.githubusercontent.com/cycloidio/cycloid-cli/master/scripts/cy-wrapper.sh > /usr/bin/cy \
	&& chmod +x /usr/bin/cy; \
    apk del .build-deps;

# runtime dependencies
RUN apk add \
	bash \
	jq \
	curl \
	wget
