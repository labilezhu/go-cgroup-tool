package test

import (
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"testing"
)

const limitBytes int64 = 400 * 1024 * 1024

func TestAbs(t *testing.T) {

	var limitBytes int64 = limitBytes
	var swappiness uint64 = 0
	var disableOOMKiller = false
	linuxMemory := specs.LinuxMemory{
		Limit:            &limitBytes,
		Swappiness:       &swappiness,
		DisableOOMKiller: &disableOOMKiller,
	}

	cgroupName := "/test_go_cgroup_tool"

	fmt.Printf("creating cgroup /sys/fs/cgroup%v \n", cgroupName)

	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath(cgroupName), &specs.LinuxResources{
		Memory: &linuxMemory,
	})
	if err != nil {
		t.Errorf("")
		panic("cgroups.New")
	}
	defer control.Delete()

	//control, err = cgroups.Load(cgroups.V1, cgroups.StaticPath("/test"))

	event := cgroups.MemoryThresholdEvent(uint64((limitBytes/100)*70), false)
	var efd uintptr
	efd, err = control.RegisterMemoryEvent(event)
	if err != nil {
		panic("RegisterMemoryEvent")
	}

	ch, err := watchMemThresholdEvent(efd)

	//control.AddTask()
	fmt.Printf("efd=%v \n", efd)

	cmd := exec.Command("/usr/bin/java",
		"-XX:NativeMemoryTracking=summary",
		"-XX:MaxDirectMemorySize=536870912",
		"-jar", "/home/labile/go-cgroup-tool/stress-java/target/stress-java-0.0.1-SNAPSHOT.jar")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()

	//cgroups.Process{
	//	//Subsystem: "",
	//	Pid:       cmd.Process.Pid,
	//	//Path:      "",
	//}

	//"/usr/bin/java -jar /home/labile/go-cgroup-tool/stress-java/target/stress-java-0.0.1-SNAPSHOT.jar"

	fmt.Printf("adding pid:%v to cgroup", cmd.Process.Pid)
	err = control.Add(cgroups.Process{
		Pid: cmd.Process.Pid,
	})
	if err != nil {
		panic("add cgroup error")
	}

	go func() {
		javaPid := cmd.Process.Pid
		select {
		case <-ch:
			fmt.Printf("Got ch\n")
			handleMemThresholdEvent(javaPid)
			//case <-time.After(100 * time.Millisecond):
			//	t.Fatal("no notification on channel after 100ms")
		}
	}()

	fmt.Printf("Waiting child process...")
	cmd.Wait()

	fmt.Printf("child process ExitCode:%v", cmd.ProcessState.ExitCode())
}

func watchMemThresholdEvent(efd uintptr) (<-chan struct{}, error) {
	eventfd := os.NewFile(efd, "eventfd")

	//eventControlPath := filepath.Join(cgDir, "cgroup.event_control")
	//data := fmt.Sprintf("%d %d %s", eventfd.Fd(), evFile.Fd(), arg)
	//if err := ioutil.WriteFile(eventControlPath, []byte(data), 0700); err != nil {
	//	eventfd.Close()
	//	evFile.Close()
	//	return nil, err
	//}
	ch := make(chan struct{})
	go func() {
		defer func() {
			eventfd.Close()
			//evFile.Close()
			close(ch)
		}()
		buf := make([]byte, 8)
		for {
			if _, err := eventfd.Read(buf); err != nil {
				fmt.Printf("eventfd.Read(buf) error: %v", err)
				return
			}
			// When a cgroup is destroyed, an event is sent to eventfd.
			// So if the control path is gone, return instead of notifying.
			//if _, err := os.Lstat(eventControlPath); os.IsNotExist(err) {
			//	return
			//}
			ch <- struct{}{}
			fmt.Printf("get cgroup mem event %v", buf)
		}
	}()
	return ch, nil
}

func handleMemThresholdEvent(javaPid int) {
	jcmdCmd := exec.Command("/usr/bin/jcmd", strconv.Itoa(javaPid), "VM.native_memory summary")
	jcmdCmd.Stdout = os.Stdout
	jcmdCmd.Stderr = os.Stderr
	jcmdCmd.Run()

	//time.Sleep(2*time.Second)
	syscall.Kill(javaPid, syscall.SIGSEGV)

	//coreDumpCmd := exec.Command("/usr/bin/kill", "-11", strconv.Itoa(javaPid))
	//fmt.Printf("kill -11 %v \n", javaPid)
	//coreDumpCmd.Stdout = os.Stdout
	//coreDumpCmd.Stderr = os.Stderr
	//coreDumpCmd.Run()
}
