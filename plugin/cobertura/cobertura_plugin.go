package cobertura

import (
	"fmt"
	"github.com/bmatcuk/doublestar/v4"
	pd "github.com/harness-community/drone-coverage-report/plugin/plugin_defs"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type CoberturaPlugin struct {
	pd.CoveragePluginArgs
	InputArgs *pd.Args
	CoberturaPluginStateStore
	Stats CoverageStats
}

type CoberturaPluginStateStore struct {
	WorkSpacePath           string
	CompleteCoverageXmlPath string
}

func (c *CoberturaPlugin) Init(args *pd.Args) error {
	c.InputArgs = args
	c.CoberturaPluginStateStore.WorkSpacePath = pd.GetTestWorkSpaceDir()
	return nil
}

func (c *CoberturaPlugin) GetWorkSpaceDir() string {
	return c.CoberturaPluginStateStore.WorkSpacePath
}

func (c *CoberturaPlugin) GetCoberturaFilesPathPattern() string {
	return c.InputArgs.ExecFilesPathPattern
}

func (c *CoberturaPlugin) SetBuildRoot(buildRootPath string) error {
	return nil
}

func (c *CoberturaPlugin) DeInit() error {
	return nil
}

func (c *CoberturaPlugin) ValidateAndProcessArgs(args pd.Args) error {
	return nil
}

func (c *CoberturaPlugin) DoPostArgsValidationSetup(args pd.Args) error {
	return nil
}

func (c *CoberturaPlugin) Run() error {
	err := c.LocateCoberturaCoverageXmlPath()
	if err != nil {
		return err
	}

	c.Stats, err = GetCoberturaCoverageMetrics(c.CompleteCoverageXmlPath)
	if err != nil {
		return err
	}

	if c.InputArgs.PluginFailOnThreshold == true {
		isGood := c.AnalyzeCoberturaThresholds()
		if !isGood {
			return pd.GetNewError("Cobertura thresholds not met")
		}
	}

	c.Stats.PrintToConsole()
	return nil
}

func (c *CoberturaPlugin) AnalyzeCoberturaThresholds() bool {

	if c.InputArgs.PluginFailOnThreshold == false {
		return true
	}

	type ThresholdsCompare struct {
		ObservedValue float64
		ExpectedValue float64
		ThresholdType string
	}

	complexityDensity := float64(c.Stats.Complexity) / float64(c.Stats.LOC)

	thresholdsCompareList := []ThresholdsCompare{
		{c.Stats.BranchCoverage, c.InputArgs.MinimumBranchCoverage, "Branch"},
		{c.Stats.ClassCoverage, c.InputArgs.MinimumClassCoverage, "Class"},
		{c.Stats.LineCoverage, c.InputArgs.MinimumLineCoverage, "Line"},
		{c.Stats.MethodCoverage, c.InputArgs.MinimumMethodCoverage, "Method"},
		{c.Stats.PackageCoverage, c.InputArgs.MinimumPackageCoverage, "Package"},
		{c.Stats.FileCoverage, c.InputArgs.MinimumFileCoverage, "File"},
	}

	for _, thresholdCompare := range thresholdsCompareList {
		if thresholdCompare.ObservedValue < thresholdCompare.ExpectedValue {
			pd.LogPrintln(c, "CoberturaPlugin "+thresholdCompare.ThresholdType+" threshold not met",
				" expected = ", thresholdCompare.ExpectedValue, " observed = ", thresholdCompare.ObservedValue)
			logrus.Printf("Threshold type: %s threshold not met expected = %.2f observed = %.2f\n",
				thresholdCompare.ThresholdType, thresholdCompare.ExpectedValue, thresholdCompare.ObservedValue)
			return false
		}
	}

	if (c.Stats.LOC < c.InputArgs.MinimumLOC) ||
		(c.Stats.Complexity > c.InputArgs.MinimumComplexityCoverage) ||
		complexityDensity > c.InputArgs.MaxComplexityDensityCoverage {

		pd.LogPrintln(c, "CoberturaPlugin Complexity threshold not met",
			" expected = ", c.InputArgs.MinimumComplexityCoverage, " observed = ", c.Stats.Complexity)
		fmt.Println("CoberturaPlugin Complexity threshold not met",
			" expected = ", c.InputArgs.MinimumComplexityCoverage, " observed = ", c.Stats.Complexity)

		return false
	}

	return true
}

func (c *CoberturaPlugin) LocateCoberturaCoverageXmlPath() error {

	workSpaceDir := c.GetWorkSpaceDir()
	if workSpaceDir == "" {
		return pd.GetNewError("Workspace dir not set")
	}

	completeWorkSpaceDir, err := filepath.Abs(workSpaceDir)
	if err != nil {
		return err
	}

	baseSearchDir := os.DirFS(completeWorkSpaceDir)
	matchedDirs, err := doublestar.Glob(baseSearchDir, c.GetCoberturaFilesPathPattern())
	if err != nil {
		return err
	}

	if len(matchedDirs) < 1 {
		return pd.GetNewError("No Cobertura report xml found")
	}

	relativeXmlReportPath := matchedDirs[0]

	completeXmlPath := filepath.Join(completeWorkSpaceDir, relativeXmlReportPath)
	_, err = os.Stat(completeXmlPath)
	if err != nil {
		return err
	}

	c.CompleteCoverageXmlPath = completeXmlPath

	return nil
}

func (c *CoberturaPlugin) WriteOutputVariables() error {

	type EnvKvPair struct {
		Key   string
		Value interface{}
	}

	var kvPairs = []EnvKvPair{
		{Key: "BRANCH_COVERAGE", Value: fmt.Sprintf("%.2f", c.Stats.BranchCoverage)},
		{Key: "LINE_COVERAGE", Value: fmt.Sprintf("%.2f", c.Stats.LineCoverage)},
		{Key: "METHOD_COVERAGE", Value: fmt.Sprintf("%.2f", c.Stats.MethodCoverage)},
		{Key: "CLASS_COVERAGE", Value: fmt.Sprintf("%.2f", c.Stats.ClassCoverage)},
		{Key: "FILE_COVERAGE", Value: fmt.Sprintf("%.2f", c.Stats.FileCoverage)},
		{Key: "PACKAGE_COVERAGE", Value: fmt.Sprintf("%.2f", c.Stats.PackageCoverage)},
		{Key: "COMPLEXITY_COVERAGE", Value: c.Stats.Complexity},
		{Key: "COMPLEXITY_DENSITY", Value: c.Stats.ComplexityDensity},
		{Key: "LOC", Value: c.Stats.LOC},
	}

	var retErr error = nil

	for _, kvPair := range kvPairs {
		err := pd.WriteEnvVariableAsString(kvPair.Key, kvPair.Value)
		if err != nil {
			retErr = err
		}
	}

	return retErr
}

func (c *CoberturaPlugin) PersistResults() error {
	return nil
}

func (c *CoberturaPlugin) GetPluginType() string {
	return "cobertura"
}

func (c *CoberturaPlugin) IsQuiet() bool {
	return false
}

func (c *CoberturaPlugin) InspectProcessArgs(argNamesList []string) (map[string]interface{}, error) {
	return nil, nil
}

func GetNewCoberturaPlugin() CoberturaPlugin {
	return CoberturaPlugin{}
}
