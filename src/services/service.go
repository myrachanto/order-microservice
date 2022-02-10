package service

import (
	"fmt"

	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/microservice/order/src/model"
	r "github.com/myrachanto/microservice/order/src/repository"
)

//orderService ...
var (
	OrderService OrderServiceInterface = &orderService{}
)

type OrderServiceInterface interface {
	Create(order *model.Order) (string, httperors.HttpErr)
	Additem(order *model.Item) (string, httperors.HttpErr)
	GetOne(code int) (*model.Order, httperors.HttpErr)
	GetAll(search string, page, pagesize int) ([]model.Order, httperors.HttpErr)
	Update(id int, order *model.Order) (*model.Order, httperors.HttpErr)
	Updateitem(id int, item *model.Item) (*model.Item, httperors.HttpErr)
	Delete(id int) (string, httperors.HttpErr)
}

type orderService struct {
	repository r.OrderRepoInterface
}

func NeworderService(repo r.OrderRepoInterface) OrderServiceInterface {
	return &orderService{
		repo,
	}
}

func (service orderService) Create(order *model.Order) (string, httperors.HttpErr) {
	if err := order.Validate(); err != nil {
		return "", err
	}
	s, err1 := r.Orderrepo.Create(order)
	if err1 != nil {
		return "", err1
	}
	return s, nil

}
func (service orderService) Additem(item *model.Item) (string, httperors.HttpErr) {
	if err := item.Validate(); err != nil {
		return "", err
	}
	s, err1 := r.Orderrepo.Additem(item)
	if err1 != nil {
		return "", err1
	}
	return s, nil

}
func (service orderService) GetOne(code int) (*model.Order, httperors.HttpErr) {
	order, err1 := r.Orderrepo.GetOne(code)
	if err1 != nil {
		return nil, err1
	}
	return order, nil
}

func (service orderService) GetAll(search string, page, pagesize int) ([]model.Order, httperors.HttpErr) {
	results, err := r.Orderrepo.GetAll(search, page, pagesize)
	return results, err
}

// func (service orderService) UpdateRole(code, admin, supervisor, employee, level, ordercode string) (string, *httperors.HttpError) {
// 	order, err1 := r.Orderrepo.UpdateRole(code, admin, supervisor, employee, level, ordercode)
// 	return order, err1
// }

func (service orderService) Update(id int, order *model.Order) (*model.Order, httperors.HttpErr) {
	fmt.Println("update1-controller")
	fmt.Println(id)
	order, err1 := r.Orderrepo.Update(id, order)
	if err1 != nil {
		return nil, err1
	}

	return order, nil
}
func (service orderService) Updateitem(id int, item *model.Item) (*model.Item, httperors.HttpErr) {
	fmt.Println("update1-controller")
	fmt.Println(id)
	item, err1 := r.Orderrepo.Updateitem(id, item)
	if err1 != nil {
		return nil, err1
	}

	return item, nil
}
func (service orderService) Delete(id int) (string, httperors.HttpErr) {
	success, failure := r.Orderrepo.Delete(id)
	return success, failure
}
