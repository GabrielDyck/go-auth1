FROM golang:1.15
#TODO MULTISTEP WITH TEST
WORKDIR /app/

COPY api api

COPY pkg pkg

COPY resources resources

COPY scripts/sql scripts/sql

COPY go.mod go.mod

COPY go.sum go.sum

COPY main.go main.go


RUN go build

EXPOSE 80

CMD ["/app/auth1"]
