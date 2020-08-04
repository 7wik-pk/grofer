/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package general

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func roundOff(num uint64) float64 {
	x := float64(num) / (1024 * 1024 * 1024)
	return math.Round(x*10) / 10
}

// PrintCPURates print the cpu rates
func PrintCPURates(cpuChannel chan []float64) {
	cpuRates, err := cpu.Percent(time.Second, true)
	if err != nil {
		log.Fatal(err)
	}
	cpuChannel <- cpuRates
}

// PrintMemRates prints stats about the memory
func PrintMemRates(dataChannel chan []float64) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	data := []float64{roundOff(memory.Total), roundOff(memory.Available), roundOff(memory.Used), roundOff(memory.Free)}
	dataChannel <- data
}

func PrintDiskRates(dataChannel chan [][]string) {

	var partitions []disk.PartitionStat
	var err error
	partitions, err = disk.Partitions(false)
	if err != nil {
		log.Fatal(err)
	}

	rows := [][]string{[]string{"Mount", "Total", "Used %", "Used", "Free", "FS Type"}}
	for _, value := range partitions {
		usageVals, _ := disk.Usage(value.Mountpoint)

		if strings.HasPrefix(value.Device, "/dev/loop") {
			continue
		} else if strings.HasPrefix(value.Mountpoint, "/var/lib/docker") {
			continue
		} else {

			path := usageVals.Path
			total := fmt.Sprintf("%.2f G", float64(usageVals.Total)/(1024*1024*1024))
			used := fmt.Sprintf("%.2f G", float64(usageVals.Used)/(1024*1024*1024))
			usedPercent := fmt.Sprintf("%.2f %s", usageVals.UsedPercent, "%")
			free := fmt.Sprintf("%.2f G", float64(usageVals.Free)/(1024*1024*1024))
			fs := usageVals.Fstype
			row := []string{path, total, usedPercent, used, free, fs}
			rows = append(rows, row)

		}
	}
	dataChannel <- rows
}

func PrintNetRates(dataChannel chan map[string][]float64) {
	netStats, err := net.IOCounters(false)
	if err != nil {
		log.Fatal(err)
	}
	IO := make(map[string][]float64)
	for _, IOStat := range netStats {
		nic := []float64{float64(IOStat.BytesSent) / (1024), float64(IOStat.BytesRecv) / (1024)}
		IO[IOStat.Name] = nic
	}
	dataChannel <- IO
}
