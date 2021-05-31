FROM golang:latest 

RUN mkdir /build
WORKDIR /build

COPY go.mod .

COPY go.sum . 

RUN go mod download 

COPY . .

RUN  export GO111MODULE=on
RUN go get github.com/Oussemabhouri00/golang-api/master
RUN cd /build && git clone https://github.com/Oussemabhouri00/golang-api.git

RUN cd/buid/temp && go build

EXPOSE 10000

RUN go build 

ENTRYPOINT ["/build/temp"]
