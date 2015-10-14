/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
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

	nan "nanocloud.com/core/lib/libnan"
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

	pPluginIaas, err := GetPlugin("iaas")
	if err != nil {
		return ""
	}

	e := pPluginIaas.Call("Iaas.ListRunningVm", vms, &res)
	if e != nil {
		LogError("when calling Iaas.ListRunningVm")
	}

	return res
}

// ========================================================================================================================
// Procedure: DownloadWindowsVm
//
// Does:
// - Launch windows download from hypervisor by calling Iaas API
// ========================================================================================================================
func DownloadWindowsVm() (bool, *nan.Err) {

	var (
		// TODO: We mustn't have value hardcoded like that
		vmName string = "winad-milli-free_use-10.20.12.20-windows-server-std-2012-x86_64"
		res    string
	)

	pPluginIaas, err := GetPlugin("iaas")
	if err != nil {
		return false, err
	}

	if e := pPluginIaas.Call("Iaas.DownloadVm", vmName, &res); e != nil {
		fmt.Println("Error calling Iaas Plugin")
		return false, nan.ErrFrom(e)
	}

	if res == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

// ========================================================================================================================
// Procedure: DownloadStatus
//
// Does:
// - Return true if a download is in progress
// ========================================================================================================================
func DownloadStatus() (bool, *nan.Err) {

	var res string

	pPluginIaas, err := GetPlugin("iaas")
	if err != nil {
		return false, err
	}

	if e := pPluginIaas.Call("Iaas.DownloadStatus", "", &res); e != nil {
		fmt.Println("Error calling Iaas Plugin")
		return false, nan.ErrFrom(e)
	}

	if res == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

// ========================================================================================================================
// Procedure: StartVm
//
// Does:
// - Start a virtual machine matching vmName
// ========================================================================================================================
func StartVm(vmName string) (bool, *nan.Err) {

	var res string

	pPluginIaas, err := GetPlugin("iaas")
	if err != nil {
		return false, err
	}

	e := pPluginIaas.Call("Iaas.StartVm", vmName, &res)
	if e != nil {
		LogError("when calling Iaas Plugin: %v", e)
		return false, nan.ErrFrom(e)
	}

	if res == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

// ========================================================================================================================
// Procedure: StopVm
//
// Does:
// - Stop a virtual machine matching vmName
// ========================================================================================================================
func StopVm(vmName string) (bool, *nan.Err) {

	var (
		res string
		err *nan.Err
	)

	var pPluginIaas *Plugin
	pPluginIaas, err = GetPlugin("iaas")
	if err != nil || pPluginIaas == nil {
		return false, err
	}

	if e := pPluginIaas.Call("Iaas.StopVm", vmName, &res); e != nil {
		LogError("when calling Iaas Plugin: %v", e)
	}

	if res == "true" {
		return true, nil
	} else {
		return false, nil
	}
}
