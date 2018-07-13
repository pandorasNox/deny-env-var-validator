
# kubernetes-deny-env-var

Kubernetes admission validator controller webhook

### generate cert.pem and key.pem for ListenAndServeTLS
`docker run --rm -v $(pwd)/certs:/certs -e SSL_SUBJECT=test.example.com -e SSL_KEY="ssl-key.pem" -e SSL_CSR="ssl-key.csr" -e SSL_CERT="ssl-cert.pem" -e K8S_NAME="pls-replace-me-kubernetes-name" paulczar/omgwtfssl`

### test in kubernetes
1. port-forward
`kubectl -n deny-env port-forward deny-env-88f744dc-tw8p6 9876:8083`

2. curl + ssl-cert + resolve hostname
`curl --cacert certs/ssl-cert.pem --resolve test.example.com:9876:127.0.0.1 https://test.example.com:9876/content`

#### resources
webhooks
- https://github.com/kubernetes/kubernetes/blob/v1.10.0-beta.1/test/images/webhook/main.go
- https://github.com/kubernetes/kubernetes/blob/master/test/e2e/apimachinery/webhook.go
- https://github.com/kelseyhightower/denyenv-validating-admission-webhook
- https://github.com/caesarxuchao/example-webhook-admission-controller
- https://de.slideshare.net/sttts?utm_campaign=profiletracking&utm_medium=sssite&utm_source=ssslideview
- https://github.com/openshift/generic-admission-server
- https://kubernetes.io/docs/tasks/access-kubernetes-api/http-proxy-access-api/
- https://github.com/kelseyhightower/denyenv-validating-admission-webhook/blob/master/index.js

initilizers
- https://github.com/kelseyhightower/kubernetes-initializer-tutorial/tree/master/envoy-initializer
- https://ahmet.im/blog/initializers/
- https://medium.com/ibm-cloud/kubernetes-initializers-deep-dive-and-tutorial-3bc416e4e13e
- https://groups.google.com/forum/#!topic/istio-users/lZxmROZxYKI
- https://groups.google.com/forum/?utm_medium=email&utm_source=footer#!msg/istio-dev/mIAbIRjCfZg/NKZfz9X8BgAJ
- https://istio.io/docs/setup/kubernetes/sidecar-injection/

cert
- https://github.com/kubernetes/kubernetes/issues/61171
