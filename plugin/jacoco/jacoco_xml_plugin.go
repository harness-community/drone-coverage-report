package jacoco

import (
	"github.com/bmatcuk/doublestar/v4"
	pd "github.com/harness-community/drone-coverage-report/plugin/plugin_defs"
	"os"
	"path/filepath"
)

type JacocoXmlPlugin struct {
	JacocoBasePlugin      JacocoPlugin
	XmlReportCompletePath string
}

func GetNewJacocoXmlPlugin() JacocoXmlPlugin {
	return JacocoXmlPlugin{}
}

func (jxp *JacocoXmlPlugin) Init(args *pd.Args) error {
	err := jxp.JacocoBasePlugin.Init(args)
	if err != nil {
		pd.LogPrintln(jxp, "Error in JacocoXmlPlugin Init: ", err.Error())
		return err
	}
	return nil
}

func (jxp *JacocoXmlPlugin) SetBuildRoot(buildRootPath string) error {
	pd.LogPrintln(jxp, "Setting build root path in JacocoXmlPlugin")
	return nil
}

func (jxp *JacocoXmlPlugin) DeInit() error {
	pd.LogPrintln(jxp, "DeInit in JacocoXmlPlugin")
	return nil
}

func (jxp *JacocoXmlPlugin) ValidateAndProcessArgs(args pd.Args) error {
	pd.LogPrintln(jxp, "Validating and processing args in JacocoXmlPlugin")
	if args.ExecFilesPathPattern == "" {
		return pd.GetNewError("JacocoXmlPlugin: No reports path pattern provided")
	}
	return nil
}

func (jxp *JacocoXmlPlugin) DoPostArgsValidationSetup(args pd.Args) error {
	pd.LogPrintln(jxp, "DoPostArgsValidationSetup in JacocoXmlPlugin")

	err := jxp.LocateJacocoReportXml()
	if err != nil {
		return err
	}

	return nil
}

func (jxp *JacocoXmlPlugin) LocateJacocoReportXml() error {
	pd.LogPrintln(jxp, "Finding Jacoco report xml")

	rootSearchDirFS := os.DirFS(jxp.JacocoBasePlugin.GetWorkspaceDir())
	matchedDirs, err := doublestar.Glob(rootSearchDirFS, jxp.JacocoBasePlugin.InputArgs.ExecFilesPathPattern)
	if err != nil {
		return err
	}

	if len(matchedDirs) < 1 {
		return pd.GetNewError("No Jacoco report xml found")
	}

	relativeXmlReportPath := matchedDirs[0]
	xmlPath := filepath.Join(jxp.JacocoBasePlugin.GetWorkspaceDir(), relativeXmlReportPath)
	_, err = os.Stat(xmlPath)
	if err != nil {
		return pd.GetNewError("Error in CheckFilesCopiedToWorkSpace: " + err.Error())
	}

	completeXmlPath, err := filepath.Abs(xmlPath)
	if err != nil {
		return pd.GetNewError("Error in getting absolute path: " + err.Error())
	}

	pd.LogPrintln(jxp, "Jacoco report xml found at: ", completeXmlPath)

	jxp.XmlReportCompletePath = completeXmlPath

	return nil
}

func (jxp *JacocoXmlPlugin) Run() error {
	pd.LogPrintln(jxp, "Running JacocoXmlPlugin")

	jacocoThresholdValues := GetJacocoCoverageThresholds(jxp.XmlReportCompletePath)
	pd.LogPrintln(jxp, "Jacoco threshold values: ", jacocoThresholdValues)

	jxp.JacocoBasePlugin.SetCoverageThresholds(jacocoThresholdValues)
	isGood := jxp.JacocoBasePlugin.IsThresholdValuesGood()
	if !isGood {
		pd.LogPrintln(jxp, "JacocoXmlPlugin: Coverage thresholds not met")
		return pd.GetNewError("JacocoXmlPlugin: Coverage thresholds not met")
	}

	return nil
}

func (jxp *JacocoXmlPlugin) WriteOutputVariables() error {
	pd.LogPrintln(jxp, "Writing output variables in JacocoXmlPlugin")
	jxp.JacocoBasePlugin.WriteOutputVariables()
	return nil
}

func (jxp *JacocoXmlPlugin) PersistResults() error {
	pd.LogPrintln(jxp, "Persisting results in JacocoXmlPlugin")
	return nil
}

func (jxp *JacocoXmlPlugin) GetPluginType() string {
	pd.LogPrintln(jxp, "Getting plugin type in JacocoXmlPlugin")
	return pd.JacocoXmlPluginType
}

func (jxp *JacocoXmlPlugin) IsQuiet() bool {
	pd.LogPrintln(jxp, "Checking if JacocoXmlPlugin is quiet")
	return false
}

func (jxp *JacocoXmlPlugin) InspectProcessArgs(argNamesList []string) (map[string]interface{}, error) {
	pd.LogPrintln(jxp, "Inspecting process args in JacocoXmlPlugin")
	return nil, nil
}
