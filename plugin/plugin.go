// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	jc "github.com/harness-community/drone-coverage-report/plugin/jacoco"
	pd "github.com/harness-community/drone-coverage-report/plugin/plugin_defs"
)

func GetNewPlugin(ctx context.Context, args pd.Args) (pd.Plugin, error) {

	pluginToolType := args.PluginToolType

	switch pluginToolType {
	case pd.JacocoPluginType:
		jcp := jc.GetNewJacocoPlugin()
		return &jcp, nil
	case pd.JacocoXmlPluginType:
		jcxp := jc.GetNewJacocoXmlPlugin()
		return &jcxp, nil

	default:
		return nil, pd.GetNewError("Unknown plugin type: " + pluginToolType)
	}
}

func Exec(ctx context.Context, args pd.Args) (pd.Plugin, error) {

	plugin, err := GetNewPlugin(ctx, args)
	if err != nil {
		return plugin, err
	}

	err = plugin.Init(&args)
	if err != nil {
		return plugin, err
	}
	defer func(p pd.Plugin) {
		err := p.DeInit()
		if err != nil {
			pd.LogPrintln(p, "Error in DeInit: "+err.Error())
		}
	}(plugin)

	err = plugin.ValidateAndProcessArgs(args)
	if err != nil {
		return plugin, err
	}

	err = plugin.DoPostArgsValidationSetup(args)
	if err != nil {
		return plugin, err
	}

	err = plugin.Run()
	if err != nil {
		return plugin, err
	}

	err = plugin.PersistResults()
	if err != nil {
		return plugin, err
	}

	err = plugin.WriteOutputVariables()
	if err != nil {
		return plugin, err
	}

	return plugin, nil
}

//
//
