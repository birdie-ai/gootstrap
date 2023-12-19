ARG GO_VERSION
ARG ALPINE_VERSION

FROM docker.io/golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder
RUN apk --no-cache add make bash
WORKDIR /app  
COPY . /app
RUN make service

FROM docker.io/alpine:${ALPINE_VERSION}
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/cmd/{{.Name}}/{{.Name}} /usr/bin/{{.Name}}
CMD ["{{.Name}}"]
