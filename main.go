package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CPUInfo struct {
	ModelName string
	Cores     string
}
type RAMInfo struct {
	Total     string
	Used      string
	Available string
}

// Mem = /proc/meminfo
// CPU = /proc/cpuinfo

func main() {
	// Get raw content of os-release file
	osReleaseContent, err := os.ReadFile("/etc/os-release")
	if err != nil {
		log.Fatal("Can't read os-release (/etc/os-release)")
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
	// Username
	username := os.Getenv("USER")
	// Hostname
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Can't get hostname. Please check your /etc/hostname file.")
	}

	fmt.Printf("\n💾 ∙ OS:       %s\n🔧 ∙ CPU:      %s | %s cores\n🧠 ∙ RAM:      %s/%s\n🧑‍💻 ∙ User:     %s\n🏠 ∙ Hostname: %s\n\n", osRelease[1], cpuInfo.ModelName, cpuInfo.Cores, memInfo.Used, memInfo.Total, username, hostname)
}
