package main

import (
	"github.com/bitmagnet-io/bitmagnet/internal/app"
	_ "github.com/joho/godotenv/autoload"
	_ "net/http/pprof"
)

func main() {
	//fCpu, _ := os.Create("bitmagnet_cpu.prof")
	//_ = pprof.StartCPUProfile(fCpu)
	//defer pprof.StopCPUProfile()
	//fHeap, _ := os.Create("bitmagnet_heap.prof")
	//defer (func() {
	//	_ = pprof.WriteHeapProfile(fHeap)
	//})()
	app.New().Run()
}
