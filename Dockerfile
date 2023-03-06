FROM golang:1.19-alpine3.17 AS build

COPY . /src
WORKDIR /src
RUN go build


FROM alpine:3.17

COPY --from=build /src/k8s-secret-admin /k8s-secret-admin
ENTRYPOINT ["/k8s-secret-admin"]
