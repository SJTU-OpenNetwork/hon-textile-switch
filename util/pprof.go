package util

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	"time"
)

// A timed task used to record pprof log.
type PprofTask struct {
	cpuDir string
	memDir string
	cpuNotice chan interface{}
	memNotice chan interface{}
}

func NewPprofTask(cpuDir string, memDir string) *PprofTask {
	return &PprofTask{
		cpuDir:    cpuDir,
		memDir:    memDir,
		cpuNotice: make(chan interface{}, 1),
		memNotice: make(chan interface{}, 1),
	}
}

func (p *PprofTask)NoticeCpu() {
	select {
	case p.cpuNotice <- struct{}{}:
	default:
	}
}

func (p *PprofTask)NoticeMem() {
	select {
	case p.memNotice <- struct{}{}:
	default:
	}
}

func (p *PprofTask) StartCpu(interval time.Duration, withNotice bool, ctx context.Context) {
	timer := time.NewTimer(interval)

	var cpuPath string
	var cpuFile *os.File
	var err error
	openTime := time.Now()
	cpuPath = path.Join(p.cpuDir, "cpu_" + openTime.Format("15_04_05") + "_1min.prof")
	cpuFile, err = os.Create(cpuPath)
	if err != nil {
		fmt.Println("Error when create cpu prof file: ", err)
	} else {
		pprof.StartCPUProfile(cpuFile)
		//defer pprof.StopCPUProfile()
	}

	for {
		select{
		case <- timer.C:
			if withNotice {
				<-p.cpuNotice
			}

			timer.Reset(interval)
			pprof.StopCPUProfile()
			if cpuFile != nil {
				err = cpuFile.Close()
				if err != nil {
					fmt.Println("Error when close cpu prof file: ", err)
				}
			}

			openTime = time.Now()
			cpuPath = path.Join(p.cpuDir, "cpu_" + openTime.Format("15_04_05") + "_1min.prof")
			cpuFile, err = os.Create(cpuPath)
			if err != nil {
				fmt.Println("Error when create cpu prof file: ", err)
			} else {
				pprof.StartCPUProfile(cpuFile)
				//defer pprof.StopCPUProfile()
			}

		case <-ctx.Done():
			fmt.Println("Cpu prof task ctx end: ", ctx.Err())
		}
	}
}

// Start the memory statistic task
// Call this func in a independent go routine
func (p *PprofTask)StartMem(interval time.Duration, withNotice bool, ctx context.Context) {
	//timeCh := time.Tick(interval)
	timer := time.NewTimer(interval)

	var memPath string
	var memFile *os.File
	var err error
	for {
		select{
		case t :=<- timer.C:
			if withNotice {
				<-p.memNotice
			}
			timer.Reset(interval)
			memPath = path.Join(p.memDir, "mem_" + t.Format("15_04_05") + ".prof")
			memFile, err = os.Create(memPath)
			if err != nil {
				fmt.Println("Error when create memory prof file: ", err)
				continue
			}
			err = pprof.WriteHeapProfile(memFile)
			if err != nil {
				fmt.Println("Error when write to memory prof file: ", err)
			}
			err = memFile.Close()
			if err != nil {
				fmt.Println("Error when close memory prof file: ", err)
			}
		case <-ctx.Done():
			fmt.Println("Memory prof task ctx end: ", ctx.Err())
		}
	}
}