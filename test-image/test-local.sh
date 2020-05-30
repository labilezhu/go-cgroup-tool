mkdir /home/labile/mycore

./test-image/build-image.sh

docker run \
-v /home/labile/mycore:/mycore \
--privileged --ulimit core=-1 \
--memory=400m --memory-swap=400m --memory-swappiness=0 --rm \
go-cgroup-too-test:latest