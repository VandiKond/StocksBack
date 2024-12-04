package query

import (
	"fmt"
	"time"

	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	InvalidQuery = "invalid query"
)

// The user field id
type UserField int

// The user fields ids using iota
const (
	ID UserField = iota
	NAME
	PASSWORD
	SOLID_BALANCE
	STOCK_BALANCE
	IS_BLOCKED
	LAST_FARMING
	CREATED_AT
)

// The map for string values of fields
var StringUserField = map[UserField]string{
	ID:            "id",
	NAME:          "name",
	PASSWORD:      "password",
	SOLID_BALANCE: "solid_balance",
	STOCK_BALANCE: "stock_balance",
	IS_BLOCKED:    "is_blocked",
	LAST_FARMING:  "last_farming",
	CREATED_AT:    "created_at",
}

// the separator id
type Separator int

// The separator ids using iota
const (
	NOT_SEPARATOR Separator = iota
	OR
	AND
)

// The map for string values of separators
var StringSeparator = map[Separator]string{
	OR:  "OR",
	AND: "AND",
}

// the sing id
type Sing int

// The sing ids using iota
const (
	EQUAL Sing = iota
	MORE
	LESS
)

// The map for string values of sings
var StringSing = map[Sing]string{
	EQUAL: "=",
	MORE:  ">",
	LESS:  "<",
}

// The settings of query
//
// Separator: Separator id
// Type: UserField id
// Sing: Sing id
// Not: need to use not
// Y: The compare value
type QuerySetting struct {
	Separator Separator `json:"separator"`
	Type      UserField `json:"type"`
	Sing      Sing      `json:"sing"`
	Not       bool      `json:"not"`
	Y         any       `json:"y"`
}

// The query setting slice, that represents a query
type Query []QuerySetting

// Runs the query
func (q QuerySetting) Run(X any) bool {
	// Sets the result
	var res bool

	// Switching by type
	switch q.Type {
	// Case of uint64 values
	case ID, SOLID_BALANCE, STOCK_BALANCE:
		// Converting x to uint64
		uintX, ok := X.(uint64)
		if !ok {
			return false
		}
		// Converting y to uint64
		uintY, ok := q.Y.(uint64)
		if !ok {
			return false
		}

		// Switching by sing
		switch q.Sing {
		// Checking equal
		case EQUAL:
			res = uintX == uintY
		// Checking more
		case MORE:
			res = uintX > uintY
		// Checking less
		case LESS:
			res = uintX < uintY
		// In default
		default:
			return false
		}
	// Case of string values
	case NAME, PASSWORD:
		// Converting x to string
		strX, ok := X.(string)
		if !ok {
			return false
		}
		// Converting y to string
		strY, ok := q.Y.(string)
		if !ok {
			return false
		}

		// Switching by sing
		switch q.Sing {
		// Checking equal
		case EQUAL:
			res = strX == strY
		// Checking more
		case MORE:
			res = strX > strY
		// Checking less
		case LESS:
			res = strX < strY
		// In default
		default:
			return false
		}
	// Case of boolean values
	case IS_BLOCKED:
		// Converting x to boolean
		boolX, ok := X.(bool)
		if !ok {
			return false
		}

		// Converting y to boolean
		boolY, ok := q.Y.(bool)
		if !ok {
			return false
		}

		// Switching by sing
		switch q.Sing {
		// Checking equal
		case EQUAL:
			res = boolX == boolY
			// Checking more
		case MORE:
			res = boolX && !boolY
			// Checking less
		case LESS:
			res = !boolX && boolY
			// In default
		default:
			return false
		}
	// Case of time.Time values
	case CREATED_AT, LAST_FARMING:
		// Converting x to time.Time
		timeX, ok := X.(time.Time)
		if !ok {
			return false
		}

		// Converting y to time.Time
		timeY, ok := q.Y.(time.Time)
		if !ok {
			return false
		}

		// Switching by sing
		switch q.Sing {
		// Checking equal
		case EQUAL:
			res = timeX.Equal(timeY)
			// Checking more
		case MORE:
			res = timeX.After(timeY)
			// Checking less
		case LESS:
			res = timeX.Before(timeY)
			// In default
		default:
			return false
		}
	}

	// Checking not
	/*
		Note:
		to create != use Equal + Not
		to create >= use Less + Not
		to create <= use More + Not
	*/
	if q.Not {
		res = !res
	}

	return res
}

// Sorting the users by query
// TODO : Write comments from this function
func (query Query) Sort(users []user_cfg.User, num int) ([]user_cfg.User, error) {
	if num == -1 {
		num = len(users) - 1
	}
	var res []user_cfg.User

	for _, u := range users {
		var do_append bool
		var cur_separator Separator = OR

		for _, qr := range query {
			switch qr.Separator {
			case NOT_SEPARATOR:
				var is_true bool
				switch qr.Type {
				case ID:
					is_true = qr.Run(u.Id)
				case NAME:
					is_true = qr.Run(u.Name)
				case PASSWORD:
					is_true = qr.Run(u.Password)
				case SOLID_BALANCE:
					is_true = qr.Run(u.SolidBalance)
				case STOCK_BALANCE:
					is_true = qr.Run(u.StockBalance)
				case IS_BLOCKED:
					is_true = qr.Run(u.IsBlocked)
				case LAST_FARMING:
					is_true = qr.Run(u.LastFarming)
				case CREATED_AT:
					is_true = qr.Run(u.CreatedAt)
				default:
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid type")
				}

				switch cur_separator {
				case OR:
					do_append = do_append || is_true
				case AND:
					do_append = do_append && is_true
				default:
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
				}

			case OR, AND:
				if cur_separator != NOT_SEPARATOR {
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
				}
			}

			cur_separator = qr.Separator
		}

		if cur_separator != NOT_SEPARATOR {
			return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
		}

		if do_append {
			res = append(res, u)
		}

		if len(res) > num {
			break
		}
	}
	return res, nil
}

func SignToString(sing Sing, not bool) string {
	if sing == EQUAL && not {
		return "!" + StringSing[sing]
	} else if sing == EQUAL && !not {
		return StringSing[sing] + StringSing[sing]
	} else if sing == MORE && not {
		return StringSing[LESS] + "="
	} else if sing == LESS && not {
		return StringSing[MORE] + "="
	} else {
		return StringSing[sing]
	}
}

func (query Query) String() string {
	var res string

	for _, qr := range query {
		switch qr.Separator {
		case NOT_SEPARATOR:
			res += fmt.Sprintf("%s %s %v", StringUserField[qr.Type], SignToString(qr.Sing, qr.Not), qr.Y)

		case OR, AND:
			res += " " + StringSeparator[qr.Separator] + " "
		}
	}
	return res
}
