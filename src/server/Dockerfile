FROM golang:alpine

RUN mkdir /app
WORKDIR /app
COPY . /app/
RUN go install
CMD ["/go/bin/server"]
