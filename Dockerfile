FROM golang:1.15-alpine as build
#TODO MULTISTEP WITH TEST
WORKDIR /app/

COPY api api

COPY pkg pkg

COPY resources resources

COPY scripts/sql scripts/sql

COPY go.mod go.mod

COPY go.sum go.sum

COPY main.go main.go

RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth1

RUN ls

FROM alpine:latest
WORKDIR /app/
EXPOSE 80

COPY --from=build /app/auth1 .

COPY resources resources

COPY pkg/routes/front/internal/templates pkg/routes/front/internal/templates
CMD ["/app/auth1"]
