# Threat modelling in kubernetes cluster

## Whats a threat model?

According to wikipedia...
```
Threat modeling is a process by which potential threats, 
such as structural vulnerabilities can be 
identified, enumerated, and prioritized – 
all from a hypothetical attacker’s point of view. 
```
More definition bikeshedding...
```
The purpose of threat modeling is to provide defenders 
with a systematic analysis of 
the probable attacker’s profile, the most likely attack vectors, and 
the assets most desired by an attacker. 
```

Specifically:

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

# CLUSTER MODEL 

## APIServer: ClusterRoleBindings

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

EASIEST THING YOU CAN DO TO AVOID 
ROOT ACLs TO THE API SERVER IS PAY
ATTENTION !

Namespace your CRBs, so you know who made them and why ! 

See https://github.com/kubernetes/minikube/issues/3825 

Example: In a managed kubernetes distribution

```
admin-and-blahblah-access                              13d
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

- Roles ~ FreeIPA: Permissions
```
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
kubeadm:bootstrap-signer-clusterinfo
rules:
- apiGroups:
  - ""
  resourceNames:
  - cluster-info
  resources:
  - configmaps
  verbs:
  - get
```

God privileges look like this: 
```
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
  - '*'
```

- ClusterRoleBindings ~ FreeIPA: Privileges
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
roleRef:
- apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admino
subjects:

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
{  
   "kind":"Event",
   "apiVersion":"audit.k8s.io/v1",
   "level":"Metadata",
   "auditID":"c6ddfa3d-953d-4aea-b36d-dd7f628bca8b",
   "stage":"ResponseComplete",
   "requestURI":"/api/v1/nodes/minikube/status?timeout=10s",
   "verb":"patch",
   "user":{  
      "username":"system:node:minikube",
      "groups":[  
         "system:nodes",
         "system:authenticated"
      ]
   },
   "sourceIPs":[  
      "127.0.0.1"
   ],
   "userAgent":"kubelet/v1.13.3 (linux/amd64) kubernetes/721bfa7",
   "objectRef":{  
      "resource":"nodes",
      "name":"minikube",
      "apiVersion":"v1",
      "subresource":"status"
   },
   "responseStatus":{  
      "metadata":{  

      },
      "code":200
   },
   "requestReceivedTimestamp":"2019-03-07T22:24:35.711681Z",
   "stageTimestamp":"2019-03-07T22:24:35.722715Z",
   "annotations":{  
      "authorization.k8s.io/decision":"allow",
      "authorization.k8s.io/reason":""
   }
}
```

Note: https://github.com/kubernetes/kubernetes/pull/71230 at some point, creation of audit policies can be done and maintained in the cluster
as a standard API object. 

## Building a webhook to model API usage in your cluster

Example log dataset for a long running cluster for auditing events:

https://gist.githubusercontent.com/jayunit100/fdcd8b5edb3f6e38191da9f435ec9d09/raw/08a3f8951fa82b9d4253f6211f83b72427b9e1a3/

## Modelling API server vulnerabilities w/ kube-hunter
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

## Runtime security at the App level

```
docker run -d --name db arminc/clair-db:2018-04-01
docker run -p 6060:6060 --link db:postgres -d --name clair arminc/clair-local-scan:v2.0.1
```

## So thats the kubernetes threat model

From an attackers perspective: The API Server, and the Kubelet, represent
to chatty components which can be influenced via API access over well known
ports.  Overall solution = SSL, enterprise support and patching.

Now, lets look at the other, continuously shifting model: Your apps.




# The Authentication / Authorization Model



# Application Level Threat Model

These are more likely to be exploited - especially if you're using a kubernetes distribution
from an enterprise grade company. 

## What are YOU doing wrong in your apps

Most likely, this is where you'll have the most churn of vulnerabilities.
Once you expose an endpoint: kubernetes can't necessarily help you very much.

### Fun part: Pick an image !

Container scanning

- CoreOS Claire
- OpsSight Connector

### Example + Demo

Find all your images:

```
	kubectl get pods --all-namespaces -o jsonpath="{.items[*].spec.containers[*].image} \n" 
```

Run the clair-scanner against them:
```
	docker pull apache/airflow:latest
	./clair-scanner_darwin_amd64 --report=vulns.json  --threshold="Critical" --ip=192.168.20.194 apache/airflow:latest	
```

#### Apache Airflow: Sample output

Apache Airflow:

```
	"description": "An issue was discovered in shadow 4.5. newgidmap (in shadow-utils) is setuid and allows an unprivileged user to be placed in a user namespace where setgroups(2) is permitted. This allows an attacker to remove themselves from a supplementary group, which may allow access to certain filesystem paths if the administrator has used \"group blacklisting\" (e.g., chmod g-rwx) to restrict access to paths. This flaw effectively reverts a security feature in the kernel (in particular, the /proc/self/setgroups knob) to prevent this sort of privilege escalation."
```


# Thanks !
