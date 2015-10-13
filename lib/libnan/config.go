package libnan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	unsetstring   = "unsetstring"
	unsetint      = math.MinInt32
	unsetduration = math.MinInt64
)

type DatabaseConfig_t struct {
	Type             string
	ConnectionString string
}

type PluginParams map[string]string

type PluginsInfo_t map[string]PluginParams

type ProxyConfig_t struct {
	MaxNumRegistrations int
	MaxNumAccounts      int
	NumRetries          int
	SleepDurationInSecs time.Duration

	WinExe string
}

type Config_t struct {
	// Role  string `json:"Role",omitempty`
	// Debug bool   `json:"Debug",omitempty`

	// CommonBaseDir string `json:"CommonBaseDir",omitempty`
	// LogFilePath   string `json:"logfilepath"`

	// ConsulPath string `json:"consulpath",omitempty`

	// Proxy   ProxyConfig_t `json:"proxy",omitempty`
	// Plugins PluginsInfo_t `json:"plugins",omitempty`

	Role  string
	Debug bool

	CommonBaseDir string
	LogFilePath   string

	ConsulPath string

	Database DatabaseConfig_t
	Proxy    ProxyConfig_t
	Plugins  PluginsInfo_t
}

var (
	g_Config Config_t
)

func Config() Config_t {
	return g_Config
}

// Looks for config file confing.json in current directory
func LoadConfig() {

	sConfigFilePath := os.Getenv("NANOCONF")

	if sConfigFilePath == "" {
		exeDir, _ := filepath.Abs(path.Dir(os.Args[0]))
		sConfigFilePath = exeDir + "/config.json"
	}

	if sDryRunConf := os.Getenv("NANODRYRUN"); sDryRunConf == "1" {
		DryRun = true
	}

	if sRef := os.Getenv("NANOREF"); sRef == "1" {
		ModeRef = true
	}

	initConfigWithUnsetValues()

	ok := true

	fileBytes, err := ioutil.ReadFile(sConfigFilePath)

	if err != nil {
		fmt.Println("Could not read from config file located at: ", sConfigFilePath)
		ok = false
	} else if err := json.Unmarshal(fileBytes, &g_Config); err != nil {
		fmt.Println("Failed to read config file, verify syntax in: ", sConfigFilePath)
		ok = false
	}

	if !ConfigFileValid() {
		LogError("Exiting program because of parameter missing from config.json")
		os.Exit(-1)
	}

	if !ok {
		os.Exit(-1)
	}
}

func initConfigWithUnsetValues() {

	g_Config = Config_t{
		Role: unsetstring,

		Debug: false,

		CommonBaseDir: unsetstring,
		LogFilePath:   unsetstring,
		ConsulPath:    unsetstring,

		Proxy: ProxyConfig_t{
			MaxNumRegistrations: unsetint,
			MaxNumAccounts:      unsetint,
			NumRetries:          unsetint,
			SleepDurationInSecs: unsetduration,
			WinExe:              unsetstring}}
}

func ConfigFileValid() bool {
	if g_Config.Role == unsetstring {
		fmt.Println(`Missing config param "Role", expected one of : [ "proxy", "plugin", "tac" ]`)
		return false
	}

	if g_Config.LogFilePath == unsetstring {
		fmt.Println(`Missing config param : "LogFilePath"`)
		return false
	}

	if g_Config.Role == "proxy" || g_Config.Role == "tac" {
		if g_Config.CommonBaseDir == unsetstring {
			fmt.Println(`Missing config param : "CommonBaseDir"`)
			return false
		}
	}

	if g_Config.Role == "proxy" {

		if g_Config.Proxy.MaxNumRegistrations == unsetint {
			fmt.Println(`Missing config param : "Proxy" : { "MaxNumRegistrations" : x }`)
			return false
		}

		if g_Config.Proxy.MaxNumAccounts == unsetint {
			fmt.Println(`Missing config param : "Proxy" : { "MaxNumAccounts" : x }`)
			return false
		}

		if g_Config.Proxy.NumRetries == unsetint {
			fmt.Println(`Missing config param : "Proxy" : { "NumRetries" : x }`)
			return false
		}

		if g_Config.Proxy.SleepDurationInSecs == unsetint {
			fmt.Println(`Missing config param : "Proxy" : { "SleepDurationInSecs" : x }`)
			return false
		}

	}

	//TODO Optional, only for plugin Talend
	// if g_Config.ConsulPath == unsetstring {
	// 	fmt.Println(`Missing config param : "ConsulPath"`)
	// 	return false
	// }

	return true
}
