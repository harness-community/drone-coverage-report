package jacoco

import (
	"encoding/xml"
	"fmt"
	plg "github.com/harness-community/drone-coverage-report/plugin/plugin_defs"
	"io"
	"log"
	"os"
)

type Report struct {
	XMLName  xml.Name  `xml:"report"`
	Counters []Counter `xml:"counter"`
	Packages []Package `xml:"package"`
}

type Counter struct {
	Type    string `xml:"type,attr"`
	Missed  int    `xml:"missed,attr"`
	Covered int    `xml:"covered,attr"`
}

type Package struct {
	Name     string    `xml:"name,attr"`
	Counters []Counter `xml:"counter"`
}

func (j *JacocoCoverageThresholds) ToFloat64() JacocoCoverageThresholdsValues {
	return JacocoCoverageThresholdsValues{
		InstructionCoverageThreshold: ParsePercentage(j.InstructionCoverageThreshold),
		BranchCoverageThreshold:      ParsePercentage(j.BranchCoverageThreshold),
		LineCoverageThreshold:        ParsePercentage(j.LineCoverageThreshold),
		ComplexityCoverageThreshold:  j.ComplexityCoverageThreshold,
		MethodCoverageThreshold:      ParsePercentage(j.MethodCoverageThreshold),
		ClassCoverageThreshold:       ParsePercentage(j.ClassCoverageThreshold),
	}
}

func ParsePercentage(percentage string) float64 {
	var value float64
	_, err := fmt.Sscanf(percentage, "%f", &value)
	if err != nil {
		log.Fatalf("Error parsing percentage: %v", err)
	}
	return value
}

func CalculatePercentage(covered, missed int) string {
	total := covered + missed
	if total == 0 {
		return "0%(0/0)"
	}
	percentage := (float64(covered) / float64(total)) * 100
	return fmt.Sprintf("%.2f%%(%d/%d)", percentage, covered, total)
}

func GetCounterValues(counters []Counter, counterType string) (int, int) {
	for _, counter := range counters {
		if counter.Type == counterType {
			return counter.Covered, counter.Missed
		}
	}
	return 0, 0
}

func CalculateCoverageMetrics(report Report) JacocoCoverageThresholds {

	instructionCoverage, instructionMiss := GetCounterValues(report.Counters, "INSTRUCTION")
	branchCoverage, branchMiss := GetCounterValues(report.Counters, "BRANCH")
	lineCoverage, lineMiss := GetCounterValues(report.Counters, "LINE")
	complexityCoverage, complexityMiss := GetCounterValues(report.Counters, "COMPLEXITY")
	methodCoverage, methodMiss := GetCounterValues(report.Counters, "METHOD")
	classCoverage, classMiss := GetCounterValues(report.Counters, "CLASS")

	return JacocoCoverageThresholds{
		InstructionCoverageThreshold: CalculatePercentage(instructionCoverage, instructionMiss),
		BranchCoverageThreshold:      CalculatePercentage(branchCoverage, branchMiss),
		LineCoverageThreshold:        CalculatePercentage(lineCoverage, lineMiss),
		ComplexityCoverageThreshold:  complexityCoverage + complexityMiss,
		MethodCoverageThreshold:      CalculatePercentage(methodCoverage, methodMiss),
		ClassCoverageThreshold:       CalculatePercentage(classCoverage, classMiss),
	}
}

func ParseXMLReport(filename string) Report {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening XML file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading XML file: %v", err)
	}

	var report Report
	err = xml.Unmarshal(data, &report)
	if err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}
	return report
}

func GetJacocoCoverageThresholds(completeXmlPath string) JacocoCoverageThresholdsValues {
	report := ParseXMLReport(completeXmlPath)
	coverageThresholds := CalculateCoverageMetrics(report)

	plg.LogPrintf(nil, "Coverage Metrics:")
	plg.LogPrintf(nil, "Instruction Coverage: %s\n", coverageThresholds.InstructionCoverageThreshold)
	plg.LogPrintf(nil, "Branch Coverage: %s\n", coverageThresholds.BranchCoverageThreshold)
	plg.LogPrintf(nil, "Line Coverage: %s\n", coverageThresholds.LineCoverageThreshold)
	plg.LogPrintf(nil, "Complexity Coverage: %d\n", coverageThresholds.ComplexityCoverageThreshold)
	plg.LogPrintf(nil, "Method Coverage: %s\n", coverageThresholds.MethodCoverageThreshold)
	plg.LogPrintf(nil, "Class Coverage: %s\n", coverageThresholds.ClassCoverageThreshold)

	coverageThresholdValues := coverageThresholds.ToFloat64()

	return coverageThresholdValues
}
