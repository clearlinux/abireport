//
// Copyright © Intel Corporation
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
//

package libabi

import (
	"debug/elf"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// A Report is used to traverse a given tree and identify any and all files
// that seem "interesting".
type Report struct {
	Root   string                        // Root directory that we're scanning
	Arches map[elf.Machine]*Architecture // Mapping of architectures

	wg        *sync.WaitGroup // Our wait group for multiprocessing
	jobChan   chan *Record    // Jobs are pushed from the walker
	storeChan chan *Record    // Single channel pulls all of the processed records
	libDirs   []string        // Valid library directories
	nRecords  int             // The total number of records encountered
	jobMutex  *sync.Mutex     // Lock for njobs decrement
	nJobs     int             // Number of jobs, decrementing through runtime
}

// NewReport will create a new walker instance, and attempt to initialise
// the magic library
func NewReport(root string) (*Report, error) {
	// Only want the absolute path here.
	if absPath, err := filepath.Abs(root); err == nil {
		root = absPath
	} else {
		return nil, err
	}
	libDirs := []string{
		filepath.Join(root, "usr", "lib64"),
		filepath.Join(root, "usr", "lib"),
		filepath.Join(root, "usr", "lib32"),
		filepath.Join(root, "usr", "lib", "x86_64-linux-gnu"),
		filepath.Join(root, "lib", "x86_64-linux-gnu"),
		filepath.Join(root, "usr", "lib", "i386-linux-gnu"),
		filepath.Join(root, "lib", "i386-linux-gnu"),
	}
	return &Report{
		Root:      root,
		jobChan:   make(chan *Record),
		storeChan: make(chan *Record),
		wg:        new(sync.WaitGroup),
		Arches:    make(map[elf.Machine]*Architecture),
		libDirs:   libDirs,
		nRecords:  0,
		jobMutex:  new(sync.Mutex),
		nJobs:     runtime.NumCPU(),
	}, nil
}

// IsAnELF determines if a file is an ELF file or not
// Loose reinterpetation of debug/elf magic checking
func IsAnELF(p string) (bool, error) {
	o, err := os.Open(p)
	if err != nil {
		return false, err
	}
	defer o.Close()
	var h [16]uint8
	if _, err := o.ReadAt(h[0:], 0); err != nil {
		return false, err
	}

	if h[0] == '\x7f' && h[1] == 'E' && h[2] == 'L' && h[3] == 'F' {
		return true, nil
	}
	return false, nil
}

// recordPath will make a record of the "interesting" path if it determines
// this path to be of value: i.e. a dynamic ELF file.
func (a *Report) recordPath(path string) *Record {
	isElf, err := IsAnELF(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking file %v: %v\n", path, err)
		return nil
	}
	if !isElf {
		return nil
	}

	record := &Record{Path: path}

	return record
}

// isPathInteresting determines whether we care about a path or not. Simply,
// is it the right *type*.
func (a *Report) isPathInteresting(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}
	// Also care not for symlinks.
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		return false
	}
	// Really don't want to examine .debug files
	if strings.Contains(path, "/debug/") && strings.HasSuffix(path, ".debug") {
		return false
	}
	return true
}

// walkTree is our callback for filepath.Walk, to actually do the initial
// tree traversal.
func (a *Report) walkTree(path string, info os.FileInfo, err error) error {
	// Determine if the path is actually worth looking at
	if !a.isPathInteresting(path, info) {
		return nil
	}
	// Send the file off for further processing
	if record := a.recordPath(path); record != nil {
		a.jobChan <- record
	}
	return nil
}

// jobDone is used to close the channel once all jobs have finished
func (a *Report) jobDone() {
	a.jobMutex.Lock()
	defer a.jobMutex.Unlock()
	a.nJobs--

	if a.nJobs == 0 {
		close(a.storeChan)
	}
}

// the jobProcessor will keep pulling record jobs created from the walker,
// which has done preliminary groundwork based on the magic type strings.
func (a *Report) jobProcessor() {
	defer a.wg.Done()
	for i := 0; i < a.nJobs; i++ {
		go func() {
			defer a.wg.Done()
			for job := range a.jobChan {
				if err := a.AnalyzeOne(job); err != nil {
					fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", job.Path, err)
					continue
				}
				a.storeChan <- job
			}
			a.jobDone()
		}()
	}
}

// storeProcessor is responsible for storing into memory
func (a *Report) storeProcessor() {
	defer a.wg.Done()

	for record := range a.storeChan {
		var symbolsMap map[string]bool
		var ok bool
		a.nRecords++

		bucket := a.GetBucket(record)
		symbolsTgt := bucket.GetSymbolsTarget(record)

		// Ensure map is here so that the .soname provider is known
		if symbolsMap, ok = symbolsTgt[record.Name]; !ok {
			symbolsMap = make(map[string]bool)
			symbolsTgt[record.Name] = symbolsMap
		}

		// Store a soname -> symbol mapping
		// Quicker than actually using lists
		for _, symbol := range record.Symbols {
			symbolsMap[symbol] = true
		}

		for _, dep := range record.Dependencies {
			bucket.Dependencies[dep] = true
		}
	}
}

// Walk will attempt to walk the preconfigured tree, and collect a set
// of "interesting" files along the way.
func (a *Report) Walk() error {
	a.wg.Add(3 + a.nJobs)
	go a.jobProcessor()
	go a.storeProcessor()
	go func() {
		defer a.wg.Done()
		filepath.Walk(a.Root, a.walkTree)
		close(a.jobChan)
	}()
	a.wg.Wait()
	return nil
}
