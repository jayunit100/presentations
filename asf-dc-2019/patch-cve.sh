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
