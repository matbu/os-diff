/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */
package godiff

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

type CompareFileNames struct {
	Origin      string
	Destination string
	reportLines []string
}

func writeReport(content []string, reportPath string) error {

	path, _ := filepath.Split(reportPath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
	reportContent := strings.Join(content, "")
	err := ioutil.WriteFile(reportPath, []byte(reportContent), 0644)
	if err != nil {
		return errors.New("Failed to write report file: '" + reportPath + "'. " + err.Error())
	}
	return nil
}

// @todo : detect config type: json / ini / yaml
func (f *CompareFileNames) makeDiff(file1 []string, file2 []string) []string {
	// Console colors
	colorRed := "\033[31m"
	colorReset := "\033[0m"
	// Check for differences
	reportLines := []string{}
	reportLines = append(reportLines, fmt.Sprintf("Compare configuration file: %s with configuration file: %s\n", f.Origin, f.Destination))

	for i, line1 := range file1 {
		found := false
		// Skip comments
		if !strings.HasPrefix(line1, "#") && len(line1) > 0 && !strings.HasPrefix(line1, "[") {
			config1 := strings.Split(line1, "=")
			for j, line2 := range file2 {
				if !strings.HasPrefix(line2, "#") && len(line2) > 0 && !strings.HasPrefix(line2, "[") {
					config2 := strings.Split(line2, "=")
					// key found in both, check if value is different
					if config1[0] == config2[0] {
						if config1[1] != config2[1] {
							fmt.Println("*** Difference detected:", string(colorRed), line1, line2, string(colorReset))
							reportLines = append(reportLines, fmt.Sprintf("Line %d: %s, different with %s\n", j+1, line1, line2))
						}
						found = true
					}
				}
			}
			if !found {
				fmt.Println("--- Line not found:", string(colorRed), line1, string(colorReset))
				reportLines = append(reportLines, fmt.Sprintf("Line %d: %s, Not found \n", i+1, line1))
			}
		}
	}
	return reportLines
}

func (f *CompareFileNames) CompareIniFiles(origin string, dest string) []string {
	// Load the first INI file
	cfg1, err := ini.Load(origin)
	if err != nil {
		fmt.Errorf("Erro while loading file %s: %s", origin, err)
	}
	cfg2, err := ini.Load(dest)
	if err != nil {
		fmt.Errorf("Erro while loading file %s: %s", dest, err)
	}

	// Compare the sections and keys in each file
	// Console colors
	colorRed := "\033[31m"
	colorReset := "\033[0m"
	// Check for differences
	for _, sec1 := range cfg1.Sections() {
		sec2, err := cfg2.GetSection(sec1.Name())
		if err != nil {
			msg := fmt.Sprintf("Section %s not found in %s \n", sec1.Name(), dest)
			if !stringInSlice(msg, f.reportLines) {
				fmt.Println(string(colorRed), "*** Difference detected -- Section: ", sec1.Name(), " not found in:", dest, string(colorReset))
				f.reportLines = append(f.reportLines, msg)
			}
		}

		for _, key1 := range sec1.Keys() {
			key2, err := sec2.GetKey(key1.Name())
			if err != nil {
				msg := fmt.Sprintf("Key %s not found in section %s of %s \n", key1.Name(), sec1.Name(), dest)
				if !stringInSlice(msg, f.reportLines) {
					fmt.Println(string(colorRed), "*** Difference detected -- Section: ", sec1.Name(), " Key ", key1.Name(), " not found in:", dest, string(colorReset))
					f.reportLines = append(f.reportLines, msg)
				}
			} else {
				if key1.Value() != key2.Value() {
					msg := fmt.Sprintf("Value of key %s in section %s is different between %s and %s \n", key1.Name(), sec1.Name(), origin, dest)
					if !stringInSlice(msg, f.reportLines) {
						fmt.Println(string(colorRed), "*** Difference detected: Values are not equal: ", key1.Value(), " and ", key2.Value(), "Section: ", sec1.Name(), " Key ", key1.Name(), dest, string(colorReset))
						f.reportLines = append(f.reportLines, msg)
					}
				}
			}
		}
	}

	return f.reportLines
}

func (f *CompareFileNames) CompareFiles() ([]string, error) {
	// Compare two files
	// Read the files
	orgContent, err := ioutil.ReadFile(f.Origin)
	if err != nil {
		return nil, errors.New("Failed to open file: '" + f.Origin + "'. " + err.Error())
	}
	destContent, err := ioutil.ReadFile(f.Destination)
	if err != nil {
		return nil, errors.New("Failed to open file: '" + f.Destination + "'. " + err.Error())
	}
	// Detect type
	reports := []string{}
	if isIni(f.Origin) {
		reports = f.CompareIniFiles(f.Origin, f.Destination)
		reports = append(reports, f.CompareIniFiles(f.Origin, f.Destination)...)
	} else {
		// Split both files into lines
		orgLines := strings.Split(string(orgContent), "\n")
		destLines := strings.Split(string(destContent), "\n")
		// Check for differences
		reports = f.makeDiff(orgLines, destLines)
		//dest_org_comparison := f.makeDiff(destLines, orgLines)
		reports = append(reports, f.makeDiff(destLines, orgLines)...)

	}
	filePath := "/tmp/" + f.Origin + ".diff"
	if len(reports) != 0 {
		err = writeReport(reports, filePath)
		if err != nil {
			fmt.Println(err)
		}
	}
	return reports, nil
}
