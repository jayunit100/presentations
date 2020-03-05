export KUBECONFIG=`kind get kubeconfig-path --name=calico-test` 

######################## 

kubectl create clusterrolebinding netpol --clusterrole=cluster-admin --serviceaccount=kube-system:netpol
kubectl create sa netpol -n kube-system


#######################

kubectl delete job netpol -n kube-system
kubectl delete ns x
kubectl delete ns y
kubectl delete ns z

sleep 2

#######################

kubectl create -f https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/netpol/install-latest.yml -n kube-system 

until kubectl logs -f `kubectl get pods -n kube-system | grep net | cut -d' ' -f 1` -n kube-system ; do 
	echo "trying again in 3 seconds..."
	sleep 3
done


