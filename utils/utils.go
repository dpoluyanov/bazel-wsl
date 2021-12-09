package utils

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func WinToWSLPath(fileName string) string {
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
	return wslPath
}

func WSLToWinPath(fileName string) string {
	var wslPathCmd = exec.Command("wsl", "wslpath", "-a", "-w", "'"+fileName+"'")
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
	return wslPath
}
