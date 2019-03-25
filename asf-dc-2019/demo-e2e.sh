mkdir ~/bin
curl -sL -o ~/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x ~/bin/gimme
export PATH=$PATH:~/bin
export GOPATH=~/go
eval "$(gimme stable)"
KUBERNETES_DIR=kubernetes
KUBERNETES_VERSION=v1.12.5
git clone --branch ${KUBERNETES_VERSION} https://github.com/kubernetes/kubernetes.git
pushd ${KUBERNETES_DIR} 
	make WHAT=test/e2e/e2e.test && ./build/run.sh make cross && make kubectl && make ginkgo
popd

 go get -u k8s.io/test-infra/kubetest
 export KUBE_MASTER_IP=ojas-deprecate-heapster-c8fb92c1-api.ojas-test.platform9.horse
 export KUBE_MASTER=ip-10-0-1-125.us-west-2.compute.internal
 export KUBECONFIG="~/kube/config.yml"
pushd kubernetes/
      $GOPATH/bin/kubetest --test --test_args="--ginkgo.focus=Secrets" --provider=skeleton
popd

