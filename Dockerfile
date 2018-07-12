FROM golang:1.10.3-alpine3.7 AS compile

# add git
RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

COPY . /go/src/k8s-resource-mutator/
WORKDIR /go/src/k8s-resource-mutator/
RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo .

# ============================================================
# ============================================================

FROM alpine:3.8
COPY --from=compile /go/bin/k8s-resource-mutator /k8s-resource-mutator
ENTRYPOINT ["/k8s-resource-mutator"]
