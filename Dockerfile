FROM golang:1.10-alpine as builder
WORKDIR $GOPATH/src/github.com/dewey/feedbridge
ADD ./ $GOPATH/src/github.com/dewey/feedbridge
RUN apk update && \
    apk upgrade && \
    apk add git
RUN go get -u github.com/gobuffalo/packr/... && \
    cd $GOPATH/src/github.com/dewey/feedbridge && \    
    packr && \
    go build -v -o /feedbridge
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /feedbridge /feedbridge
CMD ["/feedbridge"]