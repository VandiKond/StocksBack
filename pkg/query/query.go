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
	OR:  "or",
	AND: "and",
}

// the Sign id
type Sign int

// The Sign ids uSign iota
const (
	EQUAL Sign = iota
	MORE
	LESS
)

// The map for string values of Signs
var StringSign = map[Sign]string{
	EQUAL: "=",
	MORE:  ">",
	LESS:  "<",
}

// The settings of query
//
// Separator: Separator id
// Type: UserField id
// Sign: Sign id
// Not: need to use not
// Y: The compare value
type QuerySetting struct {
	Separator Separator `json:"separator"`
	Type      UserField `json:"type"`
	Sign      Sign      `json:"Sign"`
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

		// Switching by Sign
		switch q.Sign {
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

		// Switching by Sign
		switch q.Sign {
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

		// Switching by Sign
		switch q.Sign {
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

		// Switching by Sign
		switch q.Sign {
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
func (query Query) Sort(users []user_cfg.User, num int) ([]user_cfg.User, error) {
	// If num is less than zero it sets to max value
	if num < 0 {
		num = len(users) - 1
	}

	// Setting the result
	var res []user_cfg.User

	for _, u := range users {
		// Setting do append and current separator
		var do_append bool
		var cur_separator Separator = OR

		for _, qr := range query {
			// Going through separators
			switch qr.Separator {

			case NOT_SEPARATOR:
				// Setting is true
				var is_true bool

				switch qr.Type {

				case ID:
					// Running id
					is_true = qr.Run(u.Id)

				case NAME:
					// Running name
					is_true = qr.Run(u.Name)

				case PASSWORD:
					// Running password
					is_true = qr.Run(u.Password)

				case SOLID_BALANCE:
					// Running solid balance
					is_true = qr.Run(u.SolidBalance)

				case STOCK_BALANCE:
					// Running stock balance
					is_true = qr.Run(u.StockBalance)

				case IS_BLOCKED:
					// Running is blocked
					is_true = qr.Run(u.IsBlocked)

				case LAST_FARMING:
					// Running last farming
					is_true = qr.Run(u.LastFarming)
				case CREATED_AT:
					// Running created at
					is_true = qr.Run(u.CreatedAt)

				default:
					// Returning error
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid type")
				}

				// Editing do append
				switch cur_separator {

				case OR:
					do_append = do_append || is_true

				case AND:
					do_append = do_append && is_true

				default:
					// Returning error
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
				}

			case OR, AND:

				// Checking the order
				if cur_separator != NOT_SEPARATOR {
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
				}

			default:
				// Returning error
				return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
			}

			// Setting current separator
			cur_separator = qr.Separator
		}

		// Checking the order
		if cur_separator != NOT_SEPARATOR {
			return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
		}

		// Appending user
		if do_append {
			res = append(res, u)
		}

		// Checking the limit
		if len(res) > num {
			break
		}
	}

	return res, nil
}

// Gets the string Sign
func SignToString(Sign Sign, not bool) string {
	if Sign == EQUAL && not {
		// !=
		return "!" + StringSign[Sign]
	} else if Sign == EQUAL && !not {
		// ==
		return StringSign[Sign] + StringSign[Sign]
	} else if Sign == MORE && not {
		// <=
		return StringSign[LESS] + "="
	} else if Sign == LESS && not {
		// >=
		return StringSign[MORE] + "="
	} else {
		// > and <
		return StringSign[Sign]
	}
}

// Creates a string query
func (query Query) String() string {
	// Setting result as a string
	var res string

	for _, qr := range query {
		// Separator switch
		switch qr.Separator {

		case NOT_SEPARATOR:
			// Adding query setting expression
			res += fmt.Sprintf("%s %s %v", StringUserField[qr.Type], SignToString(qr.Sign, qr.Not), qr.Y)

		case OR, AND:
			// Adding separator
			res += " " + StringSeparator[qr.Separator] + " "
		}
	}

	return res
}

func (query Query) PrepareString() (string, []any) {
	// Setting result as a string
	var resStr string
	var resSlice []any

	for i, qr := range query {
		// Separator switch
		switch qr.Separator {

		case NOT_SEPARATOR:
			// Adding query setting expression
			resStr += fmt.Sprintf("%s %s $%d", StringUserField[qr.Type], SignToString(qr.Sign, qr.Not), i+1)
			resSlice = append(resSlice, qr.Y)
		case OR, AND:
			// Adding separator
			resStr += " " + StringSeparator[qr.Separator] + " "
		}
	}

	return resStr, resSlice
}
