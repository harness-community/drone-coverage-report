package plugin

import (
	"context"
	pd "github.com/harness-community/drone-coverage-report/plugin/plugin_defs"
	"testing"
)

func TestCoberturaGoodThreshold(t *testing.T) {

	envPluginInputArgs := pd.EnvPluginInputArgs{
		MinimumBranchCoverage:        0,
		MinimumClassCoverage:         30,
		MinimumLineCoverage:          5.0,
		MinimumMethodCoverage:        0,
		MinimumPackageCoverage:       100,
		MinimumFileCoverage:          50.0,
		MinimumLOC:                   15,
		MinimumComplexityCoverage:    10,
		MaxComplexityDensityCoverage: 0.50,
	}

	args := GetTestCoberturaNewArgs(envPluginInputArgs)
	plugin, err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("Expected passing threshold but got error: %s", err.Error())
	}
	_ = plugin
}

func TestCoberturaBadThreshold(t *testing.T) {

	envPluginInputArgs := pd.EnvPluginInputArgs{
		MinimumBranchCoverage:        0,
		MinimumClassCoverage:         40,
		MinimumLineCoverage:          5.0,
		MinimumMethodCoverage:        0,
		MinimumPackageCoverage:       100,
		MinimumFileCoverage:          50.0,
		MinimumLOC:                   15,
		MinimumComplexityCoverage:    10,
		MaxComplexityDensityCoverage: 0.50,
	}

	args := GetTestCoberturaNewArgs(envPluginInputArgs)
	plugin, err := Exec(context.TODO(), args)
	if err == nil {
		t.Errorf("Expected failure for high class coverage threshold but test passed: %s", err.Error())
	}
	_ = plugin
}

func TestCoberturaNoFailOnBadThreshold(t *testing.T) {

	envPluginInputArgs := pd.EnvPluginInputArgs{
		MinimumBranchCoverage:        0,
		MinimumClassCoverage:         40,
		MinimumLineCoverage:          5.0,
		MinimumMethodCoverage:        0,
		MinimumPackageCoverage:       100,
		MinimumFileCoverage:          50.0,
		MinimumLOC:                   15,
		MinimumComplexityCoverage:    10,
		MaxComplexityDensityCoverage: 0.50,
	}

	args := GetTestCoberturaNewArgs(envPluginInputArgs)
	args.PluginFailOnThreshold = false
	plugin, err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("Expected no error due to PluginFailOnThreshold being false, but got: %s", err.Error())
	}
	_ = plugin
}

func TestCoberturaBadComplexityCoverage(t *testing.T) {

	envPluginInputArgs := pd.EnvPluginInputArgs{
		MinimumBranchCoverage:        0,
		MinimumClassCoverage:         30,
		MinimumLineCoverage:          5.0,
		MinimumMethodCoverage:        0,
		MinimumPackageCoverage:       100,
		MinimumFileCoverage:          50.0,
		MinimumLOC:                   15,
		MinimumComplexityCoverage:    5,
		MaxComplexityDensityCoverage: 0.50,
	}

	args := GetTestCoberturaNewArgs(envPluginInputArgs)
	plugin, err := Exec(context.TODO(), args)
	if err == nil {
		t.Errorf("Error in TestCoberturaGoodThreshold: %s", err.Error())
	}
	_ = plugin
}

func TestCoberturaBadComplexityDensityCoverage(t *testing.T) {

	envPluginInputArgs := pd.EnvPluginInputArgs{
		MinimumBranchCoverage:        0,
		MinimumClassCoverage:         30,
		MinimumLineCoverage:          5.0,
		MinimumMethodCoverage:        0,
		MinimumPackageCoverage:       100,
		MinimumFileCoverage:          50.0,
		MinimumLOC:                   15,
		MinimumComplexityCoverage:    10,
		MaxComplexityDensityCoverage: 0.25,
	}

	args := GetTestCoberturaNewArgs(envPluginInputArgs)
	plugin, err := Exec(context.TODO(), args)
	if err == nil {
		t.Errorf("Error in TestCoberturaGoodThreshold: %s", err.Error())
	}
	_ = plugin
}

func GetTestCoberturaNewArgs(envPluginInputArgs pd.EnvPluginInputArgs) pd.Args {

	args := pd.Args{
		Pipeline: pd.Pipeline{},
		CoveragePluginArgs: pd.CoveragePluginArgs{
			PluginToolType:        pd.CoberturaPluginType,
			PluginFailOnThreshold: true,
		},
		EnvPluginInputArgs: envPluginInputArgs,
	}
	args.ExecFilesPathPattern = "**/coverage.xml"
	return args
}
