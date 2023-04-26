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
package diff

import (
  "errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type CompareFileNames struct {
	Origin       string
  Destination  string
}

// @todo : detect config type: json / ini / yaml
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func writeReport(content []string, report_path string) error {
  var err error
  // Write the report to file
  reportContent := strings.Join(content, "")
  err = ioutil.WriteFile(report_path, []byte(reportContent), 0644)
  if err != nil {
    return errors.New("Failed to write report file: '" + report_path + "'. " + err.Error())
  }
  return nil
}

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

func (f *CompareFileNames) CompareIniFiles(report_path string) error {

  // Read the files
  orgContent, err := ioutil.ReadFile(f.Origin)
  if err != nil {
    return errors.New("Failed to open file: '" + f.Origin + "'. " + err.Error())
  }
  destContent, err := ioutil.ReadFile(f.Destination)
  if err != nil {
    return errors.New("Failed to open file: '" + f.Destination + "'. " + err.Error())
  }
  // Split both files into lines
  orgLines := strings.Split(string(orgContent), "\n")
  destLines := strings.Split(string(destContent), "\n")

  // Check for differences
  reports := f.makeDiff(orgLines, destLines)
  //dest_org_comparison := f.makeDiff(destLines, orgLines)
  reports = append(reports, f.makeDiff(destLines, orgLines)...)
  writeReport(reports, report_path)

	return nil
}
