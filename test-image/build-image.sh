set -e

cd /home/labile/go-cgroup-tool/src/cgrouptool
go build .

cd /home/labile/go-cgroup-tool/stress-java
mvn -Dmaven.test.skip=true  clean package

cd /home/labile/go-cgroup-tool
docker build --no-cache  -t go-cgroup-too-test:latest -f ./test-image/Dockerfile .


