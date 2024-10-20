package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type DeviceStats struct {
	readCount  uint64
	writeCount uint64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <device_name1> [device_name2] ...")
		os.Exit(1)
	}
	deviceNames := os.Args[1:]

	lastStats := make(map[string]DeviceStats)
	for _, device := range deviceNames {
		lastStats[device] = DeviceStats{}
	}

	for {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		stats, err := getDiskStats(deviceNames)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("[%s]", timestamp)

		for device, currentStats := range stats {
			lastStat := lastStats[device]
			fmt.Printf(" %s (r: %5d w: %5d)",
				device,
				currentStats.readCount-lastStat.readCount,
				currentStats.writeCount-lastStat.writeCount)
			lastStats[device] = currentStats
		}
		fmt.Printf("\n")

		time.Sleep(5 * time.Second)
	}
}

func getDiskStats(deviceNames []string) (map[string]DeviceStats, error) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := make(map[string]DeviceStats)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		if contains(deviceNames, fields[2]) {
			readCount, err := strconv.ParseUint(fields[3], 10, 64)
			if err != nil {
				return nil, err
			}
			writeCount, err := strconv.ParseUint(fields[7], 10, 64)
			if err != nil {
				return nil, err
			}
			stats[fields[2]] = DeviceStats{readCount: readCount, writeCount: writeCount}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(stats) != len(deviceNames) {
		return nil, fmt.Errorf("some devices were not found")
	}

	return stats, nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
