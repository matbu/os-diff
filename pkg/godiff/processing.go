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
	"fmt"
	"os"
	"path/filepath"
)

// @todo : checksum before doing anything
type GoDiffDataStruct struct {
	Origin          string
	Destination     string
	missingInOrg    []string
	missingPath     []string
	missingInDest   []string
	wrongTypeInOrg  []string
	wrongTypeInDest []string
	unmatchFile     []string
}

func checkFile(path1 string, path2 string) (bool, error) {
	/*
	   return:
	   true if files are equal, error nil
	   false if there is a difference, error nil
	   false and error not nil if an error occur
	*/
	statFile1, err := os.Stat(path1)
	if err != nil {
		return false, err
	}
	statFile2, err := os.Stat(path2)
	if err != nil {
		return false, err
	}
	if statFile1.IsDir() {
		return false, fmt.Errorf("Path: %s is a directoy", path1)
	}
	if statFile2.IsDir() {
		return false, fmt.Errorf("Path: %s is a directoy", path2)
	}

	file1, err := os.Open(path1)
	if err != nil {
		return false, err
	}
	defer file1.Close()
	file2, err := os.Open(path2)
	if err != nil {
		return false, err
	}
	defer file2.Close()

	buf1 := make([]byte, 1024)
	buf2 := make([]byte, 1024)

	for {
		n1, err1 := file1.Read(buf1)
		n2, err2 := file2.Read(buf2)

		if err1 != nil || err2 != nil || n1 != n2 {
			fmt.Printf("Files %s and %s are different\n", path1, path2)
			return false, nil
			break
		}
		if n1 == 0 {
			break
		}
		if string(buf1[:n1]) != string(buf2[:n2]) {
			fmt.Printf("Files %s and %s are different\n", path1, path2)
			return false, nil
			break
		}
	}
	return true, nil
}

func (p *GoDiffDataStruct) Process(dir1 string, dir2 string) error {
	// Walk through the first directory and compare each file to the second directory
	// Start to process
	fmt.Println("Start processing: ", dir1, "and: ", dir2, "\n")
	// Walk through DIR 1
	filepath.Walk(dir1, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Get the corresponding file in the second directory
		relPath, _ := filepath.Rel(dir1, path)
		path2 := filepath.Join(dir2, relPath)
		file1, err := os.Stat(path)
		file2, err := os.Stat(path2)
		if err != nil {
			if !stringInSlice(path, p.missingPath) {
				p.missingPath = append(p.missingPath, path)
			}
		} else {
			if file1.IsDir() && !file2.IsDir() {
				if !stringInSlice(path, p.wrongTypeInOrg) {
					p.wrongTypeInOrg = append(p.wrongTypeInOrg, path)
				}
			}
			if !file1.IsDir() && file2.IsDir() {
				if !stringInSlice(path, p.wrongTypeInDest) {
					p.wrongTypeInDest = append(p.wrongTypeInDest, path2)
				}
			}
			if !file1.IsDir() && !file2.IsDir() {
				check, err := checkFile(path, path2)
				if err != nil {
					return err
				}
				if !check {
					// Compare the two files
					fmt.Println("Start compare in file: ", path, "and: ", path2, "\n")
					if !stringInSlice(path, p.unmatchFile) {
						p.unmatchFile = append(p.unmatchFile, path)

						compareFiles := CompareFileNames{
							Origin:      path,
							Destination: path2,
						}
						report, err := compareFiles.CompareFiles()
						if err != nil {
							return err
						}

						if report != nil {
							// File didnt match
							if !stringInSlice(path, p.unmatchFile) {
								p.unmatchFile = append(p.unmatchFile, path)
							}
							if !stringInSlice(path2, p.unmatchFile) {
								p.unmatchFile = append(p.unmatchFile, path2)
							}
						}
					}

				}

			}
		}
		return nil
	})
	return nil
}

func (p *GoDiffDataStruct) ProcessDirectories() error {
	// Compare origin vs destination
	p.Process(p.Origin, p.Destination)
	// Compare destination vs origin
	p.Process(p.Destination, p.Origin)
	fmt.Printf("Missing files: %s \n", p.missingPath)
	fmt.Printf("Files with differences: %s \n", p.unmatchFile)
	fmt.Printf("Different file type in origin: %s \n", p.wrongTypeInOrg)

	return nil
}
