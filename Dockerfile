
FROM golang:1.14 as builder

RUN git clone https://github.com/haibeey/toptweet
WORKDIR ./toptweet/src
RUN go test
RUN go build -o toptweet

EXPOSE 8080
CMD ["./toptweet"]