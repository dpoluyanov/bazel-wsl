package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// open output file
	var fo, err = os.OpenFile("C:\\Users\\la_d.poluyanov\\GolandProjects\\bazel\\command_log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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

	var sleepAtEnd = false
	for _, arg := range os.Args[1:] {
		// https://github.com/bazelbuild/intellij/pull/2976/files
		fo.WriteString(arg + "\n")
		if strings.HasPrefix(arg, "attr(tags") {
			var newAttr = strings.Replace(arg, "attr(tags", "attr('tags'", 1)
			newAttr = strings.Replace(newAttr, "^((?!manual).)*$", "'^((?!manual).)*$'", 1)
			bazelArgs = append(bazelArgs, newAttr)
		} else if strings.HasPrefix(arg, "--build_event_binary_file=") {
			var fileName = strings.Replace(arg, "--build_event_binary_file=", "", 1)

			var wslPathCmd = exec.Command("wsl", "wslpath", "-u", "-a", "'"+fileName+"'")
			r, w, err := os.Pipe()
			if err != nil {
				panic(err)
			}
			wslPathCmd.Stdout = w
			wslPathCmd.Stderr = os.Stderr
			var buf strings.Builder
			done := make(chan error, 1)

			go func() {
				_, err := io.Copy(&buf, r)
				r.Close()
				done <- err
			}()
			err = wslPathCmd.Run()
			if err != nil {
				panic(err)
			}
			w.Close()
			err = <-done
			wslPath := strings.TrimSpace(buf.String())

			bazelArgs = append(bazelArgs, "--build_event_binary_file="+wslPath)
			sleepAtEnd = true
		} else if strings.HasPrefix(arg, "--override_repository=intellij_aspect=") {
			var fileName = strings.Replace(arg, "--override_repository=intellij_aspect=", "", 1)

			var wslPathCmd = exec.Command("wsl", "wslpath", "-u", "-a", "'"+fileName+"'")
			r, w, err := os.Pipe()
			if err != nil {
				panic(err)
			}
			wslPathCmd.Stdout = w
			wslPathCmd.Stderr = os.Stderr
			var buf strings.Builder
			done := make(chan error, 1)

			go func() {
				_, err := io.Copy(&buf, r)
				r.Close()
				done <- err
			}()
			err = wslPathCmd.Run()
			if err != nil {
				panic(err)
			}
			w.Close()
			err = <-done
			wslPath := strings.TrimSpace(buf.String())

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
		//fo.WriteString(hex.EncodeToString([]byte(arg)) + "\n")
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
	fo.WriteString("END\n")
	if sleepAtEnd {
		//time.Sleep(60 * time.Second)
	}
}
