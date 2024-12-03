package file_db

import (
	"time"
)

const (
	ID int = iota
	NAME
	PASSWORD
	SOLID_BALANCE
	STOCK_BALANCE
	IS_BLOCKED
	CREATED_AT
)

const (
	NOT_SEPARATOR int = iota
	OR
	AND
)

const (
	EQUAL int = iota
	MORE
	LESS
)

type QuerySetting struct {
	Separator int  `json:"separator"`
	Type      int  `json:"type"`
	Sing      int  `json:"sing"`
	Not       bool `json:"not"`
	Y         any  `json:"y"`
}

func (q QuerySetting) Run(X any) bool {
	var res bool
	switch q.Type {
	case ID, SOLID_BALANCE, STOCK_BALANCE:
		uintX, ok := X.(uint64)
		if !ok {
			return false
		}
		uintY, ok := q.Y.(uint64)
		if !ok {
			return false
		}
		switch q.Sing {
		case EQUAL:
			res = uintX == uintY
		case MORE:
			res = uintX > uintY
		case LESS:
			res = uintX < uintY
		default:
			return false
		}
	case NAME, PASSWORD:
		strX, ok := X.(string)
		if !ok {
			return false
		}
		strY, ok := q.Y.(string)
		if !ok {
			return false
		}
		switch q.Sing {
		case EQUAL:
			res = strX == strY
		case MORE:
			res = strX > strY
		case LESS:
			res = strX < strY
		default:
			return false
		}
	case IS_BLOCKED:
		boolX, ok := X.(bool)
		if !ok {
			return false
		}
		boolY, ok := q.Y.(bool)
		if !ok {
			return false
		}
		switch q.Sing {
		case EQUAL:
			res = boolX == boolY
		case MORE:
			res = boolX && !boolY
		case LESS:
			res = !boolX && boolY
		default:
			return false
		}
	case CREATED_AT:
		timeX, ok := X.(time.Time)
		if !ok {
			return false
		}
		timeY, ok := q.Y.(time.Time)
		if !ok {
			return false
		}
		switch q.Sing {
		case EQUAL:
			res = timeX.Equal(timeY)
		case MORE:
			res = timeX.After(timeY)
		case LESS:
			res = timeX.Before(timeY)
		default:
			return false
		}
	}
	if q.Not {
		res = !res
	}
	return res
}
