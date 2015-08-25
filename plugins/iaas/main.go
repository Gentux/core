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

	g_IaasConfig.Url = iaasConfig["Url"]
	g_IaasConfig.Port = iaasConfig["Port"]

	return nil
}

type ListVMReply struct {
	AvailableVmList []VmInfo
	RunningVmList   []VmInfo
}

func (p *Iaas) ListRunningVm(jsonParams string, _outMsg *string) error {
	var err error
	*_outMsg, err = jsonRpcRequest(
		fmt.Sprintf("%s:%s", g_IaasConfig.Url, g_IaasConfig.Port),
		"Iaas.GetList",
		nil,
	)

	return err
}

func (p *Iaas) DownloadVm(jsonParams string, _outMsg *string) error {

	var (
		params map[string]string
		vmName string
	)

	err := json.Unmarshal([]byte(jsonParams), &vmName)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas.AccountParams : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	params["vmname"] = vmName

	go jsonRpcRequest(
		g_IaasConfig.Url,
		"Iaas.Download",
		params,
	)

	*_outMsg = "success"
	return err
}

func (p *Iaas) StartVm(jsonParams string, _outMsg *string) error {

	var (
		params map[string]string
		vmName string
	)

	err := json.Unmarshal([]byte(jsonParams), &vmName)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas.AccountParams : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	params["vmname"] = vmName

	*_outMsg, err = jsonRpcRequest(
		g_IaasConfig.Url,
		"Iaas.Start",
		params,
	)

	return err
}

func (p *Iaas) StopVm(jsonParams string, _outMsg *string) error {

	var (
		params map[string]string
		vmName string
	)

	err := json.Unmarshal([]byte(jsonParams), &vmName)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas.AccountParams : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	params["vmname"] = vmName

	*_outMsg, err = jsonRpcRequest(
		g_IaasConfig.Url,
		"Iaas.Stop",
		params,
	)

	return err
}

func jsonRpcRequest(url string, method string, param map[string]string) (string, error) {

	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      1,
		"params":  []map[string]string{0: param},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return "", err
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
		return "", err
	}

	return string(body), nil
}

func main() {

	plugin := &Iaas{}

	pingo.Register(plugin)

	pingo.Run()
}
