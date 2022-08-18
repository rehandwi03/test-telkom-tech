package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type DBENGINE string

const (
	DATABASE_ENGINE_POSTGRESQL DBENGINE = "postgresql"
	DATABASE_ENGINE_MYSQL               = "mysql"
)

type Queryer interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Rebind(query string) string
	QueryRow(query string, args ...interface{}) *sql.Row
	sqlx.Queryer
	sqlx.Execer
}

type Where struct {
	Or  []squirrel.Or
	And []squirrel.And
	Eq  []squirrel.Eq
}

type QueryBuilderCriteria struct {
	Where *Where

	// this is for ordering, use key for field and value for ordering.
	// example: map["id"] = "ASC"
	Order  map[string]string
	Limit  *uint64
	Offset *uint64
	Join   struct {
		LeftJoin  []string
		RightJoin []string
		InnerJoin []string
	}
	Group  []string
	Select []string
}

func (q *QueryBuilderCriteria) GenerateSquirrelQuery(
	tableName string,
	dbEngine DBENGINE,
) (res squirrel.SelectBuilder, err error) {
	if len(q.Select) > 0 {
		res = squirrel.Select(q.Select...).From(tableName)
	} else {
		res = squirrel.Select("*").From(tableName)
	}

	switch dbEngine {
	case DATABASE_ENGINE_POSTGRESQL:
		res = res.PlaceholderFormat(squirrel.Dollar)
	case DATABASE_ENGINE_MYSQL:
		res = res.PlaceholderFormat(squirrel.Question)
	default:
		return res, errors.New("database engine invalid")
	}

	if q.Order != nil {
		for order, sort := range q.Order {
			orderBy := fmt.Sprintf("%s %s", order, sort)
			res = res.OrderBy(orderBy)
		}
	}

	if len(q.Join.LeftJoin) > 0 {
		for _, leftJoin := range q.Join.LeftJoin {
			res = res.LeftJoin(leftJoin)
		}
	}

	if len(q.Join.RightJoin) > 0 {
		for _, rightJoin := range q.Join.RightJoin {
			res = res.RightJoin(rightJoin)
		}
	}

	if len(q.Join.InnerJoin) > 0 {
		for _, innerJoin := range q.Join.InnerJoin {
			res = res.InnerJoin(innerJoin)
		}
	}

	if q.Limit != nil {
		res = res.Limit(*q.Limit)
	}

	if q.Offset != nil {
		res = res.Offset(*q.Offset)
	}

	if q.Group != nil {
		res = res.GroupBy(q.Group...)
	}

	if q.Where != nil {
		if len(q.Where.Eq) > 0 {
			for _, eq := range q.Where.Eq {
				res = res.Where(eq)
			}
		}

		if len(q.Where.And) > 0 {
			for _, and := range q.Where.And {
				res = res.Where(and)
			}
		}

		if len(q.Where.Or) > 0 {
			for _, or := range q.Where.Or {
				res = res.Where(or)
			}
		}
	}

	return res, nil
}

func (q *QueryBuilderCriteria) GenerateSquirrelQueryCountData(
	tableName string,
	dbEngine DBENGINE,
) (res squirrel.SelectBuilder, err error) {
	if len(q.Select) > 0 {
		res = squirrel.Select(q.Select...).From(tableName)
	} else {
		res = squirrel.Select("*").From(tableName)
	}

	switch dbEngine {
	case DATABASE_ENGINE_POSTGRESQL:
		res = res.PlaceholderFormat(squirrel.Dollar)
	case DATABASE_ENGINE_MYSQL:
		res = res.PlaceholderFormat(squirrel.Question)
	default:
		return res, errors.New("db engine invalid")
	}

	if q.Order != nil {
		for order, sort := range q.Order {
			orderBy := fmt.Sprintf("%s %s", order, sort)
			res = res.OrderBy(orderBy)
		}
	}

	if len(q.Join.LeftJoin) > 0 {
		for _, leftJoin := range q.Join.LeftJoin {
			res = res.LeftJoin(leftJoin)
		}
	}

	if len(q.Join.RightJoin) > 0 {
		for _, rightJoin := range q.Join.RightJoin {
			res = res.RightJoin(rightJoin)
		}
	}

	if len(q.Join.InnerJoin) > 0 {
		for _, innerJoin := range q.Join.InnerJoin {
			res = res.InnerJoin(innerJoin)
		}
	}

	if q.Where != nil {
		if len(q.Where.Eq) > 0 {
			for _, eq := range q.Where.Eq {
				res = res.Where(eq)
			}
		}

		if len(q.Where.And) > 0 {
			for _, and := range q.Where.And {
				res = res.Where(and)
			}
		}

		if len(q.Where.Or) > 0 {
			for _, or := range q.Where.Or {
				res = res.Where(or)
			}
		}
	}

	return res, nil
}
