FROM golang as app
WORKDIR /go/src/app
COPY ./api/main.go .
RUN go get ./...
RUN CGO_ENABLED=0 go build

FROM alpine:3.9
COPY ./api/config/config.common.yml /config/
RUN chown -R 1000:1000 /config && chmod 0400 /config/*
ENTRYPOINT ["/bin/service"]
COPY --from=app /go/src/app/app /bin/service
