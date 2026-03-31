package pagination_util

import (
	"cafe/domain"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var allowedOperators = map[string]bool{
	"=":           true,
	"!=":          true,
	">":           true,
	">=":          true,
	"<":           true,
	"<=":          true,
	"LIKE":        true,
	"ILIKE":       true,
	"NOT LIKE":    true,
	"IN":          true,
	"NOT IN":      true,
	"BETWEEN":     true,
	"NOT BETWEEN": true,
	"IS NULL":     true,
	"IS NOT NULL": true,
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Filter(filters []domain.Criterion, whitelist map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(filters) == 0 {
			return db
		}

		workingFilters := filters
		if len(workingFilters) > 50 {
			workingFilters = workingFilters[:50]
		}

		subQuery := db.Session(&gorm.Session{NewDB: true})
		hasConditions := false

		isFirst := true
		for _, f := range workingFilters {
			column, ok := whitelist[f.Field]
			if !ok {
				continue
			}

			operator := strings.ToUpper(strings.TrimSpace(f.Operator))
			if !allowedOperators[operator] {
				continue
			}

			condition, args := buildCondition(db, column, operator, f.Value)
			if condition == "" {
				continue
			}

			if isFirst {
				subQuery = subQuery.Where(condition, args...)
				isFirst = false
				hasConditions = true
			} else {
				switch f.Logic {
				case domain.Or:
					subQuery = subQuery.Or(condition, args...)
				default:
					subQuery = subQuery.Where(condition, args...)
				}
			}
		}

		if !hasConditions {
			return db
		}

		return db.Where(subQuery)
	}
}

func Sort(sortField string, sortOrder domain.SortOrder, whitelist map[string]string, anchorField string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		column, ok := whitelist[sortField]
		anchorColumn, okAnchor := whitelist[anchorField]

		if sortOrder != domain.Asc && sortOrder != domain.Desc {
			sortOrder = domain.Asc
		}

		// Case 1: Neither sortField nor anchorField is valid
		if !ok && !okAnchor {
			return db
		}

		// Case 2: Only sortField is valid
		if ok && !okAnchor {
			return db.Order(fmt.Sprintf("%s %s", column, sortOrder))
		}

		// Case 3: Only anchorField is valid (or sortField omitted)
		if !ok && okAnchor {
			return db.Order(fmt.Sprintf("%s ASC", anchorColumn))
		}

		// Case 4: Both are valid - avoid duplication if they are the same column
		if column == anchorColumn {
			return db.Order(fmt.Sprintf("%s %s", column, sortOrder))
		}

		return db.Order(fmt.Sprintf("%s %s", column, sortOrder)).Order(fmt.Sprintf("%s ASC", anchorColumn))
	}
}

func buildCondition(db *gorm.DB, column, operator string, value interface{}) (string, []interface{}) {
	if operator == "ILIKE" && db.Dialector.Name() != "postgres" {
		operator = "LIKE"
	}

	switch operator {
	case "LIKE", "ILIKE", "NOT LIKE":
		strVal, ok := value.(string)
		if !ok || strings.Trim(strVal, "%") == "" {
			return "", nil
		}
		return fmt.Sprintf("%s %s ?", column, operator), []interface{}{value}

	case "IN", "NOT IN":
		return buildInCondition(db, column, operator, value)

	case "BETWEEN", "NOT BETWEEN":
		return buildBetweenCondition(column, operator, value)

	case "IS NULL", "IS NOT NULL":
		return fmt.Sprintf("%s %s", column, operator), nil

	default:
		switch value.(type) {
		case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, time.Time:
			return fmt.Sprintf("%s %s ?", column, operator), []interface{}{value}
		default:
			return "", nil
		}
	}
}

func buildInCondition(db *gorm.DB, column, operator string, value interface{}) (string, []interface{}) {
	var length int
	var truncatedValue interface{} = value

	switch v := value.(type) {
	case []interface{}:
		length = len(v)
		if length > 1000 {
			truncatedValue = v[:1000]
			db.Logger.Warn(db.Statement.Context, "PAGNATION_UTIL: IN clause truncated to 1000 elements for column: %s", column)
		}
	case []string:
		length = len(v)
		if length > 1000 {
			truncatedValue = v[:1000]
			db.Logger.Warn(db.Statement.Context, "PAGNATION_UTIL: IN clause truncated to 1000 elements for column: %s", column)
		}
	case []int:
		length = len(v)
		if length > 1000 {
			truncatedValue = v[:1000]
			db.Logger.Warn(db.Statement.Context, "PAGNATION_UTIL: IN clause truncated to 1000 elements for column: %s", column)
		}
	case []int64:
		length = len(v)
		if length > 1000 {
			truncatedValue = v[:1000]
			db.Logger.Warn(db.Statement.Context, "PAGNATION_UTIL: IN clause truncated to 1000 elements for column: %s", column)
		}
	default:
		return "", nil
	}

	if length == 0 {
		return "", nil
	}

	return fmt.Sprintf("%s %s (?)", column, operator), []interface{}{truncatedValue}
}

func buildBetweenCondition(column, operator string, value interface{}) (string, []interface{}) {
	var args []interface{}
	switch v := value.(type) {
	case []interface{}:
		if len(v) >= 2 {
			args = compareAndSwapInterface(v[0], v[1])
		}
	case []string:
		if len(v) >= 2 {
			vCopy := []string{v[0], v[1]}
			if vCopy[0] > vCopy[1] {
				vCopy[0], vCopy[1] = vCopy[1], vCopy[0]
			}
			args = []interface{}{vCopy[0], vCopy[1]}
		}
	case []int:
		if len(v) >= 2 {
			vCopy := []int{v[0], v[1]}
			if vCopy[0] > vCopy[1] {
				vCopy[0], vCopy[1] = vCopy[1], vCopy[0]
			}
			args = []interface{}{vCopy[0], vCopy[1]}
		}
	case []int64:
		if len(v) >= 2 {
			vCopy := []int64{v[0], v[1]}
			if vCopy[0] > vCopy[1] {
				vCopy[0], vCopy[1] = vCopy[1], vCopy[0]
			}
			args = []interface{}{vCopy[0], vCopy[1]}
		}
	}

	if len(args) < 2 {
		return "", nil
	}
	return fmt.Sprintf("%s %s ? AND ?", column, operator), args
}

func compareAndSwapInterface(a, b interface{}) []interface{} {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok && va > vb {
			return []interface{}{vb, va}
		}
	case int64:
		if vb, ok := b.(int64); ok && va > vb {
			return []interface{}{vb, va}
		}
	case float64:
		if vb, ok := b.(float64); ok && va > vb {
			return []interface{}{vb, va}
		}
	case string:
		if vb, ok := b.(string); ok && va > vb {
			return []interface{}{vb, va}
		}
	case time.Time:
		if vb, ok := b.(time.Time); ok && va.After(vb) {
			return []interface{}{vb, va}
		}
	}
	return []interface{}{a, b}
}
