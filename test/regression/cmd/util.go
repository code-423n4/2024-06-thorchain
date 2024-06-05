package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// trunk-ignore-all(golangci-lint/forcetypeassert)

// deepMerge merges two maps recursively - if there is an array in the map it will deep
// merge elements that match on keys in overrideKeys.
func deepMerge(a, b map[string]any, overrideKeys ...string) map[string]any {
	result := make(map[string]any)
	for k, v := range a {
		result[k] = v
	}

	for k := range b {
		switch v := b[k].(type) {
		case []any:
			if avx, ok := result[k]; ok {
				var av []any
				if av, ok = avx.([]any); ok {

					// deep merge if key in overrideKeys matches in the source and destination
					for _, ov := range overrideKeys {
						for i := range av {
							for j := range v {
								if av[i].(map[string]any)[ov] == v[j].(map[string]any)[ov] {
									av[i] = deepMerge(av[i].(map[string]any), v[j].(map[string]any), overrideKeys...)

									// remove value from v
									v = append(v[:j], v[j+1:]...)
									break
								}
							}
						}
					}

					result[k] = append(av, v...)
					continue
				}
			}
		case map[string]any:
			if avx, ok := result[k]; ok {
				var av map[string]any
				if av, ok = avx.(map[string]any); ok {
					result[k] = deepMerge(av, v, overrideKeys...)
					continue
				}
			}
		}
		result[k] = b[k]
	}
	return result
}

func dumpLogs(out io.Writer, logs chan string) {
	for {
		select {
		case line := <-logs:
			_, _ = out.Write([]byte(line + "\n"))
			continue
		default:
		}
		break
	}
}

func drainLogs(logs chan string) {
	// if DEBUG is set skip draining logs
	if os.Getenv("DEBUG") != "" {
		return
	}

	for {
		select {
		case <-logs:
			continue
		case <-time.After(100 * time.Millisecond):
		}
		break
	}
}

func processRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(os.Signal(nil))
	return err == nil
}

func getTimeFactor() time.Duration {
	tf, err := strconv.ParseInt(os.Getenv("TIME_FACTOR"), 10, 64)
	if err != nil {
		return time.Duration(1)
	}
	return time.Duration(tf)
}

func consoleLogger(w io.Writer) zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: w}).With().Timestamp().Caller().Logger()
}

var reSetVar = regexp.MustCompile(`\$\{([A-Z0-9_]+)=([^}]+)\}`)

func parseOps(localLog zerolog.Logger, path string, tmpls *template.Template, env []string) (ops []Operation, opLines []int, newEnv []string) {
	// read the file
	f, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open test file")
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read test file")
	}
	f.Close()

	// track line numbers
	opLines = []int{0}
	scanner := bufio.NewScanner(bytes.NewBuffer(fileBytes))
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if line == "---" {
			opLines = append(opLines, i+2)
		}
	}

	// parse the template
	tmpl, err := tmpls.Parse(string(fileBytes))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse template")
	}

	// render the template
	tmplBuf := &bytes.Buffer{}
	err = tmpl.Execute(tmplBuf, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to render template")
	}

	// scan all lines in buffer for variable expansion
	buf := &bytes.Buffer{}
	scanner = bufio.NewScanner(tmplBuf)
	vars := map[string]string{}
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		matches := reSetVar.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			vars[match[1]] = match[2]
		}

		// regex replace variables and set variables
		for k, v := range vars {
			re := regexp.MustCompile(`\$\{` + k + `(=([^}]+))?\}`)
			line = re.ReplaceAllString(line, v)
		}

		// write line
		buf.WriteString(line + "\n")
	}

	// track whether we complete seen state and env operations
	stateComplete := false
	envComplete := false
	envAdded := 0

	dec := yaml.NewDecoder(buf)
	for {
		// decode into temporary type
		op := map[string]any{}
		err = dec.Decode(&op)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal().Err(err).Msg("failed to decode operation")
		}

		// warn empty operations
		if len(op) == 0 {
			localLog.Warn().Msg("empty operation, line numbers may be wrong")
			continue
		}

		// env operations are special
		if op["type"] == "env" {
			if envComplete {
				log.Fatal().Msg("env operations must be first")
			}
			o := NewOperation(op)
			env = append(env, fmt.Sprintf("%s=%s", o.(*OpEnv).Key, o.(*OpEnv).Value))
			envAdded++
			continue
		}
		envComplete = true

		// state operations must be first
		if op["type"] == "state" && stateComplete {
			log.Fatal().Msg("state operations must come before all operations other than env")
		}
		if op["type"] != "state" {
			stateComplete = true
		}

		ops = append(ops, NewOperation(op))
	}

	// remove env operations from op lines so numbers are correct
	opLines = opLines[envAdded:]

	return ops, opLines, env
}

func blockCount(ops []Operation) int {
	blocks := 0
	for _, op := range ops {
		if _, ok := op.(*OpCreateBlocks); ok {
			blocks += op.(*OpCreateBlocks).Count
		}
	}
	return blocks
}
