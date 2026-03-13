package main

import (
	"os"
	"runtime/pprof"
)

func wasteCPU(n int) int {
	if n <= 1 {
		return n
	}
	return wasteCPU(n-1) + wasteCPU(n-2)
}

func wasteMemory() [][]byte {
	var data [][]byte
	for i := 0; i < 1000; i++ {
		b := make([]byte, 1024*1000)
		data = append(data, b)
	}
	return data
}

func main() {
	// CPU Profiling
	cpuFile, _ := os.Create("out/cpu.pprof")
	pprof.StartCPUProfile(cpuFile)

	wasteCPU(40)
	wasteMemory()

	pprof.StopCPUProfile()

	// Memory Profiling
	memFile, _ := os.Create("out/mem.pprof")
	pprof.WriteHeapProfile(memFile)

}
