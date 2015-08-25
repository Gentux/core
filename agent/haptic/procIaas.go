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

// ========================================================================================================================
// Procedure: DownloadWindowsVm
//
// Does:
// - Launch windows download from hypervisor by calling Iaas API
// ========================================================================================================================
func DownloadWindowsVm() bool {

	var (
		// TODO: We mustn't have value hardcoded like that
		vmName string = "winad-milli-free_use-10.20.12.20-windows-server-std-2012-x86_64"
		res    string
	)

	err := g_PluginIaas.Call("Iaas.DownloadVm", vmName, &res)
	if err != nil {
		fmt.Println("Error calling Iaas Plugin")
	}

	if res == "true" {
		return true
	} else {
		return false
	}
}

// ========================================================================================================================
// Procedure: StartVm
//
// Does:
// - Start a virtual machine matching vmName
// ========================================================================================================================
func StartVm(vmName string) bool {

	var (
		res string
	)

	err := g_PluginIaas.Call("Iaas.StartVm", vmName, &res)
	if err != nil {
		fmt.Printf("Error calling Iaas Plugin: %v", err)
	}

	if res == "true" {
		return true
	} else {
		return false
	}
}
