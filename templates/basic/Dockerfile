ARG GO_VERSION

FROM golang:${GO_VERSION}-alpine3.17 as builder
RUN apk --no-cache add make bash
WORKDIR /app  
COPY . /app
RUN make service

FROM alpine:3.17.2
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/cmd/{{.Name}}/{{.Name}} /usr/bin/{{.Name}}
CMD ["{{.Name}}"]
