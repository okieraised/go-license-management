package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"go-license-management/internal/utils"
	"log"
)

type MachineAttribute struct {
	CPUModel string   `json:"cpu_model"`
	Platform string   `json:"platform"`
	MacAddr  []string `json:"mac_addr"`
	IPAddr   string   `json:"ip_addr"`
	Serial   string   `json:"serial"`
}

func main() {
	for _ = range 10 {
		// Get host info
		hostStat, err := host.Info()
		if err != nil {
			log.Fatal(err)
		}

		// Get cpu stat
		cpuStat, err := cpu.Info()
		if err != nil {
			log.Fatal(err)
		}

		// Get mac address
		macAddress, err := utils.RetrievePhysicalMacAddr()
		if err != nil {
			log.Fatal(err)
		}

		// Get IP address
		ip, err := utils.GetOutboundIP()
		if err != nil {
			log.Fatal(err)
		}
		motherboardSerial := hostStat.HostID
		cpuModel := cpuStat[0].ModelName
		platform := hostStat.Platform

		attr := MachineAttribute{
			CPUModel: cpuModel,
			Platform: platform,
			MacAddr:  macAddress,
			IPAddr:   ip.String(),
			Serial:   motherboardSerial,
		}
		bAttr, err := json.Marshal(attr)
		if err != nil {
			log.Fatal(err)
		}
		hash := sha256.Sum256(bAttr)
		fingerprint := hex.EncodeToString(hash[:])

		fmt.Println("fingerprint", fingerprint)
	}

}
