FROM golang:1.20 as dependencies

FROM dependencies 

RUN apt-get update && apt-get install -y iputils-ping

ENV HOME /home

WORKDIR $HOME/athn

COPY pkg ./pkg
COPY cmd ./cmd
COPY go.mod \
  go.sum ./

RUN go build -o athnsrv ./cmd/athn/main.go 

EXPOSE 8080
EXPOSE 2345
EXPOSE 3456
EXPOSE 4567

CMD ["./athnsrv"]