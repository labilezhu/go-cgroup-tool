echo '/mycore/core.%h.%e.%t' > /proc/sys/kernel/core_pattern
ulimit -c unlimited

cat /proc/sys/kernel/core_pattern

/opt/java/openjdk/bin/java \
-XX:MaxDirectMemorySize=536870912 \
-XX:NativeMemoryTracking=summary \
-jar /stress-java/stress-java-0.0.1-SNAPSHOT.jar &

child=$!
echo "child=$child"

/stress-java/cgrouptool \
-level 200:"jcmd $child VM.native_memory summary" \
-level 300:"kill -3 $child" &

cgrouptool_pid=$!

wait "$child"
kill -9 $cgrouptool_pid

