FROM golang:1.14-alpine3.12 AS build-env

#Build deps
RUN apk --no-cache add build-base git

#Setup 
COPY . /src/gitea-group-sync
WORKDIR /src/gitea-group-sync

RUN go get -d -v && go build

# Final
FROM alpine:3.12

COPY --from=build-env /src/gitea-group-sync/gitea-group-sync /app/gitea-group-sync/gitea-group-sync

RUN ln -s /app/gitea-group-sync/gitea-group-sync /usr/local/bin/gitea-group-sync

ENTRYPOINT ["/usr/local/bin/gitea-group-sync"]
