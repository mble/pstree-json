package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Proc represents an entry in /proc
type Proc struct {
	PID      string  `json:"pid"`
	Command  string  `json:"command"`
	Status   string  `json:"status"`
	PPID     string  `json:"ppid"`
	Children []*Proc `json:"children"`
}

func handleErr(err error) {
	errExit := 1

	fmt.Printf(`{"error": "%+v"}\n`, err)
	os.Exit(errExit)
}

func readProc(path string) (*Proc, error) {
	root := filepath.Dir(path)
	f, err := os.Open(path)
	if err != nil {
		return &Proc{}, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return &Proc{}, err
	}

	vals := strings.Split(string(data), " ")
	proc := &Proc{
		PID:      vals[0],
		Status:   vals[2],
		PPID:     vals[3],
		Children: []*Proc{},
	}

	command, err := readCommand(root + "/cmdline")
	if err != nil {
		return &Proc{}, err
	}

	proc.Command = command

	return proc, nil
}

func readCommand(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	vals := strings.Split(string(data), "\x00")
	cmd := strings.TrimSpace(strings.Join(vals, " "))

	return cmd, nil
}

func associateChildren(proc *Proc, candidates []*Proc) {
	for _, candidate := range candidates {
		if candidate.PPID == proc.PID {
			proc.Children = append(proc.Children, candidate)

			for index := range candidates {
				candidates = candidates[:index]
			}
		}
	}

	for _, child := range proc.Children {
		associateChildren(child, candidates)
	}
}

func buildTree(procPrefix string) ([]*Proc, error) {
	pattern := procPrefix + "/[0-9]*/stat"

	procs := []*Proc{}
	files, err := filepath.Glob(pattern)

	if err != nil {
		return procs, err
	}

	for _, file := range files {
		proc, err := readProc(file)
		if err != nil {
			return procs, err
		}

		procs = append(procs, proc)
	}

	for _, proc := range procs {
		associateChildren(proc, procs)
	}

	return procs, err
}

func marshalTree(procs []*Proc, targetPID string) (string, error) {
	output := ""
	for _, proc := range procs {
		if proc.PID == targetPID {
			output, err := json.Marshal(proc)
			if err != nil {
				return "", err
			}
			return string(output), nil
		}
	}
	return output, nil
}

func main() {
	procPrefix := "/proc"
	targetPID := ""
	flag.StringVar(&targetPID, "pid", "1", "target PID for the root of the tree")
	flag.Parse()

	procs, err := buildTree(procPrefix)
	if err != nil {
		handleErr(err)
	}

	output, err := marshalTree(procs, targetPID)
	if err != nil {
		handleErr(err)
	}
	fmt.Printf("%s", output)
}
