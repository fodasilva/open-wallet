package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type FieldMapping struct {
	Column     string
	GoField    string
	IsOptional bool
}

type MethodConfig struct {
	Name       string
	Operation  string
	Fields     []FieldMapping
	Payload    string
	Alias      string
	ReturnType string
	Repo       *RepoConfig
}

type RepoConfig struct {
	Table    string
	Entity   string
	RepoName string
	Methods  []MethodConfig
	Package  string
	Dir      string
}

var (
	reRepo    = regexp.MustCompile(`@gen_repo`)
	reTable   = regexp.MustCompile(`@table:\s*([^\s|]+)`)
	reEntity  = regexp.MustCompile(`@entity:\s*([^\s|]+)`)
	reName    = regexp.MustCompile(`@name:\s*([^\s|]+)`)
	reMethod  = regexp.MustCompile(`@method:\s*([^|]+)`)
	reFields  = regexp.MustCompile(`fields:\s*([^|]+)`)
	rePayload = regexp.MustCompile(`payload:\s*([^|]+)`)
	reAlias   = regexp.MustCompile(`alias:\s*([^|]+)`)
	reStruct  = regexp.MustCompile(`type\s+([^\s]+)\s+struct`)
)

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor") || strings.Contains(path, "scripts") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if !reRepo.Match(content) {
			return nil
		}

		processFile(path)
		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func processFile(path string) {
	fmt.Printf("--- Generating repository for %s ---\n", path)
	dir := filepath.Dir(path)

	// Clean old generated files
	files, _ := filepath.Glob(filepath.Join(dir, "zz_generated_*.go"))
	for _, f := range files {
		_ = os.Remove(f)
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	packageName := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "package ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				packageName = fields[1]
			}
			break
		}
	}

	var repos []*RepoConfig
	for i := 0; i < len(lines); i++ {
		if reRepo.MatchString(lines[i]) {
			blockEnd := len(lines)
			for j := i + 1; j < len(lines); j++ {
				if reRepo.MatchString(lines[j]) {
					blockEnd = j
					break
				}
			}

			repo := &RepoConfig{
				Package: packageName,
				Dir:     dir,
			}

			startSearch := i - 5
			if startSearch < 0 {
				startSearch = 0
			}

			repoContent := strings.Join(lines[startSearch:blockEnd], "\n")

			if m := reTable.FindStringSubmatch(repoContent); len(m) > 1 {
				repo.Table = m[1]
			}
			if m := reEntity.FindStringSubmatch(repoContent); len(m) > 1 {
				repo.Entity = m[1]
			}
			if m := reName.FindStringSubmatch(repoContent); len(m) > 1 {
				repo.RepoName = m[1]
			}

			for j := i; j < blockEnd; j++ {
				if m := reStruct.FindStringSubmatch(lines[j]); len(m) > 1 {
					if repo.Entity == "" {
						repo.Entity = m[1]
					}
					break
				}
			}

			if repo.RepoName == "" {
				repo.RepoName = repo.Entity + "Repo"
			}

			for j := i; j < blockEnd; j++ {
				if m := reMethod.FindStringSubmatch(lines[j]); len(m) > 1 {
					method := MethodConfig{
						Name: strings.TrimSpace(m[1]),
						Repo: repo,
					}

					methodLine := lines[j]
					if m := reFields.FindStringSubmatch(methodLine); len(m) > 1 {
						method.Fields = parseFields(m[1])
					}
					if m := rePayload.FindStringSubmatch(methodLine); len(m) > 1 {
						method.Payload = strings.TrimSpace(m[1])
					}
					if m := reAlias.FindStringSubmatch(methodLine); len(m) > 1 {
						method.Alias = strings.TrimSpace(m[1])
					} else {
						method.Alias = method.Name
					}

					name := method.Name
					switch {
					case strings.HasPrefix(name, "Select"):
						method.Operation = "Select"
					case strings.HasPrefix(name, "Insert"):
						method.Operation = "Insert"
					case strings.HasPrefix(name, "Update"):
						method.Operation = "Update"
					case strings.HasPrefix(name, "Delete"):
						method.Operation = "Delete"
					case strings.HasPrefix(name, "Count"):
						method.Operation = "Count"
					case strings.HasPrefix(name, "GetByID"):
						method.Operation = "GetByID"
					default:
						method.Operation = name
					}

					switch method.Operation {
					case "Select":
						method.ReturnType = "[]" + repo.Entity
					case "Insert", "Update", "Delete":
						method.ReturnType = "error"
					case "Count":
						method.ReturnType = "int"
					case "GetByID":
						method.ReturnType = repo.Entity
					}

					repo.Methods = append(repo.Methods, method)
				}
			}

			repos = append(repos, repo)
			i = blockEnd - 1
		}
	}

	generateFiles(repos)
}

