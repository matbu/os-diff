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
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isIni(filePath string) bool {
	// @todo pass directly file content to avoid read duplication
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("Error reading file:", err)
	}
	if data[0] == '[' {
		return true
	}
	return false
}

func isYaml(filePath string) bool {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("Error reading file:", err)
	}
	var yamlData interface{}
	err = yaml.Unmarshal(data, &yamlData)
	if err == nil {
		return true
	}
	return false
}

func isJson(filePath string) bool {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("Error reading file:", err)
	}
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err == nil {
		return true
	}
	return false
}
