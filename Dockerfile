FROM golang:1.12-alpine

COPY bin/extrapolator.go /go/bin/

RUN rm -rf /go/src/*
