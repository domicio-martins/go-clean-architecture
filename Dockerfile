FROM golang:alpine

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ARG APP_STAGE
ENV APP_STAGE $APP_STAGE

RUN apk --no-cache add curl

RUN apk add -U --no-cache ca-certificates git make

WORKDIR /app

ADD go.mod go.sum ./

RUN go mod download -x

ADD . .

RUN make

ENTRYPOINT ["/bin/sh", "/app/docker-entrypoint.sh"]