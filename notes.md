# Threat modelling in kubernetes cluster

## Whats a threat model?

```
Threat modeling is a process by which potential threats, such as structural vulnerabilities can be identified, enumerated, and prioritized – all from a hypothetical attacker’s point of view. 
```

```
The purpose of threat modeling is to provide defenders with a systematic analysis of the probable attacker’s profile, the most likely attack vectors, and the assets most desired by an attacker. 
```

- High level assets?
- Where am i vulnerable? 
- Most relevant threats? 
- Attack vector that is unnoticed? 
- The importance of performance.

## Who am I 

- Kubernetes core contributor @ Red Hat (when it was wild and crazy)
- Cloud Native Czar at Synopsys/Blackduck (security)
- Now Kubernetes core engineer again : Platform9

## Lets start threat modelling.

High level assets:

- ETCD ~ it has your entire cluster contents.
- Kubelets ~ they run as root, and they generally have a docker daemon on them.
- API Server ~ bringing an API server down could DoS any 'cloud native' app in a cluster.
- Apps (security assets  and vulnerabilities specific to your code)
- Serices (downtime)

In this talk: Focus on kubelets, API server, ETCD, Apps.

## How do I know I haven't broken my cluster when securing it? 

Run the Conformance tests!

- Sonubuoy: https://scanner.heptio.com/
- Kube-test manually: 

```
mkdir ~/bin
curl -sL -o ~/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x ~/bin/gimme
export PATH=$PATH:~/bin
export GOPATH=~/go
eval "$(gimme stable)"
KUBERNETES_DIR=kubernetes
KUBERNETES_VERSION=v1.12.5
git clone --branch ${KUBERNETES_VERSION} https://github.com/kubernetes/kubernetes.git
cd ${KUBERNETES_DIR} && make WHAT=test/e2e/e2e.test && make kubectl && make ginkgo
 go get -u k8s.io/test-infra/kubetest
export KUBE_MASTER_IP=ojas-deprecate-heapster-c8fb92c1-api.ojas-test.platform9.horse
export KUBE_MASTER=ip-10-0-1-125.us-west-2.compute.internal
export KUBECONFIG="~/kube/config.yml"
cd kubernetes/
$GOPATH/bin/kubetest --test --test_args="--ginkgo.focus=Secrets" --provider=skeleton
```

# APIServer: ClusterRoleBindings

Allow you to do anything.  Look at them, DELETE THEM !

`kubectl get clusterrolebindings`

Example: Minikube
```
cluster-admin                                          11d
kubeadm:kubelet-bootstrap                              11d
kubeadm:node-autoapprove-bootstrap                     11d
kubeadm:node-autoapprove-certificate-rotation          11d
kubeadm:node-proxier                                   11d
minikube-rbac                                          11d
storage-provisioner                                    11d
```

Example: In a managed kubernetes distribution


```
admin-and-blahblah-access                                   13d
apiproxy-access                                        13d
cluster-admin                                          13d
default-access                                         13d
heapster-binding                                       13d
kube-proxy-access                                      13d
kubelet-access                                         13d
```

From the docs: 

```
	API servers create a set of default ClusterRole and ClusterRoleBinding objects. Many of these are  system: prefixed, which indicates that the resource is "owned" by the infrastructure. Modifications to these resources can result in non-functional clusters. One example is the system:node ClusterRole. This role defines permissions for kubelets. If the role is modified, it can prevent kubelets from working.
```
In otherwords system: is an indication you *can't* delete something.  Otherwise, go for it.

## RBAC: what does it need to do?

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  creationTimestamp: 2019-02-24T03:32:11Z
  name: kubeadm:kubelet-bootstrap
  resourceVersion: "207"
  selfLink: /apis/rbac.authorization.k8s.io/v1/clusterrolebindings/kubeadm%3Akubelet-bootstrap
  uid: c586609e-37e4-11e9-b5ac-0800272e45d7
roleRef:
- apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:node-bootstrapper
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:bootstrappers:kubeadm:default-node-token
```

## I see something fishy, whats it trying to do?

Auditing...

API Server
```
--audit-policy-file=/log/audit.json
--audit-log-format=json
```

### Example: Minikube

```
{"kind":"Event","apiVersion":"audit.k8s.io/v1","level":"Metadata","auditID":"c6ddfa3d-953d-4aea-b36d-dd7f628bca8b","stage":"ResponseComplete","requestURI":"/api/v1/nodes/minikube/status?timeout=10s","verb":"patch","user":{"username":"system:node:minikube","groups":["system:nodes","system:authenticated"]},"sourceIPs":["127.0.0.1"],"userAgent":"kubelet/v1.13.3 (linux/amd64) kubernetes/721bfa7","objectRef":{"resource":"nodes","name":"minikube","apiVersion":"v1","subresource":"status"},"responseStatus":{"metadata":{},"code":200},"requestReceivedTimestamp":"2019-03-07T22:24:35.711681Z","stageTimestamp":"2019-03-07T22:24:35.722715Z","annotations":{"authorization.k8s.io/decision":"allow","authorization.k8s.io/reason":""}}
```

Note: https://github.com/kubernetes/kubernetes/pull/71230 at some point, creation of audit policies can be done and maintained in the cluster
as a standard API object. 

## Building a webhook to model API usage in your cluster

Golang example:

```
mkdir /root/bin ; curl -sL -o ~/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x ~/bin/gimme
~/bin/gimme 1.11 # Shows the env vars 
eval "$(GIMME_GO_VERSION=1.11 gimme)" # sets them as default, now go is ready
export GOPATH=/root/go/
```

Example log dataset for a long running cluster for auditing events:

https://gist.githubusercontent.com/jayunit100/fdcd8b5edb3f6e38191da9f435ec9d09/raw/08a3f8951fa82b9d4253f6211f83b72427b9e1a3/

## What to look for ? 




# Modelling API server vulnerabilities w/ kube-hunter
```
[centos@ip-10-0-3-30 ~]$ sudo docker run -it --rm --network host aquasec/kube-hunter
```

Scans for ports that have vulnerable information on them, for example, the kubelet r/o endpoint:
```
| Kubelet API (readonly):
|   type: open service
|   service: Kubelet API (readonly)
|_  host: 10.0.3.30:10255
```
Which has data about all your workloads !

```
[centos@ip-10-0-3-30 ~]$ curl 10.20.5.1:10255/pods
{"kind":"PodList","apiVersion":"v1","metadata":{},"items":[{"metadata":{"name":"metrics-server-v0.2.1-675ccb567f-x7fbt","generateName":"metrics-server-v0.2.1-675ccb567f-","namespace":"kube-system","selfLink":"/api/v1/namespaces/k
ube-system/pods/metrics-server-v0.2.1-675ccb567f-x7fbt","uid":"d4d08bb2-3675-11e9-a097-0add60e45816","resourceVersion":"537","creationTimestamp":"2019-02-22T07:45:32Z","labels":{"k8s-app":"metrics-server","pod-template-hash":"231
7761239","version":"v0.2.1"},"annotations":{"kubernetes.io/config.seen":"2019-02-22T16:41:17.609318905Z","kubernetes.io/config.source":"api","scheduler.alpha.kubernetes.io/critical-pod":""},"ownerReferences":[{"apiVersion":"exten
s
```


