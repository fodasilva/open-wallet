package querybuilder

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

func ToSquirrel(query squirrel.SelectBuilder, b *Builder) squirrel.SelectBuilder {
	if b == nil {
		return query
	}

	for _, andCondition := range b.AndConditions {
		query = query.Where(conditionToSquirrel(andCondition))
	}

	for _, orGroup := range b.OrGroups {
		orSqlizers := squirrel.Or{}
		for _, condition := range orGroup {
			orSqlizers = append(orSqlizers, conditionToSquirrel(condition))
		}

		query = query.Where(orSqlizers)
	}

	if b.LimitValue != nil {
		query = query.Limit(uint64(*b.LimitValue))
	}

	if b.OffsetValue != nil {
		query = query.Offset(uint64(*b.OffsetValue))
	}

	for _, order := range b.Orders {
		if order.Dir == "asc" {
			query = query.OrderBy(order.Field + " ASC")
		} else {
			query = query.OrderBy(order.Field + " DESC")
		}
	}

	return query
}

func ToDeleteSquirrel(query squirrel.DeleteBuilder, b *Builder) squirrel.DeleteBuilder {
	if b == nil {
		return query
	}

	for _, andCondition := range b.AndConditions {
		query = query.Where(conditionToSquirrel(andCondition))
	}

	for _, orGroup := range b.OrGroups {
		orSqlizers := squirrel.Or{}
		for _, condition := range orGroup {
			orSqlizers = append(orSqlizers, conditionToSquirrel(condition))
		}
		query = query.Where(orSqlizers)
	}

	if b.LimitValue != nil {
		query = query.Limit(uint64(*b.LimitValue))
	}

	for _, order := range b.Orders {
		if order.Dir == "asc" {
			query = query.OrderBy(order.Field + " ASC")
		} else {
			query = query.OrderBy(order.Field + " DESC")
		}
	}

	return query
}

func ToUpdateSquirrel(query squirrel.UpdateBuilder, b *Builder) squirrel.UpdateBuilder {
	if b == nil {
		return query
	}

	for _, andCondition := range b.AndConditions {
		query = query.Where(conditionToSquirrel(andCondition))
	}

	for _, orGroup := range b.OrGroups {
		orSqlizers := squirrel.Or{}
		for _, condition := range orGroup {
			orSqlizers = append(orSqlizers, conditionToSquirrel(condition))
		}
		query = query.Where(orSqlizers)
	}

	if b.LimitValue != nil {
		query = query.Limit(uint64(*b.LimitValue))
	}

	for _, order := range b.Orders {
		if order.Dir == "asc" {
			query = query.OrderBy(order.Field + " ASC")
		} else {
			query = query.OrderBy(order.Field + " DESC")
		}
	}

	return query
}

func conditionToSquirrel(condition Condition) squirrel.Sqlizer {
	switch condition.Operator {
	case "eq":
		return squirrel.Eq{condition.Field: condition.Value}
	case "ne":
		return squirrel.NotEq{condition.Field: condition.Value}
	case "lt":
		return squirrel.Lt{condition.Field: condition.Value}
	case "lte":
		return squirrel.LtOrEq{condition.Field: condition.Value}
	case "gt":
		return squirrel.Gt{condition.Field: condition.Value}
	case "gte":
		return squirrel.GtOrEq{condition.Field: condition.Value}
	case "like":
		return squirrel.Like{fmt.Sprintf("upper(%s)", condition.Field): strings.ToUpper(fmt.Sprintf("%%%s%%", condition.Value))}
	default:
		return nil
	}
}
