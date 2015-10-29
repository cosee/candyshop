FROM golang:1.5.1-onbuild

RUN mkdir /go/src/candy-worker
ADD main.go /go/src/candy/main.go
RUN go install candy-worker

ENTRYPOINT /go/bin/candy

