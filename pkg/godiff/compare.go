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
	DiffReport  []string
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

func (f *CompareFileNames) makeDiff(file1 []string, file2 []string) error {
	// Console colors
	colorRed := "\033[31m"
	colorReset := "\033[0m"
	// Check for differences
	f.DiffReport = append(f.DiffReport, fmt.Sprintf("Compare configuration file: %s with configuration file: %s\n", f.Origin, f.Destination))

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
							f.DiffReport = append(f.DiffReport, fmt.Sprintf("Line %d: %s, different with %s\n", j+1, line1, line2))
						}
						found = true
					}
				}
			}
			if !found {
				fmt.Println("--- Line not found:", string(colorRed), line1, string(colorReset))
				f.DiffReport = append(f.DiffReport, fmt.Sprintf("Line %d: %s, Not found \n", i+1, line1))
			}
		}
	}
	return nil
}

func (f *CompareFileNames) CompareIniFiles(origin string, dest string) error {
	// Load the first INI file
	cfg1, err := ini.Load(origin)
	if err != nil {
		fmt.Errorf("Erro while loading file %s: %s", origin, err)
	}
	cfg2, err := ini.Load(dest)
	if err != nil {
		fmt.Errorf("Erro while loading file %s: %s", dest, err)
	}

	diffFound := false
	// Compare the sections and keys in each file
	// Console colors
	colorRed := "\033[31m"
	colorReset := "\033[0m"
	// Check for differences
	for _, sec1 := range cfg1.Sections() {
		sec2, err := cfg2.GetSection(sec1.Name())
		if err != nil {
			msg := fmt.Sprintf("-[%s]\n", sec1.Name())
			if !stringInSlice(msg, f.DiffReport) {
				diffFound = true
				fmt.Println(string(colorRed), "*** Difference detected -- Section: ", sec1.Name(), " not found in:", dest, string(colorReset))
				f.DiffReport = append(f.DiffReport, msg)
			}
		}
		for _, key1 := range sec1.Keys() {
			key2, err := sec2.GetKey(key1.Name())
			if err != nil {
				msg := fmt.Sprintf("[%s]\n-%s=%s\n", sec1.Name(), key1.Name(), key1.Value())
				if !stringInSlice(msg, f.DiffReport) {
					diffFound = true
					fmt.Println(string(colorRed), "*** Difference detected -- Section: ", sec1.Name(), " Key ", key1.Name(), " not found in:", dest, string(colorReset))
					f.DiffReport = append(f.DiffReport, msg)
				}
			} else {
				if key1.Value() != key2.Value() {
					msg := fmt.Sprintf("[%s]\n+%s=%s\n-%s=%s\n", sec1.Name(), key1.Name(), key1.Value(), key2.Name(), key2.Value())
					if !stringInSlice(msg, f.DiffReport) {
						diffFound = true
						fmt.Println(string(colorRed), "*** Difference detected: Values are not equal: ", key1.Value(), " and ", key2.Value(), "Section: ", sec1.Name(), " Key ", key1.Name(), dest, string(colorReset))
						f.DiffReport = append(f.DiffReport, msg)
					}
				}
			}
		}
	}
	if diffFound {
		msg := fmt.Sprintf("Source file path: %s, difference with: %s\n", origin, dest)
		f.DiffReport = append([]string{msg}, f.DiffReport...)
	}
	return nil
}

func (f *CompareFileNames) CompareFiles() ([]string, error) {
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
	if isIni(f.Origin) {
		f.CompareIniFiles(f.Origin, f.Destination)
		// reverse comparison
		f.CompareIniFiles(f.Origin, f.Destination)
	} else {
		// Split both files into lines
		orgLines := strings.Split(string(orgContent), "\n")
		destLines := strings.Split(string(destContent), "\n")
		// Check for differences
		f.makeDiff(orgLines, destLines)
		f.makeDiff(destLines, orgLines)
		//reports := []string{}
		//dest_org_comparison := f.makeDiff(destLines, orgLines)
		//reports = append(reports, f.makeDiff(destLines, orgLines)...)
	}
	filePath := "/tmp/" + f.Origin + ".diff"
	if len(f.DiffReport) != 0 {
		err = writeReport(f.DiffReport, filePath)
		if err != nil {
			fmt.Println(err)
		}
	}
	return f.DiffReport, nil
}
