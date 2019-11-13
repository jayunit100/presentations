#  How to build and release Calico

## Thanks!

-  jvyas@vmware.com / @jayunit100

## Calico Build infra

Upstream images are available.  You don't have to do this !

```
docker pull calico/node:release-v3.10
docker pull calico/kube-controllers:release-v3.10
```
You can see the latest image builds on:

```
https://semaphoreci.com/calico/
```

## Why build it then? 

- Inject metrics or logging into different parts (i.e. felix, etc) to see how individual components work.
- You may want to patch it yourself or you may have organization requirements for image base layers etc
- You want to change the behaviour of a particular use case
- You want to contribute to calico and need a way to test your patches.
- You want to reproduce or report bugs.
- You want to backport a CVE fix etc to an older branch which might not be released by upstream.

## What images run in a K8s cluster:

In a running k8s cluster, you'll see the following images running for calico.

```
    node1: NAMESPACE     NAME                                       READY   STATUS    RESTARTS   AGE
    node1: kube-system   calico-kube-controllers-6bdf954cd6-pvnx9   0/1     Running   0          13s
    node1: kube-system   calico-node-9p9vh                          1/1     Running   0          13s
    node1: kube-system   coredns-5644d7b6d9-h8l9d                   0/1     Running   0          13s
    node1: kube-system   coredns-5644d7b6d9-pg227                   0/1     Running   0          13s
    node1: kube-system   etcd-node1                                 1/1     Running   0          27s
    node1: kube-system   kube-apiserver-node1                       1/1     Running   0          27s
    node1: kube-system   kube-controller-manager-node1              1/1     Running   0          27s
    node1: kube-system   kube-proxy-xkhrk                           1/1     Running   0          13s
    node1: kube-system   kube-scheduler-node1                       1/1     Running   0          27s
    node1: Your dev VM is up, vagrant ssh to access it.   TEST RESULT: passed.
```

## What components does calico build ? 

Roughly, don't qoute me on this... 

- `RELEASE_REPOS=felix typha kube-controllers calicoctl cni-plugin app-policy pod2daemon node` (from the Calico Makefile).
- cni-plugin -> libcalico
- node -> felix, libcalico, confd (calico specific fork )
- felix -> pod2daemon, libcalico, confd (calico specific fork) 
- typha -> libcalico (optional)
- kubecontrollers -> libcalico

## Structure of the dev environment:

- hack/development-environment: Vagrantfile w/ build automation
- not specific to vagrant: install.sh can run on any centos box (easily changeable to deb if we want to)
- submodules: create a new directory, modify the vagrantfile to mount your submodules:
```
➜  hack git:(remotes/origin/release-v3.10-tanzu-master) ls -altrh vmware
total 16
drwxr-xr-x   6 jayunit100  staff   192B Nov  6 10:55 ..
-rw-r--r--   1 jayunit100  staff   511B Nov  6 10:55 README.md
-rwxr-xr-x   1 jayunit100  staff   477B Nov  6 10:55 setup.sh
drwxr-xr-x  14 jayunit100  staff   448B Nov  6 10:55 .
drwxr-xr-x  27 jayunit100  staff   864B Nov  6 10:55 app-policy
drwxr-xr-x  26 jayunit100  staff   832B Nov  6 10:55 calicoctl
drwxr-xr-x  31 jayunit100  staff   992B Nov  6 10:55 cni-plugin
drwxr-xr-x  21 jayunit100  staff   672B Nov  6 10:55 confd
drwxr-xr-x  56 jayunit100  staff   1.8K Nov  6 10:55 felix
drwxr-xr-x  26 jayunit100  staff   832B Nov  6 10:55 kube-controllers
drwxr-xr-x  21 jayunit100  staff   672B Nov  6 10:55 libcalico-go
drwxr-xr-x  20 jayunit100  staff   640B Nov  6 10:55 pod2daemon
drwxr-xr-x  23 jayunit100  staff   736B Nov  6 11:38 typha
drwxr-xr-x  27 jayunit100  staff   864B Nov  6 15:17 node
```

## How I use hack/development-environment to get my job done.

