package main

import (
	"fmt"

	"encoding/json"
	"strings"

	nan "nanocloud.com/core/lib/libnan"

	"github.com/dullgiulio/pingo" // for plugins
)

type Plugin struct {
	*pingo.Plugin
}

type PluginsMaps_t map[string]*Plugin

func SetupPlugins() {

	g_PluginsMap = make(PluginsMaps_t)

	Log("Num plugins referenced in config : %d", len(nan.Config().Plugins))

	for pluginName, _ := range nan.Config().Plugins {

		fmt.Println(pluginName)

		pluginNameLowercase := strings.ToLower(pluginName)

		pluginPath := fmt.Sprintf("%s/plugins/%s/%s", nan.Config().CommonBaseDir,
			pluginNameLowercase, pluginNameLowercase)

		pPlugin := pingo.NewPlugin("tcp", pluginPath)

		if pPlugin == nil {
			nan.ExitErrorf(0, "Failed to create plugin %s", pluginName)
		}

		pPlugin.Start()

		g_PluginsMap[pluginNameLowercase] = &Plugin{pPlugin}

		Log("Starting plugin %s", pluginName)

		ok := false
		var pluginJsonParams nan.PluginParams

		pluginRpcName := strings.ToUpper(pluginName[0:1]) + pluginName[1:len(pluginName)]

		// try finding plugin params in main config.json using lowercase name
		if pluginJsonParams, ok = nan.Config().Plugins[pluginNameLowercase]; !ok {

			// try finding plugin params in main config.json using Rpc "Pascal" name, first letter = uppercase)
			if pluginJsonParams, ok = nan.Config().Plugins[pluginRpcName]; !ok {

				err := nan.Errorf("Plugin %s doesn't have a parameters section in config.json !", pluginName)

				nan.ExitError(err)
			}
		}

		pluginParams, e := json.Marshal(pluginJsonParams)

		if e != nil {
			LogError("Failed to unmarshall %s plugin params", pluginName)
			ExitError(nan.ErrConfigError)
		}

		resp := ""

		if e := pPlugin.Call(pluginRpcName+".Configure", string(pluginParams), &resp); e != nil {
			// TODO Clarify error and string output
			LogError("while configuring plugin %s : %s", pluginRpcName, e)
			// TODO activate this line when all plugins have a Configure method
			ExitError(nan.ErrPluginError)
		}

		Log("Start plugin %s : DONE", pluginRpcName)

	}

}

func GetPlugin(pluginName string) (*Plugin, *nan.Err) {
	if pPlugin, ok := g_PluginsMap[pluginName]; !ok {
		err := nan.ErrPluginUnknown
		err.Message += (": " + pluginName)
		LogErrorCode(err)
		return nil, nan.ErrPluginUnknown
	} else {
		return pPlugin, nil
	}
}

func StopPlugins() {

	for _, pPlugin := range g_PluginsMap {
		if pPlugin != nil {
			pPlugin.Stop()
		}

		pPlugin = nil
	}
}
