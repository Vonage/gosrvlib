package sqlconn

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	sqlConditionIn    = "IN"
	sqlConditionNotIn = "NOT IN"
)

// BuildInClauseString prepares a SQL IN clause with the given list of string values.
func (c *SQLConn) BuildInClauseString(field string, values []string) string {
	return c.composeInClause(sqlConditionIn, field, c.formatStrings(values))
}

// BuildNotInClauseString prepares a SQL NOT IN clause with the given list of string values.
func (c *SQLConn) BuildNotInClauseString(field string, values []string) string {
	return c.composeInClause(sqlConditionNotIn, field, c.formatStrings(values))
}

// BuildInClauseInt prepares a SQL IN clause with the given list of integer values.
func (c *SQLConn) BuildInClauseInt(field string, values []int) string {
	return c.composeInClause(sqlConditionIn, field, formatInts(values))
}

// BuildNotInClauseInt prepares a SQL NOT IN clause with the given list of integer values.
func (c *SQLConn) BuildNotInClauseInt(field string, values []int) string {
	return c.composeInClause(sqlConditionNotIn, field, formatInts(values))
}

// BuildInClauseUint prepares a SQL IN clause with the given list of integer values.
func (c *SQLConn) BuildInClauseUint(field string, values []uint64) string {
	return c.composeInClause(sqlConditionIn, field, formatUints(values))
}

// BuildNotInClauseUint prepares a SQL NOT IN clause with the given list of integer values.
func (c *SQLConn) BuildNotInClauseUint(field string, values []uint64) string {
	return c.composeInClause(sqlConditionNotIn, field, formatUints(values))
}

func (c *SQLConn) composeInClause(condition string, field string, values []string) string {
	if len(values) == 0 {
		return ""
	}

	return fmt.Sprintf("%s %s (%s)", c.QuoteID(field), condition, strings.Join(values, ","))
}

func (c *SQLConn) formatStrings(values []string) []string {
	items := make([]string, len(values))

	for k, v := range values {
		items[k] = c.QuoteValue(v)
	}

	return items
}

func formatInts(values []int) []string {
	items := make([]string, len(values))

	for k, v := range values {
		items[k] = strconv.Itoa(v)
	}

	return items
}

func formatUints(values []uint64) []string {
	items := make([]string, len(values))

	for k, v := range values {
		items[k] = strconv.FormatUint(v, 10)
	}

	return items
}
