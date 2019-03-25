echo "What does a real K8s patch look like ? "

cat << EOF
        cp /usr/bin/docker-runc ./backup/docker-runc
        sudo mv runc-v17.03.2-amd64 /usr/bin/docker-runc
EOF

echo	"Patching docker looks somthing like this*"

cat << EOF
        pdsh -u centos -w hosts patch.sh
        cp runc-v17.06.2-amd64 /usr/bin/docker-runc
EOF


cat << EOF

Easier fix, assuming you turned off hostnetworks, enabled adm ctl
```
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: non-root
  spec:
    privileged: false
      allowPrivilegeEscalation: false
        runAsUser:
	    rule: 'MustRunAsNonRoot'
```
EOF
