package querybuilder

import (
	"fmt"
	"strings"
)

type SQLFragments struct {
	Where   string
	Args    []any
	OrderBy string
	Limit   string
	Offset  string
}

func (b *Builder) ToSQL(startIdx int) SQLFragments {
	if b == nil {
		return SQLFragments{Where: "1=1"}
	}

	where := "1=1"
	args := []any{}
	currentIdx := startIdx

	// 1. AND Conditions
	for _, cond := range b.AndConditions {
		sql, condArgs := conditionToSQL(cond, &currentIdx)
		if sql != "" {
			where += " AND " + sql
			args = append(args, condArgs...)
		}
	}

	// 2. OR Groups
	for _, group := range b.OrGroups {
		var orParts []string
		for _, cond := range group {
			sql, condArgs := conditionToSQL(cond, &currentIdx)
			if sql != "" {
				orParts = append(orParts, sql)
				args = append(args, condArgs...)
			}
		}
		if len(orParts) > 0 {
			where += fmt.Sprintf(" AND (%s)", strings.Join(orParts, " OR "))
		}
	}

	// 3. Order
	orderBy := ""
	if len(b.Orders) > 0 {
		var parts []string
		for _, o := range b.Orders {
			dir := "ASC"
			if strings.ToLower(o.Dir) == "desc" {
				dir = "DESC"
			}
			parts = append(parts, fmt.Sprintf("%s %s", o.Field, dir))
		}
		orderBy = " ORDER BY " + strings.Join(parts, ", ")
	}

	// 4. Limit
	limit := ""
	if b.LimitValue != nil {
		l := *b.LimitValue
		if l < 0 {
			l = 0
		}
		limit = fmt.Sprintf(" LIMIT %d", l)
	}

	// 5. Offset
	offset := ""
	if b.OffsetValue != nil {
		off := *b.OffsetValue
		if off < 0 {
			off = 0
		}
		offset = fmt.Sprintf(" OFFSET %d", off)
	}

	return SQLFragments{
		Where:   where,
		Args:    args,
		OrderBy: orderBy,
		Limit:   limit,
		Offset:  offset,
	}
}

func conditionToSQL(cond Condition, idx *int) (string, []any) {
	switch cond.Operator {
	case "eq":
		if cond.Value == nil {
			return fmt.Sprintf("%s IS NULL", cond.Field), nil
		}
		sql := fmt.Sprintf("%s = $%d", cond.Field, *idx)
		*idx++
		return sql, []any{cond.Value}

	case "ne":
		if cond.Value == nil {
			return fmt.Sprintf("%s IS NOT NULL", cond.Field), nil
		}
		sql := fmt.Sprintf("%s != $%d", cond.Field, *idx)
		*idx++
		return sql, []any{cond.Value}

	case "lt":
		sql := fmt.Sprintf("%s < $%d", cond.Field, *idx)
		*idx++
		return sql, []any{cond.Value}

	case "lte":
		sql := fmt.Sprintf("%s <= $%d", cond.Field, *idx)
		*idx++
		return sql, []any{cond.Value}

	case "gt":
		sql := fmt.Sprintf("%s > $%d", cond.Field, *idx)
		*idx++
		return sql, []any{cond.Value}

	case "gte":
		sql := fmt.Sprintf("%s >= $%d", cond.Field, *idx)
		*idx++
		return sql, []any{cond.Value}

	case "like":
		sql := fmt.Sprintf("upper(%s) LIKE upper($%d)", cond.Field, *idx)
		val := fmt.Sprintf("%%%v%%", cond.Value)
		*idx++
		return sql, []any{val}

	case "in":
		if slice, ok := cond.Value.([]any); ok {
			if len(slice) == 0 {
				return "FALSE", nil
			}

			var nonNulls []any
			hasNull := false
			for _, v := range slice {
				if v == nil {
					hasNull = true
				} else {
					nonNulls = append(nonNulls, v)
				}
			}

			var parts []string
			var args []any

			if len(nonNulls) > 0 {
				placeholders := make([]string, len(nonNulls))
				for i := range nonNulls {
					placeholders[i] = fmt.Sprintf("$%d", *idx)
					*idx++
				}
				parts = append(parts, fmt.Sprintf("%s IN (%s)", cond.Field, strings.Join(placeholders, ", ")))
				args = append(args, nonNulls...)
			}

			if hasNull {
				parts = append(parts, fmt.Sprintf("%s IS NULL", cond.Field))
			}

			if len(parts) == 1 {
				return parts[0], args
			}
			return "(" + strings.Join(parts, " OR ") + ")", args
		}
		// If not a slice, treat as eq
		sql := fmt.Sprintf("%s = $%d", cond.Field, *idx)
		*idx++
		return sql, []any{cond.Value}
	}

	return "", nil
}
