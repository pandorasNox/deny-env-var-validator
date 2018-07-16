
# kubernetes-deny-env-var

Kubernetes admission controller validator webhook

### generate cert.pem and key.pem for ListenAndServeTLS
`docker run --rm -v $(pwd)/certs:/certs -e SSL_SUBJECT=test.example.com -e SSL_KEY="ssl-key.pem" -e SSL_CSR="ssl-key.csr" -e SSL_CERT="ssl-cert.pem" -e K8S_NAME="pls-replace-me-kubernetes-name" paulczar/omgwtfssl`

### test local
without env
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{},"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}' \
  http://localhost:8083/content
```

with env
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","env":[{"name":"foo","value":"bar"}],"ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{},"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}' \
  http://localhost:8083/content
```

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
