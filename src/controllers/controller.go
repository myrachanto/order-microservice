package controllers

import (
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/microservice/order/src/model"
	service "github.com/myrachanto/microservice/order/src/services"
)

//orderController ..
var (
	OrderController OrdercontrollerInterface = &orderController{}
)

type orderController struct {
	service service.OrderServiceInterface
}
type OrdercontrollerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Additem(c echo.Context) error
	Updateitem(c echo.Context) error
}

func NeworderController(ser service.OrderServiceInterface) OrdercontrollerInterface {
	return &orderController{
		ser,
	}
}

/////////controllers/////////////////
func (controller orderController) Create(c echo.Context) error {
	order := &model.Order{}
	order.Customer = c.FormValue("customer")
	order.Address = c.FormValue("address")
	payamount, err := strconv.ParseFloat(c.FormValue("payamount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code(), httperror)
	}
	order.Amount = payamount

	s, err1 := service.OrderService.Create(order)
	if err1 != nil {
		return c.JSON(err1.Code(), err1)
	}
	return c.JSON(http.StatusCreated, s)
}
func (controller orderController) Additem(c echo.Context) error {
	item := &model.Item{}
	item.Name = c.FormValue("name")
	item.Prodductcode = c.FormValue("productcode")
	item.Ordercode = c.FormValue("ordercode")
	price, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code(), httperror)
	}
	item.Price = price
	priced, err := strconv.ParseInt(c.FormValue("price"), 10, 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code(), httperror)
	}
	item.Quantity = priced

	s, err1 := service.OrderService.Additem(item)
	if err1 != nil {
		return c.JSON(err1.Code(), err1)
	}
	return c.JSON(http.StatusCreated, s)
}
func (controller orderController) GetAll(c echo.Context) error {
	search := string(c.QueryParam("q"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid page number")
		return c.JSON(httperror.Code(), httperror)
	}
	pagesize, err := strconv.Atoi(c.QueryParam("pagesize"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid pagesize")
		return c.JSON(httperror.Code(), httperror)
	}

	results, err3 := service.OrderService.GetAll(search, page, pagesize)
	if err3 != nil {
		return c.JSON(err3.Code(), err3)
	}
	return c.JSON(http.StatusOK, results)
}
func (controller orderController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code(), httperror)
	}
	result, problem := service.OrderService.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code(), problem)
	}
	return c.JSON(http.StatusOK, result)
}

func (controller orderController) Update(c echo.Context) error {
	order := &model.Order{}
	order.Customer = c.FormValue("customer")
	order.Address = c.FormValue("address")
	payamount, err := strconv.ParseFloat(c.FormValue("payamount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code(), httperror)
	}
	order.Amount = payamount
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code(), httperror)
	}
	updatedorder, problem := service.OrderService.Update(id, order)
	if problem != nil {
		return c.JSON(problem.Code(), problem)
	}
	return c.JSON(http.StatusOK, updatedorder)
}
func (controller orderController) Updateitem(c echo.Context) error {
	item := &model.Item{}
	item.Name = c.FormValue("name")
	item.Prodductcode = c.FormValue("productcode")
	item.Ordercode = c.FormValue("ordercode")
	price, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code(), httperror)
	}
	item.Price = price
	priced, err := strconv.ParseInt(c.FormValue("price"), 10, 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid selling price")
		return c.JSON(httperror.Code(), httperror)
	}
	item.Quantity = priced

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code(), httperror)
	}
	s, err1 := service.OrderService.Updateitem(id, item)
	if err1 != nil {
		return c.JSON(err1.Code(), err1)
	}
	return c.JSON(http.StatusCreated, s)
}

func (controller orderController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code(), httperror)
	}
	success, failure := service.OrderService.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure)
	}
	return c.JSON(http.StatusOK, success)

}
