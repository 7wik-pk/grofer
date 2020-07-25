package general

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// GlobalStats gets stats about the mem and the CPUs and prints it.
func GlobalStats(endChannel chan os.Signal, memChannel chan []float64, cpuChannel chan []float64, wg *sync.WaitGroup) {
	for {
		select {
		case <-endChannel: // Stop execution if end signal received
			wg.Done()
			return

		default: // Get Memory and CPU rates per core periodically

			cpuUsageRates, err := cpu.Percent(time.Second, true)
			if err != nil {
				log.Fatal(err)
			}

			memoryStat, err := mem.VirtualMemory()
			if err != nil {
				log.Fatal(err)
			}

			go PrintCPURates(cpuUsageRates, cpuChannel)
			go PrintMemRates(memoryStat, memChannel)
		}
	}
}
