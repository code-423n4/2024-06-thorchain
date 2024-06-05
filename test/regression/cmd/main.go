package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"text/template"

	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////////////
// Main
////////////////////////////////////////////////////////////////////////////////////////

func main() {
	// parse the regex in the RUN environment variable to determine which tests to run
	runRegex := regexp.MustCompile(".*")
	if len(os.Getenv("RUN")) > 0 {
		runRegex = regexp.MustCompile(os.Getenv("RUN"))
	}

	// find all regression tests in path
	files := []string{}
	err := filepath.Walk("suites", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// skip files that are not yaml
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		if runRegex.MatchString(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to find regression tests")
	}

	// sort the files descending by the number of blocks created (so long tests run first)
	counts := make(map[string]int)
	for _, file := range files {
		ops, _, _ := parseOps(log.Output(io.Discard), file, template.Must(templates.Clone()), []string{})
		counts[file] = blockCount(ops)
	}
	sort.Slice(files, func(i, j int) bool {
		return counts[files[i]] > counts[files[j]]
	})

	// keep track of the results
	mu := sync.Mutex{}
	succeeded := []string{}
	failed := []string{}
	completed := 0
	total := len(files)

	// get parallelism from environment variable if DEBUG is not set
	parallelism := 1
	sem := make(chan struct{}, 1)
	wg := sync.WaitGroup{}
	if len(os.Getenv("PARALLELISM")) > 0 && len(os.Getenv("DEBUG")) == 0 {
		parallelism, err = strconv.Atoi(os.Getenv("PARALLELISM"))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse PARALLELISM")
		}
		sem = make(chan struct{}, parallelism)
	}
	log.Info().
		Int("parallelism", parallelism).
		Int("count", total).
		Msg("running tests")
	fmt.Println() // A blank line before the first regression test whether parallel or not.

	// use a channel to abort early if there is a failure in a merge request run
	abort := make(chan struct{})
	aborted := func() bool {
		select {
		case <-abort:
			return true
		default:
			return false
		}
	}

	// run tests
	for i, file := range files {
		// break if aborted
		if aborted() {
			break
		}

		sem <- struct{}{}
		wg.Add(1)
		go func(routine int, file string) {
			// create home directory
			home := "/" + strconv.Itoa(routine)
			_ = os.MkdirAll(home, 0o755)

			// create a buffer to capture the logs
			var out io.Writer = os.Stderr
			buf := new(bytes.Buffer)
			if parallelism > 1 {
				out = buf
			}

			// release semaphore and wait group
			defer func() {
				<-sem
				wg.Done()

				// write buffer to outputs
				mu.Lock()
				completed++
				localLog := consoleLogger(out)
				localLog.Info().Msg(fmt.Sprintf("%d/%d regression tests completed.", completed, total))
				if parallelism > 1 {
					fmt.Print(buf.String())
				}
				fmt.Println() // Blank line separating regression tests.
				mu.Unlock()
			}()

			// run test
			failExportInvariants, runErr := run(out, file, routine)
			if runErr != nil {
				mu.Lock()
				failed = append(failed, file)
				if os.Getenv("FAIL_FAST") != "" {
					close(abort)
				}
				mu.Unlock()
				return
			}

			// check export state
			exportErr := export(out, file, routine, failExportInvariants)
			if exportErr != nil {
				mu.Lock()
				failed = append(failed, file)
				if os.Getenv("FAIL_FAST") != "" {
					close(abort)
				}
				mu.Unlock()
				return
			}

			// success
			mu.Lock()
			succeeded = append(succeeded, file)
			mu.Unlock()
		}(i, file)
	}

	// wait for all tests to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// wait for all tests to finish or abort
	select {
	case <-done:
	case <-abort:
		fmt.Printf("%s>> FAIL_FAST: Aborting Now <<%s\n", ColorRed, ColorReset)
	}

	// lock in case this was early abort, no need to unlock in main
	mu.Lock()

	// print the results
	fmt.Printf("%sSucceeded:%s %d\n", ColorGreen, ColorReset, len(succeeded))
	for _, file := range succeeded {
		fmt.Printf("- %s\n", file)
	}
	fmt.Printf("%sFailed:%s %d\n", ColorRed, ColorReset, len(failed))
	for _, file := range failed {
		fmt.Printf("- %s\n", file)
	}
	fmt.Println()

	// exit with error code if any tests failed
	if len(failed) > 0 {
		os.Exit(1)
	}
}
