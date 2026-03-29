package utils

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

func QueryOpts() *QueryOptsBuilder {
	return &QueryOptsBuilder{
		AndConditions: make([]Condition, 0),
		OrGroups:      make([][]Condition, 0),
		Orders:        make([]Order, 0),
		LimitValue:    nil,
		OffsetValue:   nil,
	}
}

type Condition struct {
	Field    string
	Operator string
	Value    any
}

type Order struct {
	field string
	dir   string
}

type QueryOptsBuilder struct {
	AndConditions []Condition
	OrGroups      [][]Condition
	Orders        []Order
	LimitValue    *int
	OffsetValue   *int
}

func (qo *QueryOptsBuilder) And(field, operator string, value any) *QueryOptsBuilder {
	qo.AndConditions = append(qo.AndConditions, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return qo
}

type OrBuilder struct {
	conditions []Condition
	qo         *QueryOptsBuilder
}

func (ob *OrBuilder) Or(field string, operator string, value any) *OrBuilder {
	ob.conditions = append(ob.conditions, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return ob
}
func (ob *OrBuilder) EndOr() *QueryOptsBuilder {
	ob.qo.OrGroups = append(ob.qo.OrGroups, ob.conditions)
	return ob.qo
}

func (qo *QueryOptsBuilder) InitOr() *OrBuilder {
	return &OrBuilder{
		qo: qo,
	}
}

func (qo *QueryOptsBuilder) OrderBy(field, direction string) *QueryOptsBuilder {
	qo.Orders = append(qo.Orders, Order{
		field: field,
		dir:   direction,
	})
	return qo
}

func (qo *QueryOptsBuilder) Limit(limit int) *QueryOptsBuilder {
	qo.LimitValue = &limit
	return qo
}

func (qo *QueryOptsBuilder) Offset(offset int) *QueryOptsBuilder {
	qo.OffsetValue = &offset
	return qo
}

func QueryOptsToSquirrel(query squirrel.SelectBuilder, qo *QueryOptsBuilder) squirrel.SelectBuilder {
	for _, andCondition := range qo.AndConditions {
		query = query.Where(conditionToSquirrel(andCondition))
	}

	for _, orGroup := range qo.OrGroups {
		orSqlizers := squirrel.Or{}
		for _, condition := range orGroup {
			orSqlizers = append(orSqlizers, conditionToSquirrel(condition))
		}

		query = query.Where(orSqlizers)
	}

	if qo.LimitValue != nil {
		query = query.Limit(uint64(*qo.LimitValue))
	}

	if qo.OffsetValue != nil {
		query = query.Offset(uint64(*qo.OffsetValue))
	}

	for _, order := range qo.Orders {
		if order.dir == "asc" {
			query = query.OrderBy(order.field + " ASC")
		} else {
			query = query.OrderBy(order.field + " DESC")
		}
	}

	return query
}

func DeleteOptsToSquirrel(query squirrel.DeleteBuilder, qo *QueryOptsBuilder) squirrel.DeleteBuilder {
	for _, andCondition := range qo.AndConditions {
		query = query.Where(conditionToSquirrel(andCondition))
	}

	for _, orGroup := range qo.OrGroups {
		orSqlizers := squirrel.Or{}
		for _, condition := range orGroup {
			orSqlizers = append(orSqlizers, conditionToSquirrel(condition))
		}
		query = query.Where(orSqlizers)
	}

	if qo.LimitValue != nil {
		query = query.Limit(uint64(*qo.LimitValue))
	}

	for _, order := range qo.Orders {
		if order.dir == "asc" {
			query = query.OrderBy(order.field + " ASC")
		} else {
			query = query.OrderBy(order.field + " DESC")
		}
	}

	return query
}

func UpdateOptsToSquirrel(query squirrel.UpdateBuilder, qo *QueryOptsBuilder) squirrel.UpdateBuilder {
	for _, andCondition := range qo.AndConditions {
		query = query.Where(conditionToSquirrel(andCondition))
	}

	for _, orGroup := range qo.OrGroups {
		orSqlizers := squirrel.Or{}
		for _, condition := range orGroup {
			orSqlizers = append(orSqlizers, conditionToSquirrel(condition))
		}
		query = query.Where(orSqlizers)
	}

	if qo.LimitValue != nil {
		query = query.Limit(uint64(*qo.LimitValue))
	}

	for _, order := range qo.Orders {
		if order.dir == "asc" {
			query = query.OrderBy(order.field + " ASC")
		} else {
			query = query.OrderBy(order.field + " DESC")
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

func ForCount(qo *QueryOptsBuilder) *QueryOptsBuilder {
	return &QueryOptsBuilder{
		AndConditions: qo.AndConditions,
		OrGroups:      qo.OrGroups,
		Orders:        nil,
		LimitValue:    nil,
		OffsetValue:   nil,
	}
}
