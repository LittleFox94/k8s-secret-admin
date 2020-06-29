FROM golang:alpine3.12 AS build

COPY . /src
WORKDIR /src
RUN go build


FROM alpine:3.12

COPY --from=build /src/k8s-secret-admin /k8s-secret-admin
ENTRYPOINT ["/k8s-secret-admin"]
