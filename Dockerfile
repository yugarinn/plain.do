FROM golang:1.22

ENV GOPATH=/go

RUN mkdir /plain.do
WORKDIR /plain.do

ADD go.mod ./go.mod
ADD . .

RUN go mod download && go mod verify

EXPOSE 8080

ENTRYPOINT ./setup.sh