- Internal build infra has one submodule, calico, which is a mirror
- calico submodule has 10 submodules inside the hack/vmware directory.
- install.sh modified to mount hack/vmware/* into calico_all.

Example: Top level calico submodule, which recursively has all calico sub repos in its hack/vmware directory:
```
[submodule "calico_all/src"]
        path = calico_all/src
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_calico
        branch = release-v3.10-tanzu-master
        # Remember that calico is always master, even when were building 3.10.
```

Then, we mave submodules inside calico:

```
➜  hack git:(remotes/origin/release-v3.10-tanzu-master) cat ../.gitmodules 
[submodule "hack/vmware/libcalico-go"]
        path = hack/vmware/libcalico-go
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_libcalico-go
        branch = release-v3.10-tanzu
[submodule "hack/vmware/confd"]
        path = hack/vmware/confd
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_confd
        branch = release-v3.10-tanzu
[submodule "hack/vmware/felix"]
        path = hack/vmware/felix
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_felix
        branch = release-v3.10-tanzu
[submodule "hack/vmware/typha"]
        path = hack/vmware/typha
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_typha
        branch = release-v3.10-tanzu
[submodule "hack/vmware/kube-controllers"]
        path = hack/vmware/kube-controllers
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_kube-controllers
        branch = release-v3.10-tanzu
[submodule "hack/vmware/calicoctl"]
        path = hack/vmware/calicoctl
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_calicoctl
        branch = release-v3.10-tanzu
[submodule "hack/vmware/app-policy"]
        path = hack/vmware/app-policy
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_app-policy
        branch = release-v3.10-tanzu
[submodule "hack/vmware/pod2daemon"]
        path = hack/vmware/pod2daemon
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_pod2daemon
        branch = release-v3.10-tanzu
[submodule "hack/vmware/node"]
        path = hack/vmware/node
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_node
        branch = release-v3.10-tanzu
[submodule "hack/vmware/cni-plugin"]
        path = hack/vmware/cni-plugin
        url = git@gitlab.eng.vmware.com:core-build/mirrors_github_projectcalico_cni-plugin
        branch = release-v3.10-tanzu
```

## LOCAL_BUILD

A few things you should know about when building from source.

LOCAL_BUILD will modify your gomodules, so that they use local urls rather then upstream calico github repos pulled from the internet.
Airgapped builds or building from source for development use this.

PS Thanks to Casey and Rafeal for helping me get this logic sorted out. Turned out there were some upstream ordering bugs in the makefiles we had to fix.
```
local_build:
	$(DOCKER_RUN) $(CALICO_BUILD) go mod edit -replace=github.com/projectcalico/libcalico-go=../libcalico-go
	$(DOCKER_RUN) $(CALICO_BUILD) go mod edit -replace=github.com/projectcalico/confd=../confd
	$(DOCKER_RUN) $(CALICO_BUILD) go mod edit -replace=github.com/projectcalico/felix=../felix
```

How you build calico from source - enabling LOCAL_BUILD, so that it doesnt pull stuff from the internet.
## Vagrantfile useage for simple dev scenario

```
git clone https://github.com/projectcalico/libcalico-go.git
git clone https://github.com/projectcalico/confd.git
git clone https://github.com/projectcalico/felix.git
git clone https://github.com/projectcalico/typha.git
git clone https://github.com/projectcalico/kube-controllers.git
git clone https://github.com/projectcalico/calicoctl.git
git clone https://github.com/projectcalico/app-policy.git
git clone https://github.com/projectcalico/pod2daemon.git
git clone https://github.com/projectcalico/node.git
git clone https://github.com/projectcalico/cni-plugin.git
cd calico/hack/development-environment/
vagrant up
```

## Kind alternative, work in progress

```
#!/bin/bash
function build() {
    sudo apt install make
    pushd ~/calico_all/calico/
        sudo make dev-image REGISTRY=cd LOCAL_BUILD=true
        sudo make dev-manifests REGISTRY=cd
    popd
}
function install_k8s() {
    sudo kind create cluster --config calico-conf.yml
    export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
    sudo chmod 755 ~/.kube/kind-config-kind
    until kubectl cluster-info;  do
        echo "`date`waiting for cluster..."
        sleep 2
    done
}    
function install_calico() {
    kubectl get pods
    pushd ~/calico_all/calico/_output/dev-manifests
            kubectl apply -f ./calico.yaml 
            kubectl get pods -n kube-system
    popd
    sleep 5 ; kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true
    sleep 5 ; kubectl -n kube-system get pods | grep calico-node
    echo "will wait for calico to start running now... "
    while true ; do
        kubectl -n kube-system get pods
        sleep 3
    done
}
build
install_k8s
install_calico
```

# TODO

- airgapped builds: can we do them  ?
- submodules / federated repo, should we curate one in upstream ? 
- reduce code duplication across makefiles --- is there a possible way we can have a global makefile that
does the mod edits and other things generically ?
- map out / make explicit the dependencies of libraries etc (i.e. pod2daemon, libcalico-go vs node, kubecontrollers, calicoctl).

