/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reporters

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/ginkgo/types"
)

type JUnitTestSuites struct {
	XMLName    xml.Name                   `xml:"testsuites"`
	TestSuites []reporters.JUnitTestSuite `xml:"testsuite"`
	Name       string                     `xml:"name,attr"`
	Tests      int                        `xml:"tests,attr"`
	Failures   int                        `xml:"failures,attr"`
	Errors     int                        `xml:"errors,attr"`
	Time       float64                    `xml:"time,attr"`
}

type JUnitReporterWithSuites struct {
	filename string
	*reporters.JUnitReporter
}

func NewJUnitReportWithSuites(filename string) *JUnitReporterWithSuites {
	return &JUnitReporterWithSuites{
		filename:      filename,
		JUnitReporter: reporters.NewJUnitReporter(filename),
	}
}
func (reporter *JUnitReporterWithSuites) SpecSuiteDidEnd(summary *types.SuiteSummary) {
	reporter.JUnitReporter.SpecSuiteDidEnd(summary)
	bytes, err := os.ReadFile(reporter.filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %q: %v\n", reporter.filename, err)
		return
	}
	testsuite := reporters.JUnitTestSuite{}
	err = xml.Unmarshal(bytes, &testsuite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read unmarshal JUnitTestSuite object: %v\n", err)
		return
	}
	testsuites := JUnitTestSuites{
		TestSuites: []reporters.JUnitTestSuite{testsuite},
		Name:       testsuite.Name,
		Tests:      testsuite.Tests,
		Failures:   testsuite.Failures,
		Errors:     testsuite.Errors,
		Time:       testsuite.Time,
	}
	file, err := os.Create(reporter.filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %q for writing: %v\n", reporter.filename, err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(xml.Header)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write XML header to %q: %v\n", reporter.filename, err)
		return
	}
	encoder := xml.NewEncoder(file)
	err = encoder.Encode(testsuites)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write JUnitTestSuites object to %q: %v\n", reporter.filename, err)
		return
	}
}
