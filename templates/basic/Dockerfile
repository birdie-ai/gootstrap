ARG ALPINE_VERSION

FROM docker.io/alpine:${ALPINE_VERSION}
RUN apk --no-cache add ca-certificates
COPY ./cmd/{{.Name}}/{{.Name}} /app/{{.Name}}
CMD ["/app/{{.Name}}"]
