package handler

import (
	"fmt"
	"interview-telkom-6/request"
	"interview-telkom-6/response"
	"interview-telkom-6/service"
	"interview-telkom-6/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type cartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(router *gin.RouterGroup, cartService *service.CartService) {
	h := cartHandler{cartService: cartService}

	path := "/carts"
	router.POST(path, h.Store)
	router.GET(path, h.Find)
	router.DELETE(fmt.Sprintf("%s/:product_id", path), h.DeleteProduct)
}

func (h *cartHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("product_id")
	err := h.cartService.DeleteProduct(c, productID)
	if err != nil {
		log.Println(err)
		util.BuildErrorAPI(c, err)
		return
	}
	c.JSONP(http.StatusOK, response.SuccessResponse{Status: "success", Message: "success delete data"})
	return
}

func (h *cartHandler) Find(c *gin.Context) {
	req := new(request.CartCriteria)
	req.FullName = c.Query("full_name")
	req.ProductName = c.Query("product_name")
	req.Quantity = c.Query("quantity")

	res, err := h.cartService.Find(c, req)
	if err != nil {
		log.Println(err)
		util.BuildErrorAPI(c, err)
		return
	}

	c.JSONP(http.StatusOK, response.SuccessResponse{Status: "success", Message: "success get data", Data: res})
	return
}

func (h *cartHandler) Store(c *gin.Context) {
	req := new(request.CartAddRequest)
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

	res, err := h.cartService.Store(c, req)
	if err != nil {
		log.Println(err)
		util.BuildErrorAPI(c, err)
		return
	}

	c.JSONP(http.StatusOK, response.SuccessResponse{Status: "success", Message: "success saved data", Data: res})
	return
}
