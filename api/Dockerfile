FROM golang:1.5.1-onbuild

RUN mkdir /go/src/candy
ADD main.go /go/src/candy/main.go
RUN go install candy

ENTRYPOINT /go/bin/candy

EXPOSE 8080
