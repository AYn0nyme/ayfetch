package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

type CPUInfo struct {
	ModelName string
	Cores     string
}
type RAMInfo struct {
	Total     string
	Used      string
	Available string
	Percent   float64
}

type DiskInfo struct {
	Total    uint64
	Used     uint64
	MainPart string
	Percent  float64
}

// Mem = /proc/meminfo
// CPU = /proc/cpuinfo

func main() {
	// Get raw content of os-release file
	osReleaseContent, err := os.ReadFile("/etc/os-release")
	if err != nil {
		log.Fatal("Can't read os-release (/etc/os-release)")
	}
	// Get raw content of /proc/version
	kernelContent, err := os.ReadFile("/proc/version")
	if err != nil {
		log.Fatal("Can't read /proc/version")
	}
	// Get raw content of /proc/cpuinfo (CPU Infos)
	cpuContent, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		log.Fatal("Can't read /proc/cpuinfo")
	}

	// Get raw content of /proc/meminfo (Memory (RAM) Infos)
	memContent, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		log.Fatal("Can't read /proc/meminfo")
	}
	// Regex that only takes element wich are between ""
	reg := regexp.MustCompile(`"([^"]+)"`)
	// OS Release (e.g. Arch Linux)
	osRelease := reg.FindStringSubmatch(strings.Split(string(osReleaseContent), "\n")[1])
	// Kernel (e.g. Linux 6.18.13-arch1-1)
	kernel := strings.Fields(string(kernelContent))[0]
	kernelVersion := regexp.MustCompile(`version ([^\s]+)`).FindStringSubmatch(string(kernelContent))[1]
	// CPU Info

	regModel := regexp.MustCompile(`(?m)^model name\s+:\s+(.+)`)
	regCores := regexp.MustCompile(`(?m)^cpu cores\s+:\s+(\d+)`)

	cpuInfo := CPUInfo{
		ModelName: regModel.FindStringSubmatch(string(cpuContent))[1],
		Cores:     regCores.FindStringSubmatch(string(cpuContent))[1],
	}
	memInfo := RAMInfo{
		Total:     regexp.MustCompile(`(?m)^MemTotal:\s+(\d+)`).FindStringSubmatch(string(memContent))[1],
		Available: regexp.MustCompile(`(?m)^MemAvailable:\s+(\d+)`).FindStringSubmatch(string(memContent))[1],
	}
	totalMem, err := strconv.Atoi(memInfo.Total)
	avMem, err := strconv.Atoi(memInfo.Available)
	usedMem := totalMem - avMem
	memInfo.Used = fmt.Sprintf("%d MB", usedMem/1024)
	memInfo.Total = fmt.Sprintf("%d MB", totalMem/1024)
	memInfo.Available = fmt.Sprintf("%d MB", avMem/1024)
	memInfo.Percent = (float64(usedMem) / float64(totalMem)) * 100

	// Get disk space
	partitions, err := disk.Partitions(true)
	if err != nil {
		log.Fatal("Can't access partitions.")
	}
	var diskContent DiskInfo
	for _, partition := range partitions {
		if partition.Mountpoint == "/" {
			diskContent.MainPart = partition.Mountpoint
			diskSpace, err := disk.Usage(diskContent.MainPart)
			if err != nil {
				log.Fatal("Can't read disks.")
			}
			diskContent.Total = diskSpace.Total / (1024 * 1024)
			diskContent.Used = diskSpace.Used / (1024 * 1024)
			diskContent.Percent = diskSpace.UsedPercent
			break
		}

	}

	// Username
	username := os.Getenv("USER")
	// Hostname
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Can't get hostname. Please check your /etc/hostname file.")
	}

	fmt.Printf("\n💾 ∙ OS:       %s\n🐧 ∙ Kernel    %s %s\n🔧 ∙ CPU:      %s | %s cores\n🧠 ∙ RAM:      %s/%s (%.1f%%)\n💾 ∙ DISK (%s): %d MiB/%d MiB Used (%.1f%%)\n🧑 ∙ User:     %s\n🏠 ∙ Hostname: %s\n\n", osRelease[1], kernel, kernelVersion, cpuInfo.ModelName, cpuInfo.Cores, memInfo.Used, memInfo.Total, memInfo.Percent, diskContent.MainPart, diskContent.Used, diskContent.Total, diskContent.Percent, username, hostname)
}
