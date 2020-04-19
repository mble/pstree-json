package main

import (
	"reflect"
	"testing"
)

func TestBuildTree(t *testing.T) {
	procPrefix := "./testdata"
	procs, err := buildTree(procPrefix)
	if err != nil {
		t.Error(err)
	}
	if len(procs) != 4 {
		t.Errorf("expected %d, got %d", 4, len(procs))
	}
}

func TestReadProc(t *testing.T) {
	testCases := []struct {
		desc         string
		filepath     string
		expectError  bool
		expectedProc *Proc
	}{
		{
			desc:        "file present",
			filepath:    "testdata/14/stat",
			expectError: false,
			expectedProc: &Proc{
				PID:      "14",
				Command:  "sleep 1000",
				Status:   "S",
				PPID:     "1",
				Children: []*Proc{},
			},
		},
		{
			desc:        "file missing",
			filepath:    "testdata/15/stat",
			expectError: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			proc, err := readProc(tC.filepath)
			if !tC.expectError {
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				if !reflect.DeepEqual(tC.expectedProc, proc) {
					t.Errorf("expected %+v, got %+v", tC.expectedProc, proc)
				}
			}
			if tC.expectError && err == nil {
				t.Errorf("expected error")
			}
		})
	}
}

func TestMarshalTree(t *testing.T) {
	procPrefix := "./testdata"
	procs, err := buildTree(procPrefix)
	if err != nil {
		t.Error(err)
	}
	testCases := []struct {
		desc      string
		targetPID string
		expected  string
	}{
		{
			desc:      "PID 1",
			targetPID: "1",
			expected:  `{"pid":"1","command":"/bin/bash -c sleep","status":"S","ppid":"0","children":[{"pid":"14","command":"sleep 1000","status":"S","ppid":"1","children":[]},{"pid":"22","command":"sleep 1000","status":"S","ppid":"1","children":[{"pid":"45","command":"sleep 1000","status":"S","ppid":"22","children":[]}]}]}`,
		},
		{
			desc:      "PID 14",
			targetPID: "14",
			expected:  `{"pid":"14","command":"sleep 1000","status":"S","ppid":"1","children":[]}`,
		},
		{
			desc:      "PID 45",
			targetPID: "45",
			expected:  `{"pid":"45","command":"sleep 1000","status":"S","ppid":"22","children":[]}`,
		},
		{
			desc:      "PID 22",
			targetPID: "22",
			expected:  `{"pid":"22","command":"sleep 1000","status":"S","ppid":"1","children":[{"pid":"45","command":"sleep 1000","status":"S","ppid":"22","children":[]}]}`,
		},
		{
			desc:      "non-existant PID",
			targetPID: "100",
			expected:  "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output, err := marshalTree(procs, tC.targetPID)
			if err != nil {
				t.Error(err)
			}
			if output != tC.expected {
				t.Errorf("expected %s, got %s", tC.expected, output)
			}
		})
	}
}
