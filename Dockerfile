FROM golang:alpine

COPY ./ /go/src/github.com/rekkusu/gyotaku

WORKDIR /go/src/github.com/rekkusu/gyotaku
RUN apk --update --no-cache add git curl \
  && go get ./... \
  && apk del --purge git

EXPOSE 9999

CMD ["go", "run", "main.go"]
