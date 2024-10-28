package cobertura

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Coverage struct {
	XMLName  xml.Name  `xml:"coverage"`
	Packages []Package `xml:"packages>package"`
}

type Package struct {
	Name       string  `xml:"name,attr"`
	Classes    []Class `xml:"classes>class"`
	BranchRate float64 `xml:"branch-rate,attr"`
	LineRate   float64 `xml:"line-rate,attr"`
}

type Class struct {
	Name       string   `xml:"name,attr"`
	Complexity float64  `xml:"complexity,attr"`
	BranchRate float64  `xml:"branch-rate,attr"`
	LineRate   float64  `xml:"line-rate,attr"`
	Lines      []Line   `xml:"lines>line"`
	Methods    []Method `xml:"methods>method"`
}

type Method struct {
	Name       string  `xml:"name,attr"`
	BranchRate float64 `xml:"branch-rate,attr"`
	LineRate   float64 `xml:"line-rate,attr"`
	Lines      []Line  `xml:"lines>line"`
}

type Line struct {
	Number            int         `xml:"number,attr"`
	Branch            bool        `xml:"branch,attr"`
	Hits              int         `xml:"hits,attr"`
	Conditions        []Condition `xml:"conditions>condition"`
	ConditionCoverage string      `xml:"condition-coverage,attr"` // Condition coverage field
}

type Condition struct {
	Number   int    `xml:"number,attr"`
	Coverage string `xml:"coverage,attr"` // Use the coverage attribute
}

type CoverageStats struct {
	ClassCoverage     float64
	BranchCoverage    float64
	LineCoverage      float64
	MethodCoverage    float64
	PackageCoverage   float64
	FileCoverage      float64
	Complexity        int
	ComplexityDensity string
	LOC               int
}

func GetCoberturaCoverageMetrics(coverageXmlCompletePath string) (CoverageStats, error) {

	file, err := os.Open(coverageXmlCompletePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return CoverageStats{}, err
	}
	defer file.Close()

	var coverage Coverage
	if err := xml.NewDecoder(file).Decode(&coverage); err != nil {
		fmt.Println("Error decoding XML:", err)
		return CoverageStats{}, err
	}

	stats := calculateCoverage(coverage)
	return stats, nil
}

func calculateCoverage(c Coverage) CoverageStats {

	var totalLines, totalCovered int
	var totalBranches, totalCoveredBranches int
	var totalPackages, totalCoveredPackages int
	var totalComplexity float64 = 0.0
	var totalMethods, totalMethodsCovered int
	var totalClasses, totalCoveredClasses int

	for _, pkg := range c.Packages {
		var pkgLines, pkgCovered int
		packageCoveredClasses := 0

		for _, class := range pkg.Classes {
			classLines, classLinesCovered := getLineStats(class.Lines)
			totalClasses++

			if classLinesCovered > 0 {
				totalCoveredClasses++
				packageCoveredClasses++
			}
			pkgLines += classLines
			pkgCovered += classLinesCovered
			totalComplexity += class.Complexity

			for _, line := range class.Lines {
				if line.ConditionCoverage != "" {
					coveredConditions, totalConditions := parseConditionCoverage(line.ConditionCoverage)
					totalCoveredBranches += coveredConditions
					totalBranches += totalConditions
				}
			}

			classMethods, classMethodsCovered := getMethodStats(class.Methods)
			totalMethods += classMethods
			totalMethodsCovered += classMethodsCovered

		}

		totalLines += pkgLines
		totalCovered += pkgCovered
		totalPackages++

		if packageCoveredClasses > 0 {
			totalCoveredPackages++
		}
	}

	packageCoverage := calculatePercentage(totalPackages, totalCoveredPackages)
	fileCoverage := calculatePercentage(totalCovered, totalPackages)
	classCoverage := calculatePercentage(totalCoveredClasses, totalClasses)
	branchCoverage := calculatePercentage(totalCoveredBranches, totalBranches)
	lineCoverage := calculatePercentage(totalCovered, totalLines)

	fmt.Printf("Branch covered: %d Total branches: %d\n", totalCoveredBranches, totalBranches)
	fmt.Printf("Methods covered: %d Total methods: %d\n", totalMethodsCovered, totalMethods)
	fmt.Printf("Classes covered: %d Total classes: %d\n", totalCoveredClasses, totalClasses)
	fmt.Printf("Complexity: %.2f\n", totalComplexity)
	fmt.Printf("Total Lines: %d\n", totalLines)
	fmt.Printf("Method Coverage: %d\n", totalMethods)

	return CoverageStats{
		PackageCoverage:   packageCoverage,
		FileCoverage:      fileCoverage,
		ClassCoverage:     classCoverage,
		BranchCoverage:    branchCoverage,
		LineCoverage:      lineCoverage,
		Complexity:        int(totalComplexity),
		ComplexityDensity: fmt.Sprintf("%d/%d", int(totalComplexity), totalLines),
		LOC:               totalLines,
	}
}

func getLineStats(lines []Line) (int, int) {
	var totalLines, totalCovered int
	for _, line := range lines {
		totalLines++
		if line.Hits > 0 {
			totalCovered++
		}
	}
	return totalLines, totalCovered
}

func getMethodStats(methods []Method) (int, int) {
	var totalMethods, totalCovered int
	for _, method := range methods {
		totalMethods++
		_, linesCovered := getLineStats(method.Lines)
		if linesCovered > 0 {
			totalCovered++
		}
	}
	return totalMethods, totalCovered
}

func parseConditionCoverage(coverage string) (int, int) {
	re := regexp.MustCompile(`(\d+)% \((\d+)/(\d+)\)`)
	matches := re.FindStringSubmatch(coverage)
	if len(matches) != 4 {
		return 0, 0
	}

	covered, _ := strconv.Atoi(matches[2])
	total, _ := strconv.Atoi(matches[3])
	return covered, total
}

func calculatePercentage(part, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(part) / float64(total) * 100
}

//
//func printCoverageStats(stats CoverageStats) {
//	fmt.Printf("Package Coverage: %.2f%%\n", stats.PackageCoverage)
//	fmt.Printf("File Coverage: %.2f%%\n", stats.FileCoverage)
//	fmt.Printf("Class Coverage: %.2f%%\n", stats.ClassCoverage)
//	fmt.Printf("Branch Coverage: %.2f%%\n", stats.BranchCoverage)
//	fmt.Printf("Line Coverage: %.2f%%\n", stats.LineCoverage)
//	fmt.Printf("Complexity: %v\n", stats.Complexity)
//	fmt.Printf("Complexity Density: %v\n", stats.ComplexityDensity)
//	fmt.Printf("LOC: %v\n", stats.LOC)
//
//}

func (stats *CoverageStats) PrintToConsole() {
	fmt.Printf("Package Coverage: %.2f%%\n", stats.PackageCoverage)
	fmt.Printf("File Coverage: %.2f%%\n", stats.FileCoverage)
	fmt.Printf("Class Coverage: %.2f%%\n", stats.ClassCoverage)
	fmt.Printf("Branch Coverage: %.2f%%\n", stats.BranchCoverage)
	fmt.Printf("Line Coverage: %.2f%%\n", stats.LineCoverage)
	fmt.Printf("Complexity: %v\n", stats.Complexity)
	fmt.Printf("Complexity Density: %v\n", stats.ComplexityDensity)
	fmt.Printf("LOC: %v\n", stats.LOC)
}
