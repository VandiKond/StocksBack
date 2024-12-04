package query

import (
	"fmt"
	"time"

	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/vanerrors"
)

const (
	InvalidQuery = "invalid query" // invalid query
)

type UserPart int

const (
	ID UserPart = iota
	NAME
	PASSWORD
	SOLID_BALANCE
	STOCK_BALANCE
	IS_BLOCKED
	LAST_FARMING
	CREATED_AT
)

var StringUserPart = map[UserPart]string{
	ID:            "id",
	NAME:          "name",
	PASSWORD:      "password",
	SOLID_BALANCE: "solid_balance",
	STOCK_BALANCE: "stock_balance",
	IS_BLOCKED:    "is_blocked",
	LAST_FARMING:  "last_farming",
	CREATED_AT:    "created_at",
}

type Separator int

const (
	NOT_SEPARATOR Separator = iota
	OR
	AND
)

var StringSeparator = map[Separator]string{
	OR:  "OR",
	AND: "AND",
}

type Sing int

const (
	EQUAL Sing = iota
	MORE
	LESS
)

var StringSing = map[Sing]string{
	EQUAL: "=",
	MORE:  ">",
	LESS:  "<",
}

type QuerySetting struct {
	Separator Separator `json:"separator"`
	Type      UserPart  `json:"type"`
	Sing      Sing      `json:"sing"`
	Not       bool      `json:"not"`
	Y         any       `json:"y"`
}

type Query []QuerySetting

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
	case CREATED_AT, LAST_FARMING:
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
			res += fmt.Sprintf("%s %s %v", StringUserPart[qr.Type], SignToString(qr.Sing, qr.Not), qr.Y)

		case OR, AND:
			res += " " + StringSeparator[qr.Separator] + " "
		}
	}
	return res
}
