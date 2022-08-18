package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"interview-telkom-6/request"
	"interview-telkom-6/response"
	"interview-telkom-6/service"
	"interview-telkom-6/util"
	"log"
	"net/http"
)

type productHandler struct {
	productSvc *service.ProductService
}

func NewProductHandler(router *gin.RouterGroup, productSvc *service.ProductService) {
	h := productHandler{productSvc: productSvc}

	path := "/products"
	router.POST(path, h.Store)
	router.GET(path, h.Find)
}

func (h *productHandler) Find(c *gin.Context) {
	req := new(request.ProductCriteria)

	pagination := util.GeneratePaginationFromRequest(c)
	req.Search = c.Query("search")
	req.Popular = c.Query("popular")
	req.Pagination = pagination

	res, err := h.productSvc.Find(c, req)
	if err != nil {
		util.BuildErrorAPI(c, err)
		return
	}

	c.JSONP(http.StatusOK, res)
	return
}

func (h *productHandler) Store(c *gin.Context) {
	req := new(request.ProductAddRequest)
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(
			http.StatusBadRequest, response.ErrorResponse{
				Message: "StatusBadRequest",
				Error:   err.Error(),
				Status:  "failed",
			},
		)

		return
	}

	if err := binding.Validator.ValidateStruct(req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.ErrorResponse{
				Message: "StatusBadRequest",
				Error:   err.Error(),
				Status:  "failed",
			},
		)
		return
	}

	err := h.productSvc.Store(c, req)
	if err != nil {
		log.Println(err)
		util.BuildErrorAPI(c, err)
		return
	}

	c.JSONP(http.StatusOK, response.SuccessResponse{Status: "success", Message: "success saved data", Data: nil})
	return
}
