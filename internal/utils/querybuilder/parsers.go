package querybuilder

import (
	"fmt"
	"strconv"
	"strings"
)

type Results struct {
	Builder *Builder
	Page    int
	PerPage int
}

type FieldConfig struct {
	AllowedOperators []string
}

type ParseConfig struct {
	AllowedFields     map[string]FieldConfig
	AllowedSortFields []string
}

func ParseRequest(filter, pageStr, perPageStr, orderBy string, config *ParseConfig) (*Results, error) {
	pageNum, err := strconv.Atoi(pageStr)
	if err != nil {
		pageNum = 1
	}

	perPageNum, err := strconv.Atoi(perPageStr)
	if err != nil {
		perPageNum = 10
	}

	builder := New()
	builder.Offset((pageNum - 1) * perPageNum)
	builder.Limit(perPageNum + 1)

	if orderBy != "" {
		splittedByComma := strings.Split(orderBy, ",")

		for _, field := range splittedByComma {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}

			splittedByColon := strings.Split(field, ":")
			fieldName := strings.TrimSpace(splittedByColon[0])

			if fieldName == "" {
				return nil, fmt.Errorf("invalid order_by: empty field name")
			}

			// Validate sort field if config is provided
			if config != nil && len(config.AllowedSortFields) > 0 {
				found := false
				for _, f := range config.AllowedSortFields {
					if f == fieldName {
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("field '%s' not allowed for sorting", fieldName)
				}
			}

			direction := "asc"
			if len(splittedByColon) > 1 {
				direction = strings.ToLower(strings.TrimSpace(splittedByColon[1]))
				if direction != "asc" && direction != "desc" {
					return nil, fmt.Errorf("invalid order: must be 'asc' or 'desc'")
				}
			}

			builder.OrderBy(fieldName, direction)
		}
	}

	if filter != "" {
		if _, err := parseFilter(filter, builder, config); err != nil {
			return nil, err
		}
	}

	return &Results{
		Builder: builder,
		Page:    pageNum,
		PerPage: perPageNum,
	}, nil
}

func parseFilter(filter string, b *Builder, config *ParseConfig) (*Builder, error) {
	splitted := splitByDelimiterOutsideQuotesAndParens(filter, " and ")

	allowedOperators := map[string]bool{
		"eq":   true,
		"ne":   true,
		"gt":   true,
		"gte":  true,
		"lt":   true,
		"lte":  true,
		"like": true,
	}

	for _, filterPart := range splitted {
		isOrGroup := strings.HasPrefix(filterPart, "(") && strings.HasSuffix(filterPart, ")")

		if isOrGroup {
			trimmedFilter := strings.TrimPrefix(filterPart, "(")
			trimmedFilter = strings.TrimSuffix(trimmedFilter, ")")

			splittedOrGroup := splitByDelimiterOutsideQuotesAndParens(trimmedFilter, " or ")

			orQuery := b.InitOr()

			for _, filterSubPart := range splittedOrGroup {
				splittedFilter, err := splitFilter(filterSubPart)

				if err != nil {
					return nil, err
				}
				if len(splittedFilter) != 3 {
					return nil, fmt.Errorf("filter syntax error at '%s'", filterSubPart)
				}

				field := splittedFilter[0]
				operator := splittedFilter[1]
				valueRaw := splittedFilter[2]

				// Validate field and operator if config is provided
				if config != nil && config.AllowedFields != nil {
					fieldConfig, allowed := config.AllowedFields[field]
					if !allowed {
						return nil, fmt.Errorf("field '%s' not allowed for filtering", field)
					}
					if len(fieldConfig.AllowedOperators) > 0 {
						found := false
						for _, op := range fieldConfig.AllowedOperators {
							if op == operator {
								found = true
								break
							}
						}
						if !found {
							return nil, fmt.Errorf("operator '%s' not allowed for field '%s'", operator, field)
						}
					}
				}

				if !allowedOperators[operator] {
					return nil, fmt.Errorf("operator '%s' not allowed at '%s'", operator, filterSubPart)
				}

				condValue, err := parseFilterValue(valueRaw)
				if err != nil {
					return nil, err
				}

				orQuery.Or(field, operator, condValue)
			}

			orQuery.EndOr()
		} else {
			splittedFilter, err := splitFilter(filterPart)
			if err != nil {
				return nil, err
			}

			if len(splittedFilter) != 3 {
				return nil, fmt.Errorf("filter syntax error at '%s'", filterPart)
			}

			field := splittedFilter[0]
			operator := splittedFilter[1]
			valueRaw := splittedFilter[2]

			// Validate field and operator if config is provided
			if config != nil && config.AllowedFields != nil {
				fieldConfig, allowed := config.AllowedFields[field]
				if !allowed {
					return nil, fmt.Errorf("field '%s' not allowed for filtering", field)
				}
				if len(fieldConfig.AllowedOperators) > 0 {
					found := false
					for _, op := range fieldConfig.AllowedOperators {
						if op == operator {
							found = true
							break
						}
					}
					if !found {
						return nil, fmt.Errorf("operator '%s' not allowed for field '%s'", operator, field)
					}
				}
			}

			if !allowedOperators[operator] {
				return nil, fmt.Errorf("operator '%s' not allowed at '%s'", operator, filterPart)
			}

			condValue, err := parseFilterValue(valueRaw)
			if err != nil {
				return nil, err
			}

			b.And(field, operator, condValue)
		}
	}

	return b, nil
}

func parseFilterValue(value string) (any, error) {
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		trimmed := strings.TrimPrefix(value, "'")
		trimmed = strings.TrimSuffix(trimmed, "'")
		return strings.ReplaceAll(trimmed, "''", "'"), nil
	}

	if value == "true" || value == "false" {
		return value == "true", nil
	}

	if value == "null" {
		return nil, nil
	}

	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, nil
	}

	return nil, fmt.Errorf("invalid value %s", value)
}

func splitFilter(filter string) ([]string, error) {
	firstSpace := strings.Index(filter, " ")
	if firstSpace == -1 {
		return nil, fmt.Errorf("filter syntax error (missing field separator) at '%s'", filter)
	}

	field := filter[:firstSpace]
	remaining := strings.TrimSpace(filter[firstSpace:])

	secondSpace := strings.Index(remaining, " ")
	if secondSpace == -1 {
		return nil, fmt.Errorf("filter syntax error (missing operator separator) at '%s'", filter)
	}

	operator := remaining[:secondSpace]
	value := strings.TrimSpace(remaining[secondSpace:])

	return []string{field, operator, value}, nil
}

func splitByDelimiterOutsideQuotesAndParens(input, delimiter string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	parenLevel := 0

	for i := 0; i < len(input); i++ {
		char := input[i]
		if char == '\'' {
			inQuotes = !inQuotes
		} else if char == '(' && !inQuotes {
			parenLevel++
		} else if char == ')' && !inQuotes {
			parenLevel--
		}

		if !inQuotes && parenLevel == 0 && strings.HasPrefix(input[i:], delimiter) {
			if current.Len() > 0 {
				result = append(result, strings.TrimSpace(current.String()))
			}
			current.Reset()
			i += len(delimiter) - 1
			continue
		}

		current.WriteByte(char)
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}
	return result
}
