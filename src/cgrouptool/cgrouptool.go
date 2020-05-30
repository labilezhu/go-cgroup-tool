package main

import (
	"flag"
	"github.com/labilezhu/go-cgroup-tool/core"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	var levels arrayFlags
	var help bool
	//var targetPid    int
	var memCGroupPath string

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	//flagSet.IntVar(&targetPid, "targetPid", 0, "The pid of target process")
	flagSet.Var(&levels, "level", "format is : $MemoryThresholdMegabytes:command . E.g "+
		"-level 500:\"/usr/bin/kill -11 97632\" "+
		"-level 400:\"/usr/bin/jcmd 97632 VM.native_memory summary\" ")

	flagSet.StringVar(&memCGroupPath, "memCGroupPath", "/",
		"The path related to `/sys/fs/cgroup/` prefix. All static path should not include `/sys/fs/cgroup/` prefix. When run in conatiner, should be / . Default / ")

	flagSet.Parse(os.Args[1:])
	if help || len(levels) == 0 {
		flagSet.Usage()
		os.Exit(-1)
	}

	var parsedLevels []core.Level

	levelPattern, _ := regexp.Compile(`(.+):(.+)`)

	//for k, v := range result {
	//	result:= re1.FindStringSubmatch(s)
	//	fmt.Printf("%d. %s\n", k, v)
	//}
	//
	for _, level := range levels {
		matchArray := levelPattern.FindStringSubmatch(level)

		memoryThresholdMegabytes, err := strconv.Atoi(matchArray[1])
		if err != nil {
			panic("parse memoryThresholdMegabytes error")
		}

		command := matchArray[2]
		parsedLevels = append(parsedLevels, core.Level{
			MemoryThresholdMegabytes: memoryThresholdMegabytes,
			Command:                  command,
		})
	}

	core.Core(parsedLevels, memCGroupPath)

	time.Sleep(time.Hour)
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
