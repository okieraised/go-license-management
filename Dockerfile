FROM golang:1.23.4-alpine3.21 AS build

RUN mkdir /opt/app
WORKDIR /opt/app

COPY ./*.go ./
COPY ./docs ./docs
COPY ./internal ./internal
COPY ./server ./server

COPY go.mod go.sum ./
RUN go mod download

RUN go install golang.org/x/vuln/cmd/govulncheck@latest
RUN govulncheck ./...

RUN go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
RUN fieldalignment -json -fix ./...

WORKDIR /opt/app

RUN go build -o go-license-management .

FROM alpine:latest

RUN mkdir /opt/app
WORKDIR  /opt/app

RUN apk add tzdata
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=build /opt/app/go-license-management /opt/app/go-license-management
RUN mkdir ./conf
RUN touch ./conf/config.toml

CMD ["./go-license-management"]
