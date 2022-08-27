package persistence

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"interview-telkom-6/entity"
	"log"
)

type cartRepository struct {
	Conn      Queryer
	TableName string
}

type CartRepository interface {
	WithTx(conn *sqlx.Tx) CartRepository
	Get(ctx context.Context, builder *QueryBuilderCriteria) (
		res entity.Cart, err error,
	)
	Find(ctx context.Context, builder *QueryBuilderCriteria) (
		res []entity.Cart, err error,
	)
	Store(ctx context.Context, data *entity.Cart) (res entity.Cart, err error)
	Update(ctx context.Context, data *entity.Cart) (res entity.Cart, err error)
	Delete(ctx context.Context, data *entity.Cart) (err error)
	Count(ctx context.Context, builder *QueryBuilderCriteria) (totalRow int64, err error)
}

func NewCartRepository(conn *sqlx.DB) CartRepository {
	return &cartRepository{Conn: conn, TableName: "carts"}
}

func (r cartRepository) WithTx(conn *sqlx.Tx) CartRepository {
	if conn == nil {
		log.Println("transaction database not found")
		return &r
	}

	return &cartRepository{Conn: conn}
}

func (r cartRepository) Get(ctx context.Context, builder *QueryBuilderCriteria) (
	res entity.Cart, err error,
) {
	sq, err := builder.GenerateSquirrelQuery(r.TableName, DATABASE_ENGINE_POSTGRESQL)
	if err != nil {
		log.Println(err)
		return res, err
	}

	query, args, err := sq.ToSql()
	if err != nil {
		log.Println(err)
		return res, err
	}

	log.Println(query)
	log.Println(args)

	err = r.Conn.Get(&res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (r cartRepository) Find(ctx context.Context, builder *QueryBuilderCriteria) (
	res []entity.Cart, err error,
) {
	sq, err := builder.GenerateSquirrelQuery(r.TableName, DATABASE_ENGINE_POSTGRESQL)
	if err != nil {
		log.Println(err)
		return res, err
	}

	query, args, err := sq.ToSql()
	if err != nil {
		log.Println(err)
		return res, err
	}

	log.Println(query)
	log.Println(args)

	err = r.Conn.Select(&res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (r cartRepository) Store(ctx context.Context, data *entity.Cart) (res entity.Cart, err error) {
	data.GenerateUUID()
	query := fmt.Sprintf(
		"INSERT INTO carts (id, full_name) VALUES (:id, :full_name)",
	)
	log.Println(query)
	_, err = r.Conn.NamedExecContext(ctx, query, data)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return *data, err
}

func (r cartRepository) Update(ctx context.Context, data *entity.Cart) (res entity.Cart, err error) {
	query := fmt.Sprintf(
		"UPDATE carts SET full_name=:full_name WHERE id=:id",
	)
	log.Println(query)
	_, err = r.Conn.NamedExecContext(ctx, query, data)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return *data, nil
}

func (r cartRepository) Delete(ctx context.Context, data *entity.Cart) (err error) {
	query := fmt.Sprintf("DELETE FROM carts WHERE id = $1")
	log.Println(query)
	_, err = r.Conn.Exec(query, data.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r cartRepository) Count(ctx context.Context, builder *QueryBuilderCriteria) (totalRow int64, err error) {
	sq, err := builder.GenerateSquirrelQueryCountData(r.TableName, DATABASE_ENGINE_POSTGRESQL)
	if err != nil {
		log.Println(err)
		return totalRow, err
	}

	queryFrom, args, err := sq.ToSql()
	if err != nil {
		log.Println(err)
		return totalRow, err
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS totalRow", queryFrom)

	log.Println(query)
	log.Println(args)

	rows, err := r.Conn.Queryx(query, args...)
	if err != nil {
		log.Println(err)
		return totalRow, err
	}

	for rows.Next() {
		err = rows.Scan(&totalRow)
		if err != nil {
			log.Println(err)
			return totalRow, err
		}
	}

	return totalRow, nil
}
