//+build profile

package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var (
	profileHeap = flag.String("heapprof", "", "if non-empty, the filename for a pprof heap profile")
	profileCPU  = flag.String("cpuprof", "", "if non-empty, the filename for a pprof CPU profile")

	profileLock = make(chan struct{}, 1)
)

func init() {
	profileLock <- struct{}{}
	startup = profileStartup
	cleanup = profileCleanup
}

func profileStartup() {
	if *profileHeap != "" {
		f, err := os.Create(*profileHeap)
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for {
				time.Sleep(time.Second * 5)

				l := <-profileLock
				f.Seek(0, os.SEEK_SET)
				f.Truncate(0)
				pprof.Lookup("heap").WriteTo(f, 0)
				profileLock <- l
			}
		}()
	}

	if *profileCPU != "" {
		f, err := os.Create(*profileCPU)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func profileCleanup() {
	<-profileLock
	if *profileCPU != "" {
		pprof.StopCPUProfile()
	}
}
