FROM adoptopenjdk/openjdk8:x86_64-debian-jdk8u252-b09-slim

ADD stress-java/target/stress-java-0.0.1-SNAPSHOT.jar /stress-java/stress-java-0.0.1-SNAPSHOT.jar

ADD src/cgrouptool/cgrouptool /stress-java/cgrouptool

ADD test-image/docker-entrypoint.sh /stress-java/docker-entrypoint.sh

CMD ["/bin/bash", "/stress-java/docker-entrypoint.sh"]