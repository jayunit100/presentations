kubectl get pods --all-namespaces -o jsonpath=" {.items[*].spec.containers[*].image } "| tr ' ' "\n"


cat << EOF 
### Example output 

quay.io/coreos/clair
postgres:latest
k8s.gcr.io/coredns:1.2.6
k8s.gcr.io/coredns:1.2.6
k8s.gcr.io/etcd:3.2.24
k8s.gcr.io/kube-addon-manager:v8.6
k8s.gcr.io/kube-apiserver:v1.13.3
k8s.gcr.io/kube-controller-manager:v1.13.3
k8s.gcr.io/kube-proxy:v1.13.3
k8s.gcr.io/kube-scheduler:v1.13.3
gcr.io/k8s-minikube/storage-provisioner:v1.8.1
gcr.io/kubernetes-helm/tiller:v2.13.0
EOF
