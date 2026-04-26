package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reRouter = regexp.MustCompile(`@Router`)
	reID     = regexp.MustCompile(`@ID\s+([^\s]+)`)
)

func main() {
	searchDir := "internal/resources"
	errorFound := false
	idsSeen := make(map[string]string) // ID -> FilePath

	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer func() { _ = file.Close() }()

		routerCount := 0
		idCount := 0
		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if reRouter.MatchString(line) {
				routerCount++
			}

			if m := reID.FindStringSubmatch(line); len(m) > 1 {
				idCount++
				id := m[1]
				if firstFile, seen := idsSeen[id]; seen {
					fmt.Printf("❌ ERROR: Duplicate @ID '%s' found in %s (first seen in %s)\n", id, path, firstFile)
					errorFound = true
				}
				idsSeen[id] = path
			}
		}

		if routerCount != idCount {
			fmt.Printf("❌ ERROR: %s has %d @Router(s) but %d @ID(s)\n", path, routerCount, idCount)
			errorFound = true
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking path: %v\n", err)
		os.Exit(1)
	}

	if errorFound {
		fmt.Println("Summary: Swagger documentation issues found. Please fix @ID tags.")
		os.Exit(1)
	}

	fmt.Println("✅ All endpoints have matching and unique @ID tags.")
}
