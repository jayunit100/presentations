# Threat modelling in kubernetes clusters

https://jayunit100asf.platform9.horse/clarity/index.html#/infrastructureK8s#clusters

*- I'm @jayunit100, and I work on K8s @platform9sys*
   - By the way we can manage this entire problem for you :)
- https://kubernetes.io/docs/tasks/administer-cluster/securing-a-cluster/
- https://kubernetes.io/docs/reference/access-authn-authz/node/
- https://kubernetes.io/docs/reference/access-authn-authz/rbac/

*- I will not attempt to tell you stuff you already know.*

---

https://jayunit100asf.platform9.horse/clarity/index.html#/infrastructureK8s#clusters

# Example: What does a real K8s patch look like? Patching runc

`cat patch-cve.sh`

- *https://github.com/rancher/runc-cve/releases*
- `--enable-admission-plugins=PodSecurityPolicy,NodeRestriction`
- Easier fix: `privileged: false,allowPrivilegeEscalation: false, runAsUser:rule: 'MustRunAsNonRoot'`

--- 
# First do no harm: Build a threat model

- Apps (security assets  and vulnerabilities specific to your code)

"HTTPS is the universal firewall bypass protocol" - John Morello (CTO @ Twistlock)

- in K8s num of IPTables rules ~ O(services)
- SSL isn't *always* required (sidecars, prototypes)
- CVE what? `exec` sometimes IS required (i.e. apache https!)
- Distributing certs w/o rotating is dumb.
- The only way to do security is to work w/ developers not against them
- Would you have thought about *runc* as a potential threat vector?

# First do no harm: Change the Question

'Sometimes the best kind of engineering involves changing the question" - Joe Beda (Heptio)

- Pay someone: https://platform9.com/blog/the-seamless-upgrade-for-kubernetes-first-major-security-hole-cve-2018-1002105/
- Go serverless+ingress level as a way to unify your threat model.

'Security at scale has to be transparent to applications' - Samrat Ray (Google)

- Automate tooling around citadel/zero-trust models

`cat e2e-logs.txt`

# First do no harm: Make sure you didnt break engineers and developers when you fixed security

https://github.com/cncf/k8s-conformance/blob/master/v1.9/platform9/e2e.log
https://github.com/cncf/k8s-conformance/blob/master/v1.9/platform9/junit_01.xml

KUBERNETES HAS END TO END TESTS TO VALIDATE ALL OF ITS FUNCTIONALITY
THAT CAN BE RUN BY ANYONE FROM ANYWHERE WITH HUMAN READABLE RESULTS

- Run kube-hunter or aqua or opssight or twistlock or whatever 
- Lock something down
- Run the conformance tests `https://scanner.heptio.com/`
- Run a kubetest variant, even easier, like this

There are ~ 130 Conformance tests / 30 minutes to run.

--- 

# First do no harm: How to run the kubetests from source to target specific functionality

`cat demo-e2e.sh`

Stuff you might have broke? The stuff in your threat model.


Master:
- ETCD ~ it has your entire cluster contents.
- Kubelets ~ they run as root, and they generally have a docker daemon on them.
- API Server ~ bringing an API server down could DoS any 'cloud native' app in a cluster.
Apps:
- Container exec
- API calls that make apps do weird stuff

---

# Cluster Threat Model:  APIServer: ClusterRoleBindings

Allow you to do anything.  Look at them, DELETE THEM !

`kubectl get clusterrolebindings`

```
+------------------------------------+
|              Minikube              |
+------------------------------------+
| cluster-admin                      |
| kubeadm:kubelet-bootstrap          |
| kubeadm:node-autoapprove-bootstrap |
| minikube-rbac                      |
+------------------------------------+
```

--- 

# ClusterThreatModel: RBAC: what does it need to do?

`cat rbac-example.yml`

- Read config maps
- God privileges
- rbac-to-allow (selinux ~ audit-to-allow)
- helm PSP example

---
 
# ClusterThreatModel:  finding RBAC anomolies

Auditing...

API Server
```
--audit-policy-file=/log/audit.json
--audit-log-format=json
```

# Example: Minikube ~ audit logs

- Breifly: How minikube actually works, static pods

`cat auditlogs-record.yml`

- Note: https://github.com/kubernetes/kubernetes/pull/71230 at some point, creation of audit policies can be done and maintained in the cluster
as a standard API object.

- Example log dataset for a long running cluster for auditing events:

- https://gist.githubusercontent.com/jayunit100/fdcd8b5edb3f6e38191da9f435ec9d09/raw/08a3f8951fa82b9d4253f6211f83b72427b9e1a3/

---

# ClusterThreatModel: Modelling RBAC vulnerabilities w/ kube-hunter ~ pen-testing
`cat kube-hunter.sh`

Scans for ports that have vulnerable information on them, for example, the kubelet r/o endpoint:

```
	Kubelet API (readonly), type: open service, service: Kubelet API (readonly),  host: 10.0.3.30:10255
```

Force the Kubelet to require HTTPS auth from metrics server.
But you just broke metrics, nodemetadata, logging... 

```
 unable to fully collect metrics: 
 unable to fully scrape metrics from source 
   kubelet_summary:ip-10-0-1-185.us-west-2.compute.internal: 
unable to fetch metrics from Kubelet ip-10-0-1-185.us-west-2.compute.internal 
```

- turn off insecure-tls
- update the SANs 
---

# So thats the kubernetes threat model

- ETCD (everything)
- API Server (get metadata about whats running, logs, exec)
- Kubelet (CAs, Node privileges)
- Individual containers (app level vulns triggered over endpoints)

--- 

# Application Level Threat Model

Lets scan an app (Airflow) ! 
`cat claire.sh`, `cat vulns.json`

- Most likely, this is where you'll have the most churn of "vulnerabilities".
- Once you expose an endpoint: kubernetes can't necessarily help you very much.
- False positives everywhere ! i.e. cgi scripts and *apache httpd*

3 rules of crafting your threat model:

- Know what your application does and why, before locking things down.
- Understand where your important assets are, worry about those first.
- Decouple your developer workflow from your security model,
so they can evolve orthogonally... unless your goal is to kill feature
velocity.

3 paradigms to explore::

- CoreOS Claire: Quickly extract an apses profile by layers.
- Blackduck / OpsSight : Scan everything deeply.
- Twistlock : Anomolous behavior, listening on sockets, weird files.

--- 

# How to continously build a scan queue that monitors all applications in your k8s cluster.

`cat build-scan-q.sh`

You could totally automate this into a workqueue: See Blackduck OpsSight product for details!

- perceptor: infinite scan queue
- pod-perceivers (k8s)
- image-perceivers (openshift)
- pod scraper perceiver (prototype: openshift, unprivileged)
- skopeo based pod perceiver (opsight 3x, openshift) 

- Example of how threat detection with fine granularity becomes very expensive...
https://github.com/blackducksoftware/opssight-connector/wiki/Image-Starvation

- If security scanning needs to happen at physical locations, scheduling primitives
or node labels:
```
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: "my.pds.opssightnode"
            operator: In
            values: ["true"]
```
