package persistence

import (
	"context"
	"fmt"
	"interview-telkom-6/entity"
	"log"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	WithTx(conn *sqlx.Tx) ProductRepository
	Get(ctx context.Context, builder *QueryBuilderCriteria) (
		res entity.Product, err error,
	)
	Find(ctx context.Context, builder *QueryBuilderCriteria) (
		res []entity.Product, err error,
	)
	Store(ctx context.Context, data *entity.Product) (res entity.Product, err error)
	Update(ctx context.Context, data *entity.Product) (res entity.Product, err error)
	Delete(ctx context.Context, data *entity.Product) (err error)
	Count(ctx context.Context, builder *QueryBuilderCriteria) (totalRow int64, err error)
}

type productRepository struct {
	Conn      Queryer
	TableName string
}

func NewProductRepository(conn *sqlx.DB) ProductRepository {
	return &productRepository{Conn: conn, TableName: "products"}
}

func (r productRepository) WithTx(conn *sqlx.Tx) ProductRepository {
	if conn == nil {
		log.Println("transaction database not found")
		return &r
	}

	return &productRepository{Conn: conn, TableName: "products"}
}

func (r productRepository) Get(ctx context.Context, builder *QueryBuilderCriteria) (
	res entity.Product, err error,
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

func (r productRepository) Find(ctx context.Context, builder *QueryBuilderCriteria) (
	res []entity.Product, err error,
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

func (r productRepository) Store(ctx context.Context, data *entity.Product) (res entity.Product, err error) {
	data.GenerateUUID()
	query := fmt.Sprintf(
		"INSERT INTO %s (id,  name, price, description, is_discount, "+
			"discount_value, start_date_discount, end_date_discount) "+
			"VALUES (:id, :name, :price, :description, :is_discount, :discount_value, "+
			":start_date_discount, :end_date_discount)",
		r.TableName,
	)
	log.Println(query)
	_, err = r.Conn.NamedExecContext(ctx, query, data)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return *data, err
}

func (r productRepository) Update(ctx context.Context, data *entity.Product) (res entity.Product, err error) {
	query := fmt.Sprintf(
		"UPDATE %s SET name=:name, price=:price, description=:description, "+
			"is_discount=:is_discount, start_date_discount=:start_date_discount, end_date_discount=:end_date_discount WHERE id=:id",
		r.TableName,
	)
	log.Println(query)
	_, err = r.Conn.NamedExecContext(ctx, query, data)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return *data, nil
}

func (r productRepository) Delete(ctx context.Context, data *entity.Product) (err error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.TableName)
	log.Println(query)
	_, err = r.Conn.Exec(query, data.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r productRepository) Count(ctx context.Context, builder *QueryBuilderCriteria) (totalRow int64, err error) {
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
