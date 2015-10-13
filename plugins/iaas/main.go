package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dullgiulio/pingo"

	nan "nanocloud.com/zeroinstall/lib/libnan"
)

type IaasConfig struct {
	Url  string
	Port string
}

type Iaas struct{}

type VmInfo struct {
	Id   string
	Ico  string
	Name string
	VM   string
}

var (
	g_IaasConfig IaasConfig
)

func (p *Iaas) Configure(jsonConfig string, _outMsg *string) error {
	var iaasConfig map[string]string

	err := json.Unmarshal([]byte(jsonConfig), &iaasConfig)
	if err != nil {
		r := fmt.Sprintf("ERROR: failed to unmarshal Iaas Plugin configuration : %s", err.Error())
		log.Printf(r)
		os.Exit(0)
		*_outMsg = r
		return nil
	}

	g_IaasConfig.Url = iaasConfig["url"]
	g_IaasConfig.Port = iaasConfig["port"]

	return nil
}

func (p *Iaas) ListRunningVm(jsonParams string, _outMsg *string) error {
	var fakeData [2]VmInfo

	fakeData[0] = VmInfo{
		Id:   "2",
		Ico:  "windows",
		Name: "Test RDP",
		VM:   "winad",
	}

	fakeData[1] = VmInfo{
		Id:   "1",
		Ico:  "view_module",
		Name: "Haptic Test",
		VM:   "proxy",
	}

	res, err := json.Marshal(fakeData)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to marshal VM list : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	*_outMsg = string(res)
	return nil
}

func (p *Iaas) DownloadVm(jsonParams string, _outMsg *string) error {

	var (
		params map[string]string
		vmName string
	)

	if err := json.Unmarshal([]byte(jsonParams), &vmName); err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas.AccountParams : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	params["vmname"] = vmName

	jsonRpcRequest(
		g_IaasConfig.Url,
		"Iaas.Download",
		params,
	)

	return nil
}

func (p *Iaas) StartVm(jsonParams string, _outMsg *string) error {

	var (
		params map[string]string
		vmName string
	)

	if err := json.Unmarshal([]byte(jsonParams), &vmName); err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas.AccountParams : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	params["name"] = vmName

	jsonRpcRequest(
		g_IaasConfig.Url,
		"Iaas.Start",
		params,
	)

	return nil
}

func (p *Iaas) StopVm(jsonParams string, _outMsg *string) error {

	var (
		params map[string]string
		vmName string
	)

	if err := json.Unmarshal([]byte(jsonParams), &vmName); err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas.AccountParams : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	params["name"] = vmName

	jsonRpcRequest(
		g_IaasConfig.Url,
		"Iaas.Stop",
		params,
	)

	return nil
}

func jsonRpcRequest(url string, method string, params map[string]string) {
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      1,
		"params": []map[string]string{
			0: params,
		},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}
	result := make(map[string]interface{})
	// TODO Check result and do something about it
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func main() {

	plugin := &Iaas{}

	pingo.Register(plugin)

	pingo.Run()
}
