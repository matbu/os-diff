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
package cmd

import (
	"fmt"

	"os-diff/pkg/diff"
	"github.com/spf13/cobra"
)

// compareCmd represents the compare command
var origin string
var destination string
var output string

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two files or directories",
	Long: `Compare files or directories from two different paths.
				 For example:
				 os-diff compare --origin=tests/podman/keystone.conf --destination=tests/ocp/keystone.conf --output=output.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		compareFiles := &diff.CompareFileNames{
			Origin: origin,
			Destination: destination,
		}
		err := compareFiles.CompareIniFiles(output)
		if err != nil {
			panic(err)
		}

		fmt.Println("Compare called")
	},
}

func init() {
	compareCmd.Flags().StringVar(&origin, "origin", "", "Origin file or directory.")
	compareCmd.Flags().StringVar(&destination, "destination", "", "Destination file or directory")
	compareCmd.Flags().StringVar(&output, "output", "output.txt", "Output file (default is $PWD/output.txt)")
	rootCmd.AddCommand(compareCmd)
}
