/*
 * Nanocloud community -- transform any application into SaaS solution
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
// Procedure: DownloadStatus
//
// Does:
// - Return true if a download is in progress
// ========================================================================================================================
func DownloadStatus() bool {

	var res string

	err := g_PluginIaas.Call("Iaas.DownloadStatus", "", &res)
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