func parseFields(s string) []FieldMapping {
	var fields []FieldMapping
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		kv := strings.Split(p, ":")
		if len(kv) == 2 {
			col := strings.TrimSpace(kv[0])
			goFieldRaw := strings.TrimSpace(kv[1])
			isOptional := strings.HasSuffix(goFieldRaw, "?")
			goField := strings.TrimSuffix(goFieldRaw, "?")
			fields = append(fields, FieldMapping{
				Column:     col,
				GoField:    goField,
				IsOptional: isOptional,
			})
		}
	}
	return fields
}

func generateFiles(repos []*RepoConfig) {
	if len(repos) == 0 {
		return
	}

	methodsByName := make(map[string][]MethodConfig)
	for _, repo := range repos {
		for _, method := range repo.Methods {
			name := strings.ToLower(method.Name)
			methodsByName[name] = append(methodsByName[name], method)
		}
	}

	for name, methods := range methodsByName {
		filename := filepath.Join(repos[0].Dir, fmt.Sprintf("zz_generated_%s.go", name))

		var sb strings.Builder
		sb.WriteString("// Code generated. DO NOT EDIT.\n\n")
		sb.WriteString(fmt.Sprintf("package %s\n\n", repos[0].Package))
		sb.WriteString("import (\n")
		sb.WriteString("\t\"context\"\n")
		sb.WriteString("\t\"github.com/Masterminds/squirrel\"\n")
		sb.WriteString("\t\"github.com/felipe1496/open-wallet/internal/util\"\n")
		if name != "insert" {
			sb.WriteString("\t\"github.com/felipe1496/open-wallet/internal/util/querybuilder\"\n")
		}
		sb.WriteString(")\n\n")

		for _, method := range methods {
			repo := method.Repo
			templatePath := filepath.Join("templates/repository", strings.ToLower(method.Operation)+".txt")
			content, err := os.ReadFile(templatePath)
			if err != nil {
				fmt.Printf("Warning: Template not found %s\n", templatePath)
				continue
			}

			tpl := string(content)
			tpl = strings.ReplaceAll(tpl, "{{RepoName}}", repo.RepoName)
			tpl = strings.ReplaceAll(tpl, "{{MethodName}}", method.Alias)
			tpl = strings.ReplaceAll(tpl, "{{TableName}}", repo.Table)
			tpl = strings.ReplaceAll(tpl, "{{StructName}}", repo.Entity)
			tpl = strings.ReplaceAll(tpl, "{{ReturnType}}", method.ReturnType)
			tpl = strings.ReplaceAll(tpl, "{{PayloadType}}", method.Payload)

			switch method.Operation {
			case "Select":
				cols := []string{}
				scans := []string{}
				for _, f := range method.Fields {
					cols = append(cols, f.Column)
					scans = append(scans, "&item."+f.GoField)
				}
				tpl = strings.ReplaceAll(tpl, "{{Columns}}", `"`+strings.Join(cols, `", "`)+`"`)
				tpl = strings.ReplaceAll(tpl, "{{ScanFields}}", strings.Join(scans, ",\n\t\t\t")+",")

			case "Insert":
				inColsVals := "var columns []string\n\tvar values []interface{}\n"
				for _, f := range method.Fields {
					if f.IsOptional {
						inColsVals += fmt.Sprintf("\tif data.%s.Set { columns = append(columns, \"%s\"); values = append(values, data.%s.Value) }\n", f.GoField, f.Column, f.GoField)
					} else {
						inColsVals += fmt.Sprintf("\tcolumns = append(columns, \"%s\"); values = append(values, data.%s)\n", f.Column, f.GoField)
					}
				}
				inColsVals += "\tquery = query.Columns(columns...).Values(values...)"
				tpl = strings.ReplaceAll(tpl, "{{InColsVals}}", inColsVals)

			case "Update":
				sets := ""
				for _, f := range method.Fields {
					if f.IsOptional {
						sets += fmt.Sprintf("\tif data.%s.Set { query = query.Set(\"%s\", data.%s.Value) }\n", f.GoField, f.Column, f.GoField)
					} else {
						sets += fmt.Sprintf("\tquery = query.Set(\"%s\", data.%s)\n", f.Column, f.GoField)
					}
				}
				tpl = strings.ReplaceAll(tpl, "{{UpdateSets}}", sets)
			}

			sb.WriteString(tpl)
			sb.WriteString("\n\n")
		}

		if err := os.WriteFile(filename, []byte(sb.String()), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filename, err)
			continue
		}
		_ = exec.Command("gofmt", "-w", filename).Run()
		fmt.Printf("    - Generated %s\n", filename)
	}
}
