package plugin_defs

type Plugin interface {
	Init(args Args) error
	SetBuildRoot(buildRootPath string) error
	DeInit() error
	ValidateAndProcessArgs(args Args) error
	DoPostArgsValidationSetup(args Args) error
	Run() error
	WriteOutputVariables() error
	PersistResults() error
	GetPluginType() string
	IsQuiet() bool
	InspectProcessArgs(argNamesList []string) (map[string]interface{}, error)
}

type Args struct {
	Pipeline
	CoveragePluginArgs
	EnvPluginInputArgs
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`
}

type CoveragePluginArgs struct {
	PluginToolType        string `envconfig:"PLUGIN_TOOL"`
	PluginFailOnThreshold bool   `envconfig:"PLUGIN_FAIL_ON_THRESHOLD"`
	PluginFailIfNoReports bool   `envconfig:"PLUGIN_FAIL_IF_NO_REPORTS"`
}

type EnvPluginInputArgs struct {
	ExecFilesPathPattern string `envconfig:"PLUGIN_REPORTS_PATH_PATTERN"`

	ClassPatterns          string `envconfig:"PLUGIN_CLASS_DIRECTORIES"`
	ClassInclusionPatterns string `envconfig:"PLUGIN_CLASS_INCLUSION_PATTERN"`
	ClassExclusionPatterns string `envconfig:"PLUGIN_CLASS_EXCLUSION_PATTERN"`

	SourcePattern          string `envconfig:"PLUGIN_SOURCE_DIRECTORIES"`
	SourceInclusionPattern string `envconfig:"PLUGIN_SOURCE_INCLUSION_PATTERN"`
	SourceExclusionPattern string `envconfig:"PLUGIN_SOURCE_EXCLUSION_PATTERN"`

	SkipCopyOfSrcFiles bool `envconfig:"PLUGIN_SKIP_SOURCE_COPY"`

	MinimumInstructionCoverage float64 `envconfig:"PLUGIN_THRESHOLD_INSTRUCTION"`
	MinimumBranchCoverage      float64 `envconfig:"PLUGIN_THRESHOLD_BRANCH"`
	MinimumComplexityCoverage  int     `envconfig:"PLUGIN_THRESHOLD_COMPLEXITY"`
	MinimumLineCoverage        float64 `envconfig:"PLUGIN_THRESHOLD_LINE"`
	MinimumMethodCoverage      float64 `envconfig:"PLUGIN_THRESHOLD_METHOD"`
	MinimumClassCoverage       float64 `envconfig:"PLUGIN_THRESHOLD_CLASS"`
}

const (
	JacocoPluginType    = "jacoco"
	JacocoXmlPluginType = "jacoco-xml"
	CorbeturaPluginType = "corbetura"
)
