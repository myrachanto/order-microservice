package repository

import (
	"encoding/json"
	"strconv"

	fetcher "github.com/myrachanto/Fetcher"
	httperors "github.com/myrachanto/custom-http-error"
	pubsub "github.com/myrachanto/microservice/order/src/events"
	"github.com/myrachanto/microservice/order/src/model"
)

//orderrepo ...
var (
	Orderrepo OrderRepoInterface = &orderrepo{}
)

type Key struct {
	EncryptionKey string `mapstructure:"EncryptionKey"`
}

type orderrepo struct {
	Token string
}

func Tokenin(token string) *orderrepo {
	return &orderrepo{Token: token}
}

type OrderRepoInterface interface {
	Create(order *model.Order) (string, httperors.HttpErr)
	Additem(order *model.Item) (string, httperors.HttpErr)
	all() (t []model.Order, r httperors.HttpErr)
	GetOne(id int) (*model.Order, httperors.HttpErr)
	orderExistbycode(code string) bool
	orderbycode(code string) *model.Order
	GetAll(search string, page, pagesize int) ([]model.Order, httperors.HttpErr)
	Update(id int, order *model.Order) (*model.Order, httperors.HttpErr)
	Updateitem(id int, item *model.Item) (*model.Item, httperors.HttpErr)
	Delete(id int) (string, httperors.HttpErr)
	geneCode() (string, httperors.HttpErr)
	orderExist(email string) bool
	orderExistByid(id int) bool
}

func NeworderRepo() *orderrepo {
	return &orderrepo{}
}
func (orderRepo orderrepo) Create(order *model.Order) (string, httperors.HttpErr) {
	if err := order.Validate(); err != nil {
		return "", err
	}
	code, x := orderRepo.geneCode()
	if x != nil {
		return "", x
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	order.Ordercode = code
	GormDB.Create(&order)
	IndexRepo.DbClose(GormDB)
	//pub sub the information about the order creation being filed up
	pubsub.Produce("order_topic", "Order_created", order)
	return "order created successifully", nil
}
func (o orderrepo) Additem(item *model.Item) (string, httperors.HttpErr) {
	if err := item.Validate(); err != nil {
		return "", err
	}
	//fetch data from product microservice
	fetch := &fetcher.Fetcher{Endpoint: "http://product_backend:4001", Token: o.Token}
	response, err := fetch.Request("GET", "/products/"+item.Prodductcode, nil)
	if err != nil {
		return "", httperors.NewNotFoundError("could not get the data on that product")
	}
	product := &model.Product{}
	json.NewDecoder(response.Body).Decode(product)
	item.Name = product.Name
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	GormDB.Create(&item)
	//pub sub the information about the order creation
	IndexRepo.DbClose(GormDB)
	return "order created successifully", nil
}
func (orderRepo orderrepo) all() (t []model.Order, r httperors.HttpErr) {

	order := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&order).Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (orderRepo orderrepo) geneCode() (string, httperors.HttpErr) {
	order := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	err := GormDB.Last(&order)
	if err.Error != nil {
		var c1 uint = 1
		code := "orderCode" + strconv.FormatUint(uint64(c1), 10)
		return code, nil
	}
	c1 := order.ID + 1
	code := "orderCode" + strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil

}
func (orderRepo orderrepo) GetOne(id int) (*model.Order, httperors.HttpErr) {
	ok := orderRepo.orderExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("order with that code does not exists!")
	}
	order := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}

	GormDB.Model(&order).Where("id = ?", id).First(&order)
	IndexRepo.DbClose(GormDB)
	return &order, nil
}
func (orderRepo orderrepo) orderExistbycode(code string) bool {
	u := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	GormDB.Where("ordercode = ?", code).First(&u)
	if u.ID == 0 {
		return false
	}
	IndexRepo.DbClose(GormDB)
	return true

}
func (orderRepo orderrepo) orderbycode(code string) *model.Order {
	u := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	GormDB.Where("ordercode = ?", code).First(&u)
	if u.ID == 0 {
		return nil
	}
	IndexRepo.DbClose(GormDB)
	return &u

}
func (orderRepo orderrepo) GetAll(search string, page, pagesize int) ([]model.Order, httperors.HttpErr) {
	results := []model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	if search == "" {
		GormDB.Find(&results)
	}
	// db.Scopes(Paginate(r)).Find(&orders)
	GormDB.Scopes(Paginate(page, pagesize)).Where("name LIKE ?", "%"+search+"%").Or("email LIKE ?", "%"+search+"%").Or("company LIKE ?", "%"+search+"%").Find(&results)

	IndexRepo.DbClose(GormDB)
	return results, nil
}

func (orderRepo orderrepo) Update(id int, order *model.Order) (*model.Order, httperors.HttpErr) {
	ok := orderRepo.orderExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("order with that id does not exists!")
	}

	uorder := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&order).Where("id = ?", id).First(&uorder)
	if order.Customer == "" {
		order.Customer = uorder.Customer
	}
	if order.Phone == "" {
		order.Phone = uorder.Phone
	}
	if order.Address == "" {
		order.Address = uorder.Address
	}
	if order.Amount == 0 {
		order.Amount = uorder.Amount
	}
	if order.Usercode == "" {
		order.Usercode = uorder.Usercode
	}
	if order.Ordercode == "" {
		order.Ordercode = uorder.Ordercode
	}
	GormDB.Save(&order)

	IndexRepo.DbClose(GormDB)

	return order, nil
}
func (i orderrepo) Updateitem(id int, item *model.Item) (*model.Item, httperors.HttpErr) {

	uitem := model.Item{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&item).Where("id = ?", id).First(&uitem)
	if item.Name == "" {
		item.Name = uitem.Name
	}
	if item.Price == 0 {
		item.Price = uitem.Price
	}
	if item.Prodductcode == "" {
		item.Prodductcode = uitem.Prodductcode
	}
	if item.Quantity == 0 {
		item.Quantity = uitem.Quantity
	}
	if item.Total == 0 {
		item.Total = uitem.Total
	}
	if item.Ordercode == "" {
		item.Ordercode = uitem.Ordercode
	}
	GormDB.Save(&item)

	IndexRepo.DbClose(GormDB)

	return item, nil
}
func (orderRepo orderrepo) Delete(id int) (string, httperors.HttpErr) {
	ok := orderRepo.orderExistByid(id)
	if !ok {
		return "", httperors.NewNotFoundError("order with that id does not exists!")
	}
	order := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	GormDB.Model(&order).Where("id = ?", id).First(&order)
	GormDB.Delete(order)
	IndexRepo.DbClose(GormDB)
	return "deleted successfully", nil
}
func (orderRepo orderrepo) orderExist(email string) bool {
	order := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&order, "email =?", email)
	if res.Error != nil {
		return false
	}
	IndexRepo.DbClose(GormDB)
	return true

}
func (orderRepo orderrepo) orderExistByid(id int) bool {
	order := model.Order{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&order, "id =?", id)
	if res.Error != nil {
		return false
	}
	IndexRepo.DbClose(GormDB)
	return true

}
