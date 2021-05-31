FROM golang:latest 

WORKDIR /app

COPY go.mod .

COPY go.sum . 

RUN go mod download 

COPY . .

ENV port 8082

RUN go build 

ENTRYPOINT ["/bin/bash"]
