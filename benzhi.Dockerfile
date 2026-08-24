FROM golang:1.26.2
WORKDIR /src
ENV GOPROXY=off GOSUMDB=off
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN go build -mod=vendor -o /usr/local/bin/fireline ./cmd/fireline
CMD ["/usr/local/bin/fireline"]
