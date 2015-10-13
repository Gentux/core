package main

import (
	"fmt"
)

type VmListParams struct {
}

// ========================================================================================================================
// Procedure: ListVMs
//
// Does:
// - Return list of running VMs
// ========================================================================================================================
func ListVMs() string {

	var (
		vms string
		res string
	)

	err := g_PluginIaas.Call("Iaas.ListRunningVm", vms, &res)
	if err != nil {
		fmt.Println("Error calling Iaas Plugin")
	}

	return res
}
