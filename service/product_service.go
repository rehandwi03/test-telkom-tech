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
	"time"

	"github.com/Masterminds/squirrel"
)

type ProductService struct {
	productRepo *persistence.ProductRepository
}

func NewProductService(
	productRepo *persistence.ProductRepository,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) Find(ctx context.Context, req *request.ProductCriteria) (
	res *util.PaginationResponse, err error,
) {
	builder := persistence.QueryBuilderCriteria{}
	builder.Where = &persistence.Where{}

	if req.Search != "" {
		and := squirrel.And{squirrel.ILike{"name": "%" + req.Search + "%"}}
		builder.Where.And = append(builder.Where.And, and)
	}

	page := uint64(req.Page)
	limit := uint64(req.Limit)
	offset := (page - 1) * limit

	builder.Limit = &limit
	builder.Offset = &offset
	responses := make([]response.ProductResponse, 0)

	results, err := s.productRepo.Find(ctx, &builder)
	if err != nil {
		log.Println(err)
		return res, err
	}

	for _, val := range results {

		data := response.ProductResponse{
			ID:                val.ID,
			Name:              val.Name,
			Price:             val.Price,
			Description:       val.Description,
			IsDiscount:        val.IsDiscount,
			StartDateDiscount: &val.StartDateDiscount.Time,
			EndDateDiscount:   &val.EndDateDiscount.Time,
			DiscountValue:     val.DiscountValue.Float64,
		}

		responses = append(responses, data)

	}

	totalRow, err := s.productRepo.Count(ctx, &builder)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return util.BuildPagination(req.Pagination, responses, totalRow), nil
}

func (s *ProductService) Store(ctx context.Context, req *request.ProductAddRequest) (err error) {
	// check product
	productBuilder := persistence.QueryBuilderCriteria{}
	productBuilder.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"name": req.Name}}}}
	product, err := s.productRepo.Get(ctx, &productBuilder)
	if err != sql.ErrNoRows && err != nil {
		log.Println(err)
		return err
	}

	if product.Name != "" {
		return &util.BadRequestError{Message: "produk sudah ada"}
	}

	// insert product
	productEntity := entity.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}
	if req.IsDiscount {
		startDD, err := time.Parse("2006-01-02", req.StartDateDiscount)
		if err != nil {
			log.Println(err)
			return err
		}

		endDD, err := time.Parse("2006-01-02", req.EndDateDiscount)
		if err != nil {
			log.Println(err)
			return err
		}

		productEntity.StartDateDiscount = sql.NullTime{
			Time:  startDD,
			Valid: true,
		}
		productEntity.EndDateDiscount = sql.NullTime{
			Time:  endDD,
			Valid: true,
		}
		productEntity.IsDiscount = req.IsDiscount
		productEntity.DiscountValue = sql.NullFloat64{Float64: req.DiscountValue, Valid: true}
	}

	_, err = s.productRepo.Store(ctx, &productEntity)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
