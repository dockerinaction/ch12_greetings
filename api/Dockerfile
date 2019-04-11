FROM golang as app
WORKDIR /go/src/app
COPY ./api/main.go .
RUN go get ./...
RUN CGO_ENABLED=0 go build

FROM scratch
COPY ./api/config/config.common.yml .
ENTRYPOINT ["/bin/service"]
COPY --from=app /go/src/app/app /bin/service
