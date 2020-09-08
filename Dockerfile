FROM golang:1.12-alpine

RUN go get github.com/InSitu-Software/extrapolator/bin && \
	cp /go/bin/bin /bin/extrapolator && \
	rm -rf /go/src/*
