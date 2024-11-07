package main

import (
	"fmt"
	"go-license-management/internal/utils"
)

func main() {
	allMAC, err := utils.RetrieveMacAddr()
	if err != nil {
		fmt.Println(err)
		return
	}

	physicalMAC, err := utils.RetrievePhysicalMacAddr()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(allMAC)
	fmt.Println(physicalMAC)
}
