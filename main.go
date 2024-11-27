package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var subcommand string
	if len(os.Args) > 1 {
		subcommand = os.Args[1]
	}
	switch subcommand {
	case "generate":
		generateLogs()
	case "analyze":
		analyzeEvents()
	default:
		panic("unknown subcommand")
	}
}

func generateLogs() {
	logsDir := "logs"
	err := os.Mkdir(logsDir, 0755)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		logFile := filepath.Join(logsDir, fmt.Sprintf("file-%d.log", i))
		fp, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		for j := 0; j < 1000; j++ {
			_, err := fp.WriteString(fmt.Sprintf("file-%d-line-%d\n", i, j))
			if err != nil {
				_ = fp.Close()
				panic(err)
			}
		}
		_ = fp.Close()
	}
}

type logEntry struct {
	Message string `json:"message"`
}

func analyzeEvents() {
	logLines := make(map[string]bool, 1000*1000)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(path, "events") && strings.HasSuffix(path, ".ndjson") {
			fp, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("could not open file %s: %w", path, err)
			}
			defer fp.Close()
			scanner := bufio.NewScanner(fp)
			for scanner.Scan() {
				var entry logEntry
				err := json.Unmarshal(scanner.Bytes(), &entry)
				if err != nil {
					return fmt.Errorf("could not parse file %s: %w", path, err)
				}
				if !strings.HasPrefix(entry.Message, "file-") {
					return fmt.Errorf("unknown message: %s", entry.Message)
				}
				_, ok := logLines[entry.Message]
				if ok {
					return fmt.Errorf("duplicated log entry %s", entry.Message)
				}
				logLines[entry.Message] = true
			}
			return nil
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	if len(logLines) != 1000*1000 {
		panic(fmt.Sprintf("expected %d lines, got %d", 1000*1000, len(logLines)))
	}
}
