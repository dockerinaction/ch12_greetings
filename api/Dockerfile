FROM golang:1.12 as app
WORKDIR /go/src/app
COPY ./api/*.go /go/src/app/
RUN go get ./...
RUN CGO_ENABLED=0 go build

# Using alpine as a base image because it makes it easy to inspect files inside the container
# `scratch` functions just fine
FROM alpine:3.9
COPY ./api/config/config.common.yml /config/
RUN chown -R 1000:1000 /config && chmod 0400 /config/*
ENTRYPOINT ["/bin/service"]
COPY --from=app /go/src/app/app /bin/service
