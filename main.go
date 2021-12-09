/*
   Copyright 2021 Dmitriy Poluyanov

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
package main

import (
	bep "bazel-wsl/bep"
	"bazel-wsl/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// open output file
	var fo, err = os.OpenFile("C:\\Users\\la_d.poluyanov\\GolandProjects\\bazel\\command_log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	//var fo, err = os.OpenFile("/Users/d.poluyanov/workspace/bazel-wsl/command_log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	var argsLine = ""

	for _, arg := range os.Args {
		argsLine += arg + " "
	}
	_, err = fo.WriteString(argsLine + "\n")
	if err != nil {
		panic(err)
	}

	bazelArgs := make([]string, 1)
	bazelArgs[0] = "bazel"

	var bepIDEAOutputPath = ""
	var bepOutputPath = ""
	for _, arg := range os.Args[1:] {
		// https://github.com/bazelbuild/intellij/pull/2976/files
		fo.WriteString(arg + "\n")
		if strings.HasPrefix(arg, "attr(tags") {
			var newAttr = strings.Replace(arg, "attr(tags", "attr('tags'", 1)
			newAttr = strings.Replace(newAttr, "^((?!manual).)*$", "'^((?!manual).)*$'", 1)
			bazelArgs = append(bazelArgs, newAttr)
		} else if strings.HasPrefix(arg, "same_pkg_direct_rdeps(") {
			var newAttr = strings.Replace(arg, "same_pkg_direct_rdeps(", "'same_pkg_direct_rdeps(", 1)
			newAttr = newAttr + "'"
			bazelArgs = append(bazelArgs, newAttr)
		} else if strings.HasPrefix(arg, "--build_event_binary_file=") {
			var fileName = strings.Replace(arg, "--build_event_binary_file=", "", 1)

			var wslPath = utils.WinToWSLPath(fileName)

			bepOutputPath = fileName + "_bazel_wsl"
			bepIDEAOutputPath = fileName
			bazelArgs = append(bazelArgs, "--build_event_binary_file="+wslPath+"_bazel_wsl")
		} else if strings.HasPrefix(arg, "--override_repository=intellij_aspect=") {
			var fileName = strings.Replace(arg, "--override_repository=intellij_aspect=", "", 1)

			wslPath := utils.WinToWSLPath(fileName)

			bazelArgs = append(bazelArgs, "--override_repository=intellij_aspect="+wslPath)
		} else {
			bazelArgs = append(bazelArgs, arg)
		}
	}

	var translatedArgsLine = ""

	for _, arg := range bazelArgs {
		translatedArgsLine += arg + " "
	}

	fo.WriteString(fmt.Sprintf("Translated to wsl %s\n", translatedArgsLine))
	for _, arg := range bazelArgs {
		fo.WriteString(arg + "\n")
	}

	//logFile, err := os.OpenFile("C:\\Users\\la_d.poluyanov\\GolandProjects\\bazel\\log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer func() {
	//	logFile.Close()
	//}()
	//mw := io.MultiWriter(os.Stdout, logFile)
	//logFileErr, err := os.OpenFile("C:\\Users\\la_d.poluyanov\\GolandProjects\\bazel\\log_err.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	logFileErr.Close()
	//}()
	//mwErr := io.MultiWriter(os.Stderr, logFileErr)

	cmd := exec.Command("wsl", bazelArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fo.WriteString(fmt.Sprint(err))
		os.Exit(err.(*exec.ExitError).ExitCode())
	}

	if bepOutputPath != "" && bepIDEAOutputPath != "" {
		bepFrom, err1 := os.OpenFile(bepOutputPath, os.O_RDONLY, 0600)
		if err1 != nil {
			panic(err1)
		}
		defer func() {
			os.Remove(bepOutputPath)
			if err := bepFrom.Close(); err != nil {
				panic(err)
			}
		}()

		var bepTo, err = os.OpenFile(bepIDEAOutputPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := bepTo.Close(); err != nil {
				panic(err)
			}
		}()

		bep.RewriteBep(bepFrom, bepTo, fo)
		//bepTo.Close()

		//time.Sleep(60 * time.Second)
	}
}
