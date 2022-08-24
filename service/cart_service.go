package service

import (
	"context"
	"database/sql"
	"interview-telkom-6/entity"
	"interview-telkom-6/repository/persistence"
	"interview-telkom-6/request"
	"interview-telkom-6/response"
	"interview-telkom-6/util"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CartService struct {
	ctx             context.Context
	cartRepo        *persistence.CartRepository
	cartProductRepo *persistence.CartProductRepository
	productRepo     persistence.ProductRepository
}

func NewCartService(
	ctx context.Context, cartRepo *persistence.CartRepository,
	cartProductRepo *persistence.CartProductRepository,
	productRepo persistence.ProductRepository,
) *CartService {
	return &CartService{
		ctx: ctx, cartRepo: cartRepo, cartProductRepo: cartProductRepo, productRepo: productRepo,
	}
}

func (s *CartService) DeleteProduct(ctx context.Context, productID string) error {
	if productID == "" {
		return &util.BadRequestError{Message: "product id not found"}
	}

	cpBuilder := persistence.QueryBuilderCriteria{}
	cpBuilder.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"product_id": productID}}}}
	cp, err := s.cartProductRepo.Get(ctx, &cpBuilder)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return &util.BadRequestError{Message: "product not found"}
		}

		return err
	}

	err = s.cartProductRepo.Delete(ctx, &cp)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *CartService) Find(ctx context.Context, req *request.CartCriteria) (
	res *response.CartResponse, err error,
) {
	builder := persistence.QueryBuilderCriteria{}
	builder.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"full_name": req.FullName}}}}
	if req.FullName == "" {
		return res, &util.BadRequestError{Message: "full name can't be null"}
	}

	result, err := s.cartRepo.Get(ctx, &builder)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, err
	}

	data := response.CartResponse{
		ID:       result.ID,
		FullName: result.FullName,
		Products: make([]response.CartResponseProduct, 0),
	}

	cpBuilder := persistence.QueryBuilderCriteria{}
	cpBuilder.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"cart_id": result.ID}}}}
	cpBuilder.Select = []string{"cart_products.*"}
	if req.Quantity != "" {
		cpBuilder.Where.And = append(cpBuilder.Where.And, squirrel.And{squirrel.Eq{"quantity": req.Quantity}})
	}

	if req.ProductName != "" {
		cpBuilder.Join.LeftJoin = []string{"products ON products.id = cart_products.product_id"}
		cpBuilder.Where.And = append(
			cpBuilder.Where.And, squirrel.And{
				squirrel.ILike{
					"products.name": "%" + req.
						ProductName + "%",
				},
			},
		)
	}
	cartProducts, err := s.cartProductRepo.Find(ctx, &cpBuilder)
	if err != nil {
		log.Println(err)
		return res, err
	}

	for _, cp := range cartProducts {
		p := response.CartResponseProduct{
			Quantity: cp.Quantity,
		}

		pBuilder := persistence.QueryBuilderCriteria{}
		pBuilder.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"id": cp.ProductID}}}}
		product, err := s.productRepo.Get(ctx, &pBuilder)
		if err != nil {
			log.Println(err)
			return res, err
		}

		p.Product.ID = product.ID
		p.Product.Name = product.Name
		p.Product.Price = product.Price
		p.Product.Description = product.Description
		p.Product.IsDiscount = product.IsDiscount
		p.Product.StartDateDiscount = &product.StartDateDiscount.Time
		p.Product.EndDateDiscount = &product.EndDateDiscount.Time

		data.Products = append(data.Products, p)
	}

	return &data, nil
}

func (s *CartService) Store(ctx context.Context, req *request.CartAddRequest) (res *response.CartResponse, err error) {
	tx, err := s.ctx.Value("db").(*sqlx.DB).Beginx()
	if err != nil {
		log.Println(err)
		return res, err
	}

	builder := persistence.QueryBuilderCriteria{}
	builder.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"full_name": req.FullName}}}}
	checkCart, err := s.cartRepo.Get(ctx, &builder)
	if err != sql.ErrNoRows && err != nil {
		log.Println(err)
		return res, err
	}

	pBuilder := persistence.QueryBuilderCriteria{}
	pBuilder.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"id": req.Product.ProductID}}}}
	_, err = s.productRepo.Get(ctx, &pBuilder)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return res, &util.BadRequestError{Message: "product not found"}
		}
		return res, err
	}

	cartTx := s.cartRepo.WithTx(tx)
	cartProductTx := s.cartProductRepo.WithTx(tx)

	// if cart already exist, just insert products to cart
	if checkCart.ID != uuid.Nil {
		cPBuilder := persistence.QueryBuilderCriteria{}
		cPBuilder.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"product_id": req.Product.ProductID}}}}
		cp, err := s.cartProductRepo.Get(ctx, &cPBuilder)
		if err != sql.ErrNoRows && err != nil {
			log.Println(err)
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				return res, err
			}
			return res, err
		}

		// if product isn't exist in cart, just insert product id to cart
		if err == sql.ErrNoRows {
			err = s.insertCartProduct(ctx, checkCart.ID, &req.Product, cartProductTx)
			if err != nil {
				log.Println(err)
				if err := tx.Rollback(); err != nil {
					log.Println(err)
					return res, err
				}
				return res, err
			}
			if err := tx.Commit(); err != nil {
				log.Println(err)
				return res, err
			}
			return res, nil
		}

		// if product is exist in cart, just update quantity
		cp.Quantity += req.Product.Quantity
		_, err = s.cartProductRepo.Update(ctx, &cp)
		if err != nil {
			log.Println(err)
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				return res, err
			}
			return res, err
		}
		if err := tx.Commit(); err != nil {
			log.Println(err)
			return res, err
		}

		reqFind := request.CartCriteria{FullName: req.FullName}
		res, err = s.Find(ctx, &reqFind)
		if err != nil {
			log.Println(err)
			return res, err
		}

		return res, nil
	}

	// if cart isn't exist, create cart and insert products to cart
	cartEntity := entity.Cart{FullName: req.FullName}
	cart, err := cartTx.Store(ctx, &cartEntity)
	if err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			return res, err
		}
		return res, err
	}

	err = s.insertCartProduct(ctx, cart.ID, &req.Product, cartProductTx)
	if err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			return res, err
		}
		return res, err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return res, err
	}

	reqFind := request.CartCriteria{FullName: req.FullName}
	res, err = s.Find(ctx, &reqFind)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (s *CartService) insertCartProduct(
	ctx context.Context, cartID uuid.UUID, req *request.CartAddProductRequest,
	cartRepoProduct *persistence.CartProductRepository,
) (err error) {
	cp := entity.CartProduct{
		CartID:    cartID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	_, err = cartRepoProduct.Store(ctx, &cp)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
