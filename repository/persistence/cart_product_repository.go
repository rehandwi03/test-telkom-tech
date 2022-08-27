package persistence

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"interview-telkom-6/entity"
	"log"
)

type cartProductRepo struct {
	Conn      Queryer
	TableName string
}

type CartProductRepository interface {
	WithTx(conn *sqlx.Tx) CartProductRepository
	Get(ctx context.Context, builder *QueryBuilderCriteria) (
		res entity.CartProduct, err error,
	)
	Find(ctx context.Context, builder *QueryBuilderCriteria) (
		res []entity.CartProduct, err error,
	)
	Store(ctx context.Context, data *entity.CartProduct) (
		res entity.CartProduct, err error,
	)
	Update(ctx context.Context, data *entity.CartProduct) (
		res entity.CartProduct, err error,
	)
	Delete(ctx context.Context, data *entity.CartProduct) (err error)
	Count(ctx context.Context, builder *QueryBuilderCriteria) (totalRow int64, err error)
}

func NewCartProductRepo(conn *sqlx.DB) CartProductRepository {
	return &cartProductRepo{Conn: conn, TableName: "cart_products"}
}

func (r cartProductRepo) WithTx(conn *sqlx.Tx) CartProductRepository {
	if conn == nil {
		log.Println("transaction database not found")
		return &r
	}

	return &cartProductRepo{Conn: conn}
}

func (r cartProductRepo) Get(ctx context.Context, builder *QueryBuilderCriteria) (
	res entity.CartProduct, err error,
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

func (r cartProductRepo) Find(ctx context.Context, builder *QueryBuilderCriteria) (
	res []entity.CartProduct, err error,
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

func (r cartProductRepo) Store(ctx context.Context, data *entity.CartProduct) (
	res entity.CartProduct, err error,
) {
	query := fmt.Sprintf(
		"INSERT INTO cart_products (cart_id, product_id, quantity) VALUES (:cart_id, :product_id, :quantity)",
	)
	log.Println(query)
	_, err = r.Conn.NamedExecContext(ctx, query, data)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return *data, err
}

func (r cartProductRepo) Update(ctx context.Context, data *entity.CartProduct) (
	res entity.CartProduct, err error,
) {
	query := fmt.Sprintf(
		"UPDATE cart_products SET product_id=:product_id, quantity=:quantity WHERE cart_id=:cart_id",
	)
	log.Println(query)
	_, err = r.Conn.NamedExecContext(ctx, query, data)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return *data, nil
}

func (r cartProductRepo) Delete(ctx context.Context, data *entity.CartProduct) (err error) {
	query := fmt.Sprintf("DELETE FROM cart_products WHERE product_id = $1")
	log.Println(query)
	_, err = r.Conn.Exec(query, data.ProductID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r cartProductRepo) Count(ctx context.Context, builder *QueryBuilderCriteria) (totalRow int64, err error) {
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
