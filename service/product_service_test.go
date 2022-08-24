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
	"log"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	mockCrtl := gomock.NewController(t)
	defer mockCrtl.Finish()

	productMock := mocks.NewMockProductRepository(mockCrtl)

	product := entity.Product{
		// ID:                uuid.New(),
		Name:              "Makanan",
		Price:             10000,
		Description:       "Makanan Enak",
		IsDiscount:        false,
		DiscountValue:     sql.NullFloat64{},
		StartDateDiscount: sql.NullTime{},
		EndDateDiscount:   sql.NullTime{},
	}
	req := request.ProductAddRequest{
		Name:              "Makanan",
		Price:             10000,
		Description:       "Makanan Enak",
		IsDiscount:        false,
		StartDateDiscount: "",
		EndDateDiscount:   "",
		DiscountValue:     0,
	}

	resProduct := entity.Product{}
	w := persistence.QueryBuilderCriteria{}
	w.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"name": req.Name}}}}
	productMock.EXPECT().Get(context.TODO(), &w).Return(resProduct, nil)
	productMock.EXPECT().Store(context.TODO(), &product).Return(product, nil)

	productSvc := service.NewProductService(productMock)

	err := productSvc.Store(context.TODO(), &req)
	assert.NoError(t, err)
}

func TestStoreWithDiscount(t *testing.T) {
	mockCrtl := gomock.NewController(t)
	defer mockCrtl.Finish()

	productMock := mocks.NewMockProductRepository(mockCrtl)

	req := request.ProductAddRequest{
		Name:              "Makanan",
		Price:             10000,
		Description:       "Makanan Enak",
		IsDiscount:        true,
		StartDateDiscount: "2022-08-24",
		EndDateDiscount:   "2022-08-30",
		DiscountValue:     10000.00,
	}

	sdParse, err := time.Parse("2006-01-02", req.StartDateDiscount)
	assert.NoError(t, err)
	edParse, err := time.Parse("2006-01-02", req.EndDateDiscount)
	assert.NoError(t, err)

	product := entity.Product{
		Name:              req.Name,
		Price:             req.Price,
		Description:       req.Description,
		IsDiscount:        req.IsDiscount,
		DiscountValue:     sql.NullFloat64{Float64: req.DiscountValue, Valid: true},
		StartDateDiscount: sql.NullTime{Time: sdParse, Valid: true},
		EndDateDiscount:   sql.NullTime{Time: edParse, Valid: true},
	}

	resProduct := entity.Product{}
	w := persistence.QueryBuilderCriteria{}
	w.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"name": req.Name}}}}
	productMock.EXPECT().Get(context.TODO(), &w).Return(resProduct, nil)
	productMock.EXPECT().Store(context.TODO(), &product).Return(product, nil)

	productSvc := service.NewProductService(productMock)

	err = productSvc.Store(context.TODO(), &req)
	assert.NoError(t, err)
}

func TestStoreErrorProductExist(t *testing.T) {
	mockCrtl := gomock.NewController(t)
	defer mockCrtl.Finish()

	productMock := mocks.NewMockProductRepository(mockCrtl)

	req := request.ProductAddRequest{
		Name:              "Makanan",
		Price:             10000,
		Description:       "Makanan Enak",
		IsDiscount:        true,
		StartDateDiscount: "2022-08-24",
		EndDateDiscount:   "2022-08-30",
		DiscountValue:     10000.00,
	}

	sdParse, err := time.Parse("2006-01-02", req.StartDateDiscount)
	assert.NoError(t, err)
	edParse, err := time.Parse("2006-01-02", req.EndDateDiscount)
	assert.NoError(t, err)

	product := entity.Product{
		Name:              req.Name,
		Price:             req.Price,
		Description:       req.Description,
		IsDiscount:        req.IsDiscount,
		DiscountValue:     sql.NullFloat64{Float64: req.DiscountValue, Valid: true},
		StartDateDiscount: sql.NullTime{Time: sdParse, Valid: true},
		EndDateDiscount:   sql.NullTime{Time: edParse, Valid: true},
	}

	w := persistence.QueryBuilderCriteria{}
	w.Where = &persistence.Where{And: []squirrel.And{{squirrel.Eq{"name": req.Name}}}}
	productMock.EXPECT().Get(context.TODO(), &w).Return(product, nil)

	productSvc := service.NewProductService(productMock)

	err = productSvc.Store(context.TODO(), &req)
	log.Println(err)
	assert.Error(t, err)
	assert.NotNil(t, err)
}

func TestFind(t *testing.T) {
	ctx := context.TODO()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	productMock := mocks.NewMockProductRepository(mockCtrl)

	req := request.ProductCriteria{}
	req.Limit = 1
	req.Page = 1

	b := persistence.QueryBuilderCriteria{}
	b.Where = &persistence.Where{}
	limit := uint64(req.Limit)
	page := uint64(req.Page)
	offset := (page - 1) * limit
	b.Limit = &limit
	b.Offset = &offset

	res := []entity.Product{
		{
			ID:                uuid.New(),
			Name:              "Makanan",
			Price:             10000,
			Description:       "Makanan Enak",
			IsDiscount:        false,
			StartDateDiscount: sql.NullTime{},
			EndDateDiscount:   sql.NullTime{},
			DiscountValue:     sql.NullFloat64{},
		},
	}
	var totalRow int64 = 1
	productMock.EXPECT().Find(ctx, &b).Return(res, nil)
	productMock.EXPECT().Count(ctx, &b).Return(totalRow, nil)

	productSvc := service.NewProductService(productMock)

	_, err := productSvc.Find(ctx, &req)
	assert.NoError(t, err)
}
