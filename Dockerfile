FROM golang:1.23

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN mkdir /app
ADD . /app/
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go build -o main .
CMD ["/app/main"]