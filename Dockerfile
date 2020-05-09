FROM golang:alpine

RUN apk add --update alpine-sdk
RUN go get github.com/anacrolix/torrent
RUN go get github.com/anacrolix/utp
RUN go get github.com/gizak/termui/v3

WORKDIR /go/src/app
COPY . .

RUN mkdir -p /go/src/github.com/jeffchannell
RUN ln -s /go/src/app /go/src/github.com/jeffchannell/lbdl
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["make"]