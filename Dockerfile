FROM golang:latest AS build
RUN mkdir -p /go/src/github.com/yanzay/picasso
ADD . /go/src/github.com/yanzay/picasso
RUN cd /go/src/github.com/yanzay/picasso && CGO_ENABLED=0 go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
EXPOSE 8014
COPY --from=build /go/src/github.com/yanzay/picasso/picasso /
ENTRYPOINT ["/picasso"]
CMD ["-m 60"]
