package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor") {
			processFile(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking path: %v\n", err)
	}
}

func processFile(path string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return
	}

	var filterConfigVar *ast.ValueSpec
	var varName string

	// Find the variable annotated with @gen_swagger_filter
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		hasAnnotation := false
		if genDecl.Doc != nil {
			for _, comment := range genDecl.Doc.List {
				if strings.Contains(comment.Text, "@gen_swagger_filter") {
					hasAnnotation = true
					break
				}
			}
		}

		if hasAnnotation {
			for _, spec := range genDecl.Specs {
				vSpec, ok := spec.(*ast.ValueSpec)
				if ok && len(vSpec.Names) > 0 {
					filterConfigVar = vSpec
					varName = vSpec.Names[0].Name
					break
				}
			}
		}
	}

	if filterConfigVar == nil {
		return
	}

	fmt.Printf("Found @gen_swagger_filter for variable %s in %s\n", varName, path)

	// Extract AllowedFields and AllowedSortFields
	allowedFields, allowedSortFields := extractFilterInfo(filterConfigVar)

	// Generate Swagger documentation strings
	filterDoc := generateFilterDoc(allowedFields)
	orderByDoc := generateOrderByDoc(allowedSortFields)

	// Update the file
	updateSwaggerDocs(path, filterDoc, orderByDoc)
}

func extractFilterInfo(spec *ast.ValueSpec) (map[string][]string, []string) {
	allowedFields := make(map[string][]string)
	var allowedSortFields []string

	if len(spec.Values) == 0 {
		return allowedFields, allowedSortFields
	}

	compLit, ok := spec.Values[0].(*ast.CompositeLit)
	if !ok {
		return allowedFields, allowedSortFields
	}

	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		keyIdent, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		if keyIdent.Name == "AllowedFields" {
			fieldsLit, ok := kv.Value.(*ast.CompositeLit)
			if ok {
				for _, fieldElt := range fieldsLit.Elts {
					fieldKv, ok := fieldElt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					fieldName := strings.Trim(fieldKv.Key.(*ast.BasicLit).Value, "\"")

					var ops []string
					fieldConfigLit, ok := fieldKv.Value.(*ast.CompositeLit)
					if ok {
						for _, configElt := range fieldConfigLit.Elts {
							configKv, ok := configElt.(*ast.KeyValueExpr)
							if ok && configKv.Key.(*ast.Ident).Name == "AllowedOperators" {
								opsLit, ok := configKv.Value.(*ast.CompositeLit)
								if ok {
									for _, opElt := range opsLit.Elts {
										ops = append(ops, strings.Trim(opElt.(*ast.BasicLit).Value, "\""))
									}
								}
							}
						}
					}
					allowedFields[fieldName] = ops
				}
			}
		} else if keyIdent.Name == "AllowedSortFields" {
			sortLit, ok := kv.Value.(*ast.CompositeLit)
			if ok {
				for _, sortElt := range sortLit.Elts {
					allowedSortFields = append(allowedSortFields, strings.Trim(sortElt.(*ast.BasicLit).Value, "\""))
				}
			}
		}
	}

	return allowedFields, allowedSortFields
}

func generateFilterDoc(fields map[string][]string) string {
	var sb strings.Builder
	sb.WriteString("// @Param filter query string false \"Filter expression. \\n- Allowed fields & ops:\\n")

	var sortedKeys []string
	for k := range fields {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		ops := fields[k]
		sb.WriteString(fmt.Sprintf("  - %s: %s\\n", k, strings.Join(ops, ", ")))
	}
	sb.WriteString("\"")
	return sb.String()
}

func generateOrderByDoc(fields []string) string {
	var sb strings.Builder
	sb.WriteString("// @Param order_by query string false \"Sort field. \\n- Allowed: ")
	sb.WriteString(strings.Join(fields, ", "))
	sb.WriteString("\"")

	if len(fields) > 0 {
		sb.WriteString(fmt.Sprintf(" example(%s:asc)", fields[0]))
	}

	return sb.String()
}

func updateSwaggerDocs(path, filterDoc, orderByDoc string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "// @Param filter query string") {
			newLines = append(newLines, filterDoc)
		} else if strings.HasPrefix(trimmed, "// @Param order_by query string") {
			newLines = append(newLines, orderByDoc)
		} else {
			newLines = append(newLines, line)
		}
	}

	err = os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		fmt.Printf("Error writing file %s: %v\n", path, err)
	}
}
