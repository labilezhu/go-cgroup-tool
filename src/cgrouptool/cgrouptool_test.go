package main

import (
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"testing"
)

func TestAbs(t *testing.T) {
	//oldArgs := os.Args
	//defer func() { os.Args = oldArgs }()

	//-level 500:"/usr/bin/kill -11 97632" -level 400:"/usr/bin/jcmd 97632 VM.native_memory summary"

	var limitBytes int64 = 512 * 1024 * 1024
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

	err = control.Add(cgroups.Process{
		Pid: os.Getpid(),
	})

	os.Args = []string{"cgrouptool",
		"-memCGroupPath",
		cgroupName,
		"-level",
		"500:/usr/bin/kill -11 97632",
		"-level",
		"400:/usr/bin/jcmd 97632 VM.native_memory summary",
	}

	main()
}
