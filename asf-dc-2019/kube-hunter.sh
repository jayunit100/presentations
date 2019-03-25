echo "enter 1, then 127.0.0.1"

sudo docker run -it --rm --network host aquasec/kube-hunter

echo "For example..."
cat << EOF

```
| Kubelet API (readonly):
|   type: open service
|   service: Kubelet API (readonly)
|_  host: 10.0.3.30:10255
```
And then you'll see something like this:
```
| Kubelet API          | 127.0.0.1:10255 | The read-only port   |
| (readonly)           |                 | on the kubelet       |
|                      |                 | serves health        |
|                      |                 | probing endpoints,   |
|                      |                 | and is relied upon   |
|                      |                 | by many kubernetes   |
|                      |                 | componenets          |

| 127.0.0.1:10250 | Remote Code          | Anonymous            | The kubelet is
|                 | Execution            | Authentication       | misconfigured,
|                 |                      |                      | potentially allowin
|                 |                      |                      | secure access to all
|                 |                      |                      | requests on the
|                 |                      |                      | kubelet, without the
|                 |                      |                      | need to authenticate
```


[centos@ip-10-0-3-30 ~]$ curl my-kubelet-ip:10255/pods
{"kind":"PodList","apiVersion":"v1","metadata":{},"items":[{"metadata":{"name":"metrics-server-v0.2.1-675ccb567f-x7fbt","generateName":"metrics-server-v0.2.1-675ccb567f-","namespace":"kube-system","selfLink":"/api/v1/namespaces/k
ube-system/pods/metrics-server-v0.2.1-675ccb567f-x7fbt","uid":"d4d08bb2-3675-11e9-a097-0add60e45816","resourceVersion":"537","creationTimestamp":"2019-02-22T07:45:32Z","labels":{"k8s-app":"metrics-server","pod-template-hash":"231
7761239","version":"v0.2.1"},"annotations":{"kubernetes.io/config.seen":"2019-02-22T16:41:17.609318905Z","kubernetes.io/config.source":"api","scheduler.alpha.kubernetes.io/critical-pod":""},"ownerReferences":[{"apiVersion":"exten
s

EOF


