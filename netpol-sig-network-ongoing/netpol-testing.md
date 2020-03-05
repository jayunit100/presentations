# NetworkPolicy testing framework with truth tables

- Incompleteness: The tests API should be written to support increasing test coverage over time
- Understandability: The tests should be easy to read and help interpret policy semantics
- Dynamic scale: The test API be declarative enough to support larger scale tests for enterprises/cloud providers
- Documentation and Community: We should have a broader documentation offering around network policies
- Continous Integration: The network policy tests should run in automation either in K8s CI or in builders curated by the community to catch bugs.  Perf tests for CNIs might be nice to (since some issues might manifest in the Kubelet)

## Motivation and Goals

Set a standard for continously running, easily adaptable network policy semantics and API validation that is readily adoptable by cloud providers providing confidence to users that all networkpolicy semantics work well at any scale.

## Design

Excerpt from our existing KEP, living https://github.com/kubernetes/enhancements/pull/1568.
These tests work by probing all 81 possible connections between containers, ( 3 namespaces, 3 containers
identical in each namespace).

```
+-------------------------------------------------------------------------+
|  +------+              +------+                                         |
|  |      |              |      |                                         |
|  |   cA |              |  cB  |     Figure 1b: The above test           |
|  +--+---+              +----X-+     is only complete if a permutation   |
|     |   +---------------+   X       of other test scenarios which       |
|     |   |    server     |   X       guarantee that (1) There is no      |
|     +--->    80,81      XXXXX       namespace that whitelists traffic   |
|         |               |           and that (2) there is no "container"| TODO: test "default" namespace
|         +----X--X-------+           which whitelists traffic.           |       check for dropped namespaces
| +------------X--X---------------+                                       |       make test instances bidirectional
| |            X  X               |   We limit the amount of namespaces   |          (client/servers)
| |   +------XXX  XXX-------+  nsB|   to test to 3 because 3 is the union |
| |   |      | X  X |       |     |   of all namespaces.                  |
| |   |  cA  | X  X |   cB  |     |                                       |
| |   |      | X  X |       |     |   By leveraging the union of all      |
| |   +------+ X  X +-------+     |   namespaces we make *all* network    |
| |            X  X               |   policy tests comparable,            |
| +-------------------------------+   to one another via a simple         |
|  +-----------X--X---------------+   truth table.                        |
|  |           X  X               |                                       |
|  |  +------XXX  XXX-------+  nsC|   This fulfills one of the core       |
|  |  |      |      |       |     |   requirements of this proposal:      |
|  |  |  cA  |      |   cB  |     |   comparing and reasoning about       |
|  |  |      |      |       |     |   network policy test completeness    |
|  |  +------+      +-------+     |   in a deterministic manner which     |
|  |                              |   doesn't require reading the code.   |
|  +------------------------------+                                       |
|                                      Note that the tests above are all  |
|                                      done in the "framework" namespace, |
|                                                  similar to Figure 1.   |
+-------------------------------------------------------------------------+
```

kubectl exec -t -i zb-659ddf6cd9-fdpqs -c c80 -n z -- wget --spider --tries 4 --timeout 0.5 --waitretry 0 http://192.168.242.197:80
## Results

- Coverage for all corner cases (81 connections) is implemented and working, covering the same semantics as upstream (roughly).
- The `Probe` functionality results in very rapid discerning of connectivity.  To run 30 of our tests takes about 20 minutes or less, pass or fail.
- Running the existing policy stack takes ~ one hour.  If ETCD is slow or pod startup times are slow, this can be much longer, and in broken clusters NetworkPolicy E2Es can take ~ 2 hours.
- Tests are now hackable - you can leave the pods up and reproduce failures from the outputted kubectl commands:
`kubectl exec -t -i zb-659ddf6cd9-fdpqs -c c80 -n z -- wget --spider --tries 4 --timeout 0.5 --waitretry 0 http://192.168.242.197:80`.

## Current state

- An implementation can be found [here](https://github.com/vmware-tanzu/antrea/tree/master/hack/netpol).
- Runs about 14 test cases in about 10mins.
- In comparision, e2e tests focused on NetworkPolicy takes close to an hour.
- Recently integrated with Antrea CI.
- Run it:

```
kubectl create clusterrolebinding netpol --clusterrole=cluster-admin --serviceaccount=kube-system:netpol
kubectl create sa netpol -n kube-system
kubectl create -f https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/netpol/install-latest.yml -n kube-system 
kubectl get pods -n kube-system # <- results will be in the netpol pod
```

## Next steps

- Goal is to consolidate more requirements around our KEP  :https://github.com/kubernetes/enhancements/pull/1568/files
- More test cases for a full coverage of NetworkPolicy spec, especially scale.
- Cleanup flag to determine if resources created by NetPol must be deleted.
- Node specific test cases.
- Extend framework to run scale tests for NetworkPolicy.
- Extend framework to test other K8s resources with truth tables.
- Porting to Ginkgo (oh no)!
- Working w/ Sig-testing on automation requisites and CI parts
