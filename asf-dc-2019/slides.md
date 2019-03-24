# Threat modelling in kubernetes clusters

--- 

- This is a 20 minute talk ! No time for bikeshedding.
- @jayunit100 @platform9sys
- https://kubernetes.io/docs/tasks/administer-cluster/securing-a-cluster/
- https://kubernetes.io/docs/reference/access-authn-authz/node/
- https://kubernetes.io/docs/reference/access-authn-authz/rbac/

---

# Some opinions that you might not like

- Model your threats first before making dumb policies.
- SSL isn't *always* required (sidecars, prototypes)
- `exec` sometimes IS required (i.e. apache https!)
- Distributing certs w/o rotating is dumb.
- The only way to do security is to work w/ developers not against them

---

# Some more opinions: Change the Question

- https://platform9.com/blog/the-seamless-upgrade-for-kubernetes-first-major-security-hole-cve-2018-1002105/
- Secure at the serverless+ingress level  as a way to unify your threat model.
- Give developers awesome tooling that abstracts security.

---

### K8s Hardening Workflow that doesn't suck 

- Run kube-hunter or aqua or opssight or twistlock or whatever 
- Lock something down
- Run the conformance tests `https://scanner.heptio.com/`
- Run a kubetest variant, even easier, like this

```
KUBE_MASTER_IP=pmdb-stage-api.snn1.pf9.io 
KUBE_MASTER_NAME=pmdb-stage kubetest 
--test --provider=skeleton --cluster=pmdb-stage --test_args 
--ginkgo.focus="Networking"
```

--- 

To see full test logs, you can view the public archives:

i.e. https://github.com/cncf/k8s-conformance/blob/master/v1.9/platform9/e2e.log

---

## Example conformance results

```
kubernetes git:(3a1c9449a9) 
# should be consumable from pods in volume 
as non-root with defaultMode and fsGroup 
set [Feature:FSGroup]
```

- ETCD ~ it has your entire cluster contents.
- Kubelets ~ they run as root, and they generally have a docker daemon on them.
- API Server ~ bringing an API server down could DoS any 'cloud native' app in a cluster.
- Apps (security assets  and vulnerabilities specific to your code)
- Serices (downtime)

In this talk: Focus on kubelets, API server, ETCD, Apps.

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
--- 
# Which has data about all your workloads ! ANYONE CAN SEE THIS!

```
[centos@ip-10-0-3-30 ~]$ curl 10.20.5.1:10255/pods
{"kind":"PodList","apiVersion":"v1","metadata":{},"items":[{"metadata":{"name":"metrics-server-v0.2.1-675ccb567f-x7fbt","generateName":"metrics-server-v0.2.1-675ccb567f-","namespace":"kube-system","selfLink":"/api/v1/namespaces/k
ube-system/pods/metrics-server-v0.2.1-675ccb567f-x7fbt","uid":"d4d08bb2-3675-11e9-a097-0add60e45816","resourceVersion":"537","creationTimestamp":"2019-02-22T07:45:32Z","labels":{"k8s-app":"metrics-server","pod-template-hash":"231
7761239","version":"v0.2.1"},"annotations":{"kubernetes.io/config.seen":"2019-02-22T16:41:17.609318905Z","kubernetes.io/config.source":"api","scheduler.alpha.kubernetes.io/critical-pod":""},"ownerReferences":[{"apiVersion":"exten
s
```

--- 

The fix?

Force the Kubelet to require HTTPS auth from metrics server.

But you just broke metrics, nodemetadata, logging... 

```
 unable to fully collect metrics: 
[unable to fully scrape metrics from source 
kubelet_summary:ip-10-0-1-185.us-west-2.compute.internal: 
unable to fetch metrics from Kubelet ip-10-0-1
-185.us-west-2.compute.internal 
```

- turn off insecure-tls
- update the SANs 

--- 

# Runtime security at the App level

NOT THAT HARD.. Start a local vulnerability scanner that's FAST (we'll use it in a sec)

```
docker run -d --name db arminc/clair-db:2018-04-01
docker run -p 6060:6060 --link db:postgres -d --name clair arminc/clair-local-scan:v2.0.1
```

--- 
## So thats the kubernetes threat model

- ETCD (everything)
- API Server (get metadata about whats running, logs, exec)
- Kubelet (CAs, Node privileges)
- Individual containers (app level vulns triggered over endpoints)

--- 

# Application Level Threat Model

- Most likely, this is where you'll have the most churn of "vulnerabilities".

- Once you expose an endpoint: kubernetes can't necessarily help you very much.

- False positives everywhere ! Plug ~ Twistlock !

---

### Fun part: Pick an image !

3 paradigms for building a cloud native threat model

- CoreOS Claire: Quickly extract an apses profile by layers.
- Blackduck / OpsSight : Scan everything deeply.
- Twistlock : Anomolous behavior, listening on sockets, weird files.

--- 

#### Example + Demo

Find all your images:

```
	kubectl get pods --all-namespaces -o jsonpath="{.items[*].spec.containers[*].image} \n"
```

You could totally automate this into a workqueue: See Blackduck OpsSight product for details!

--- 


Get the IP of your docker local:

```
Ifconfig | grep -B 4 -A 4 en0
```

Then scan stuff :

```
docker run -d --name db arminc/clair-db:2018-04-01
docker run -p 6060:6060 --link db:postgres -d --name clair arminc/clair-local-scan:v2.0.1
docker pull apache/airflow:latest
./clair-scanner_darwin_amd64 --report=vulns.json  --threshold="Critical" --ip=192.168.20.194 apache/airflow:latest

```


---

## Apache Airflow: Sample output

Apache Airflow:

```
	"description": "An issue was discovered in shadow 4.5."
```

And it goes on... 

```
... newgidmap (in shadow-utils) is setuid and allows an 
unprivileged user to be placed in a user namespace where setgroups(2) 
is permitted. This allows an attacker to remove themselves from a 
supplementary group, which may allow access to certain filesystem paths...
```

# 
