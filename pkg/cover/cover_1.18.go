// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build go1.18

package cover

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// FlushProfiles flushes test profiles to disk. It works by build and executing
// a dummy list of 1 test. This is to ensure we execute the M.after() function
// (a function internal to the testing package) where all reporting (cpu, mem,
// cover, ... profiles) is flushed to disk.
func FlushProfiles() {
	// Redirect Stdout/err temporarily so the testing code doesn't output the
	// regular:
	//   PASS
	//   coverage: 21.4% of statements
	// Thanks to this, we can test the output of the instrumented binary the same
	// way the normal binary is.
	oldstdout := os.Stdout
	oldstderr := os.Stderr
	os.Stdout, _ = os.Open(os.DevNull)
	os.Stderr, _ = os.Open(os.DevNull)

	tests := []testing.InternalTest{}
	benchmarks := []testing.InternalBenchmark{}
	fuzzes := []testing.InternalFuzzTarget{}
	examples := []testing.InternalExample{}
	var f dummyTestDeps
	dummyM := testing.MainStart(f, tests, benchmarks, fuzzes, examples)
	dummyM.Run()

	// restore stdout/err
	os.Stdout = oldstdout
	os.Stderr = oldstderr
}

// lifted from go's internal/fuzz package
type corpusEntry = struct {
	Parent     string
	Path       string
	Data       []byte
	Values     []interface{}
	Generation int
	IsSeed     bool
}

func (f dummyTestDeps) CheckCorpus(vs []interface{}, ts []reflect.Type) error { return nil }
func (f dummyTestDeps) ReadCorpus(dir string, types []reflect.Type) ([]corpusEntry, error) {
	return nil, nil
}
func (f dummyTestDeps) ResetCoverage()                              {}
func (f dummyTestDeps) RunFuzzWorker(func(corpusEntry) error) error { return nil }
func (f dummyTestDeps) SnapshotCoverage()                           {}

func (f dummyTestDeps) CoordinateFuzzing(
	timeout time.Duration,
	limit int64,
	minimizeTimeout time.Duration,
	minimizeLimit int64,
	parallel int,
	seed []corpusEntry,
	types []reflect.Type,
	corpusDir,
	cacheDir string) error {
	return nil
}