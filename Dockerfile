FROM golang:1.10.3-alpine3.7 AS compile

# add git
RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

COPY . /go/src/github.com/pandorasnox/kubernetes-deny-env-var/
WORKDIR /go/src/github.com/pandorasnox/kubernetes-deny-env-var/
RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo .

# ============================================================
# ============================================================

FROM alpine:3.8
COPY --from=compile /go/bin/kubernetes-deny-env-var /kubernetes-deny-env-var
ENTRYPOINT ["/kubernetes-deny-env-var"]
