package service_test

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"interview-telkom-6/entity"
	"interview-telkom-6/repository/persistence"
	"interview-telkom-6/repository/persistence/mocks"
	"interview-telkom-6/request"
	"interview-telkom-6/service"
	"testing"
)

func TestFindCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	cartRepo := mocks.NewMockCartRepository(ctrl)
	cartProductRepo := mocks.NewMockCartProductRepository(ctrl)
	productRepo := mocks.NewMockProductRepository(ctrl)

	res := entity.Cart{
		ID:       uuid.New(),
		FullName: "Rehan",
	}

	req := request.CartCriteria{
		FullName: "Rehan",
	}

	resCP := []entity.CartProduct{
		{
			CartID:    uuid.New(),
			ProductID: uuid.New(),
			Quantity:  1,
		},
	}

	resProduct := entity.Product{
		ID:                uuid.New(),
		Name:              "Makanan",
		Price:             1000,
		Description:       "",
		IsDiscount:        false,
		DiscountValue:     sql.NullFloat64{},
		StartDateDiscount: sql.NullTime{},
		EndDateDiscount:   sql.NullTime{},
	}

	ctx := context.TODO()
	b := persistence.QueryBuilderCriteria{}
	b.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"full_name": req.FullName}}}}
	cartRepo.EXPECT().Get(ctx, &b).Return(res, nil)

	bc := persistence.QueryBuilderCriteria{}
	bc.Where = &persistence.Where{}
	bc.Select = []string{"cart_products.*"}
	// bc.Join.LeftJoin = []string{"products ON products.id = cart_products.product_id"}
	// bc.Where.And = append(bc.Where.And, squirrel.And{squirrel.ILike{"product_name": "%" + req.ProductName + "%"}})
	bc.Where.And = append(bc.Where.And, squirrel.And{squirrel.Eq{"cart_id": res.ID}})
	cartProductRepo.EXPECT().Find(ctx, &bc).Return(resCP, nil)

	for _, val := range resCP {
		bp := persistence.QueryBuilderCriteria{}
		bp.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"id": val.ProductID}}}}
		productRepo.EXPECT().Get(ctx, &bp).Return(resProduct, nil)
	}

	cartSvc := service.NewCartService(ctx, cartRepo, cartProductRepo, productRepo)
	_, err := cartSvc.Find(ctx, &req)
	assert.NoError(t, err)

}
