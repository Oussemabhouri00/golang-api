# Start from golang base image
FROM golang:alpine 

#ENV GO111MODULE=on
RUN mkdir /app
RUN mkdir /app/temp

COPY . /app/temp

WORKDIR /app/temp


COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o main .

CMD ["/app/temp/main"]