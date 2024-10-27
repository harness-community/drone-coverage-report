package plugin

import (
	"context"
	"fmt"
	pd "github.com/harness-community/drone-coverage-report/plugin/plugin_defs"
	"testing"
)

func TestJacocoXmlGoodThreshold(t *testing.T) {

	fmt.Println("Running TestJacocoXml")

	reportsPathPattern := "**/jacoco.xml"

	envPluginInputArgs := pd.EnvPluginInputArgs{
		ExecFilesPathPattern:       TestBuildRootPath,
		MinimumInstructionCoverage: 20,
		MinimumBranchCoverage:      20,
		MinimumComplexityCoverage:  75,
		MinimumLineCoverage:        20,
		MinimumMethodCoverage:      20,
		MinimumClassCoverage:       20,
	}

	args := GetTestJacocoXmlNewArgs(envPluginInputArgs)
	args.ExecFilesPathPattern = reportsPathPattern

	plugin, err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("Error in TestJcXml: %s", err.Error())
	}
	_ = plugin
}

func TestJacocoXmlBadMinimumThreshold(t *testing.T) {

	fmt.Println("Running TestJacocoXml")

	reportsPathPattern := "**/jacoco.xml"

	envPluginInputArgs := pd.EnvPluginInputArgs{
		ExecFilesPathPattern:       TestBuildRootPath,
		MinimumInstructionCoverage: 100,
		MinimumBranchCoverage:      20,
		MinimumComplexityCoverage:  75,
		MinimumLineCoverage:        20,
		MinimumMethodCoverage:      20,
		MinimumClassCoverage:       20,
	}

	args := GetTestJacocoXmlNewArgs(envPluginInputArgs)
	args.ExecFilesPathPattern = reportsPathPattern

	plugin, err := Exec(context.TODO(), args)
	if err == nil {
		t.Errorf("Error: TestJacocoXmlBadMinimumThreshold " +
			"MinimumInstructionCoverage check should have failed")
	}
	_ = plugin
}

func TestJacocoXmlBadComplexityThreshold(t *testing.T) {

	fmt.Println("Running TestJacocoXml")

	reportsPathPattern := "**/jacoco.xml"

	envPluginInputArgs := pd.EnvPluginInputArgs{
		ExecFilesPathPattern:       TestBuildRootPath,
		MinimumInstructionCoverage: 20,
		MinimumBranchCoverage:      20,
		MinimumComplexityCoverage:  60,
		MinimumLineCoverage:        20,
		MinimumMethodCoverage:      20,
		MinimumClassCoverage:       20,
	}

	args := GetTestJacocoXmlNewArgs(envPluginInputArgs)
	args.ExecFilesPathPattern = reportsPathPattern

	plugin, err := Exec(context.TODO(), args)
	if err == nil {
		t.Errorf("Error in TestJacocoXmlBadComplexityThreshold" +
			"MinimumComplexityCoverage check should have failed")
	}

	_ = plugin
}

func GetTestJacocoXmlNewArgs(envPluginInputArgs pd.EnvPluginInputArgs) pd.Args {

	args := pd.Args{
		Pipeline:           pd.Pipeline{},
		CoveragePluginArgs: pd.CoveragePluginArgs{PluginToolType: pd.JacocoXmlPluginType},
		EnvPluginInputArgs: envPluginInputArgs,
	}
	args.ExecFilesPathPattern = TestExecPathPattern01
	return args
}
