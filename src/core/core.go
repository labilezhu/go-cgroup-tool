package core

import (
	"fmt"
	"github.com/containerd/cgroups"
	"os"
	"os/exec"
)

func Core(levels []Level, memCGroupPath string) {
	fmt.Printf("cgroups.load path=%v\n", memCGroupPath)
	control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(memCGroupPath))
	if err != nil {
		panic(fmt.Sprintf("cgroups.Load:%s failed", memCGroupPath))
	}

	fmt.Printf("cgroup subsys: %v \n", control.Subsystems())

	metrics, err := control.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		panic("control.Stat() error")
	}

	fmt.Printf("cgroup limitMem Megabytes: %v \n", metrics.Memory.HierarchicalMemoryLimit)

	processes, err := control.Processes(cgroups.Memory, false)
	if err != nil {
		panic(fmt.Sprintf("list Processes of cgroup fail\n"))
	}

	fmt.Print("List Processes of cgroup:\n")
	for _, process := range processes {
		fmt.Printf("pid=%v path=%v, subsys=%v \n", process.Pid, process.Pid, process.Subsystem)
	}

	for _, level := range levels {
		watchLevel(control, level)
	}
}

func watchLevel(control cgroups.Cgroup, level Level) {
	fmt.Printf("watchLevel %v:%v \n", level.MemoryThresholdMegabytes, level.Command)

	event := cgroups.MemoryThresholdEvent(uint64(level.MemoryThresholdMegabytes*1024*1024), false)
	var efd uintptr
	efd, err := control.RegisterMemoryEvent(event)
	if err != nil {
		panic(fmt.Sprintf("RegisterMemoryEvent error, %v \n", err))
	}

	ch, err := watchMemThresholdEvent(efd)

	//control.AddTask()
	fmt.Printf("efd=%v \n", efd)

	go func() {
		select {
		case <-ch:
			fmt.Printf("***Got ch***\n")
			handleMemThresholdEvent(level)
			//case <-time.After(100 * time.Millisecond):
			//	t.Fatal("no notification on channel after 100ms")
		}
	}()
}

func watchMemThresholdEvent(efd uintptr) (<-chan struct{}, error) {
	eventfd := os.NewFile(efd, fmt.Sprintf("eventfd_%v", efd))

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

			fmt.Printf("Reading efd %v \n", efd)

			if _, err := eventfd.Read(buf); err != nil {
				fmt.Printf("eventfd.Read(buf) error: %v", err)
				return
			}

			fmt.Printf("Read return efd %v \n", efd)
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

func handleMemThresholdEvent(level Level) {
	fmt.Printf("handleMemThresholdEvent %v:%v \n", level.MemoryThresholdMegabytes, level.Command)

	path, err := exec.LookPath("bash")
	if err != nil {
		fmt.Print("exec.LookPath(\"bash\") fail")
		path = "/bin/bash"
	}

	fmt.Printf("exec %v -c %v", path, level.Command)

	jcmdCmd := exec.Command(path, "-c", level.Command)
	jcmdCmd.Stdout = os.Stdout
	jcmdCmd.Stderr = os.Stderr
	jcmdCmd.Run()
}

type Level struct {
	MemoryThresholdMegabytes int
	Command                  string
}
