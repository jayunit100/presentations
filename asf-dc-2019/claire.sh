echo "set your IP to 10.0.0.28, or whatever the en0 IP is!"

docker run -d --name db arminc/clair-db:2018-04-01
docker run -p 6060:6060 --link db:postgres -d --name clair arminc/clair-local-scan:v2.0.1
docker pull apache/airflow:latest
./clair-scanner_darwin_386 --report=vulns.json  --threshold="Critical" --ip=10.0.0.28 apache/airflow:latest
