# Threat modelling in kubernetes clusters ++

## Reach out + thanks to the ASF

-  https://kubernetes.io/community/
-  platform9.com / @platform9sys 
-  jay@apache.org / @jayunit100 

---

# Two qoutes

- "HTTPS is the universal firewall bypass protocol" - John Morello (CTO @ Twistlock)
- "Security at scale has to be transparent to applications" - Samrat Ray (PM @ Google)

https://jayunit100asf.platform9.horse/clarity/index.html#/infrastructureK8s#clusters

- I'm @jayunit100, and I work on K8s @platform9sys
- By the way we can manage this entire problem for you :)

---

# *TL;DR* ++

- https://kubernetes.io/docs/tasks/administer-cluster/securing-a-cluster/
- https://kubernetes.io/docs/reference/access-authn-authz/rbac/
- https://kubernetes.io/blog/2018/07/18/11-ways-not-to-get-hacked/

---

*What does a real K8s patch look like? Patching runc...* ++

- https://jayunit100asf.platform9.horse/clarity/index.html#/infrastructureK8s#clusters

`cat patch-cve.sh` bait-&-switch

- *https://github.com/rancher/runc-cve/releases*
- --enable-admission-plugins=PodSecurityPolicy
- Easier fix: `privileged: false,allowPrivilegeEscalation: false, runAsUser:rule: 'MustRunAsNonRoot'`

--- 

Modelling Threats != blocking ports + reading CVEs all day + boring rules ++

- in K8s num of IPTables rules ~ O(services)
- NetworkPolicy API (cheaper) vs.  multiple clusters (easier).
- Vanity SSL vs. Certificate chains, rotation, expiration
- App CVE false positives, boooring? `exec` sometimes IS required (i.e. apache https!)
- mismatched vulnerabilities

- minikube ssh `iptables -L | grep default` 56
- `sudo iptables --table nat --list  | wc -l` 148 -> 158 

---
*First do no harm: Change the Question* ++

*'Sometimes the best kind of engineering involves changing the question" - Joe Beda (Heptio)*

- *Pay someone*: https://platform9.com/blog/the-seamless-upgrade-for-kubernetes-first-major-security-hole-cve-2018-1002105/
- Level-up:  *serverless+ingress* level as a way to unify your threat model.

ISTIO ISTIO ISTIO
## Automate tooling around citadel/*zero-trust* including your developers

---

First do no harm: Make sure you didnt break engineers and developers when you fixed security ++

- cat e2e-logs.txt
- https://github.com/cncf/k8s-conformance/blob/master/v1.9/platform9/e2e.log
- https://github.com/cncf/k8s-conformance/blob/master/v1.9/platform9/junit_01.xml

- Run kube-hunter or aqua or opssight or twistlock or whatever 
- Lock something down
- Run the conformance tests `https://scanner.heptio.com/`
- Run a kubetest variant, even easier... 
There are ~ 130 Conformance tests / 30 minutes to run.

--- 

*First do no harm: How to run the kubetests from source to target specific functionality* ++

- cat demo-e2e.sh 

Stuff you might have broke? The stuff in your threat model.
Whats in your threat model to start with?  The stuff you think is important.

APIServer,ETCD,Kublet,Secrets,Network DS,Volumes, EmptyDir,VolPerms,...

--- 
App security: ++

- Twistlock, Blackduck, Claire
- Watch out for Container exec + consider scratch
- API calls that make apps do weird stuff
*Application security*
 
(((More on this later)))

3 paradigms to explore:

- CoreOS Claire: Quickly extract an apses profile by layers.
- Blackduck / OpsSight : Scan everything deeply.
- Twistlock : Anomolous behavior, listening on sockets, weird files.

---
*Cluster Threat Model:  APIServer: ClusterRoleBindings* ++

`cat rbac-example.yml`

- Roles are good, ClusterRoles are OK, ClusterRoleBindings are dangerous.

- Allow you to do things in ANY NAMESPACE  Look at them, maybe delete them...

.. Compare ...

- kubectl get role --all-namespaces
- kubectl get clusterrole --all-namespaces
- kubectl get clusterrolebindings

--- 

Operators: ++

- Intelligently modify CRD permissions using new aggreagators (rather
then editing user roles)

Example of how to declaratively cascade operator RBAC rules

- labels:
  - rbac.authorization.k8s.io/aggregate-to-admin: "true"

Thanks @liggit

--- 

ClusterThreatModel: RBAC: what does it need to do? ++

`cat rbac-example.yml`

- Read config maps, God privileges, rbac-to-allow (selinux ~ audit-to-allow)
- helm doesnt need god-privileges PSP example

```yaml
	apiVersion: rbac.authorization.k8s.io/v1
	kind: ClusterRoleBinding
	metadata:
	  name: tiller
	roleRef:
	  apiGroup: rbac.authorization.k8s.io
	  kind: ClusterRole
	  name: cluster-admin
	subjects:
	  - kind: ServiceAccount
	    name: tiller
	    namespace: kube-system
```

---
 
*ClusterThreatModel:  finding RBAC anomolies* ++

`cat  auditlogs-record.yml`

- Auditing API !!!

API Server

- --audit-policy-file=/log/audit.json
- --audit-log-format=json

Also, look at audit-to-rbac by jason ligget !

---

*ClusterThreatModel: Modelling RBAC vulnerabilities w/ kube-hunter* ++

- pen-testing
- `cat kube-hunter.sh`

Scans for ports that have vulnerable information on them, for example, the kubelet r/o endpoint...

---

*Example vulnerability* ++

- Kubelet API (readonly), type: open service, service: Kubelet API (readonly),  host: 10.0.3.30:10255

*Force the Kubelet to require HTTPS auth from metrics server.*
But you just broke metrics, nodemetadata, logging... 

```
unable to fully collect metrics: 
unable to fully scrape metrics from source 
kubelet_summary:ip-10-0-1-185.us-west-2.compute.internal: 
unable to fetch metrics from Kubelet ip-10-0-1-185.us-west-2.compute.internal 
```

Solution:

- update the SANs 

---

Kubernetes threat model assets ==  ++
 
- RBAC
- Volumes
- ETCD (everything)
- API Server (get metadata about whats running, logs, exec)
- Services and NodePorts
- Kubelet (CAs, Node privileges)
- Individual containers (app level vulns triggered over endpoints)

--- 

AppSec Threat Modelling ++

Lets scan an app (Airflow) ! 
`cat claire.sh`, `cat vulns.json`

- Most likely, this is where you'll have the most churn of "vulnerabilities".
- Once you expose an endpoint: kubernetes can't necessarily help you very much.
- False positives everywhere ! i.e. cgi scripts and *apache httpd*

---

Rather then micromanaging policies and making people hate you: ++

- Service Mesh
- Go serverless
- API gateways 

All of these do one thing:

DECOUPLE YOUR DEVELOPER WORKFLOW FROM YOUR SECURITY MODEL

---

Blackduck Perceptor: Building an infinite scan queue

`cat build-scan-q.sh`

--- 

Build your own container scanning platform: Perceptor.

- pod-perceivers (k8s)
- image-perceivers (openshift)
- pod scraper perceiver (prototype: openshift, unprivileged)
- skopeo based pod perceiver (opsight 3x, openshift) 
- Example of how threat detection with fine granularity becomes very expensive...
https://github.com/blackducksoftware/opssight-connector/wiki/Image-Starvation
- If security scanning needs to happen at physical locations, scheduling primitives
or node labels:

--- 

Example of how to localize scan locality

```yaml
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

---
