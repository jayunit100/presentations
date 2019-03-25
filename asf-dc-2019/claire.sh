docker run -d --name db arminc/clair-db:2018-04-01
docker run -p 6060:6060 --link db:postgres -d --name clair arminc/clair-local-scan:v2.0.1
docker pull apache/airflow:latest
./clair-scanner_darwin_amd64 --report=vulns.json  --threshold="Critical" --ip=192.168.20.194 apache/airflow:latest
