package model

import (
	"time"

	httperors "github.com/myrachanto/custom-http-error"
	"gorm.io/gorm"
)

var ExpiresAt = time.Now().Add(time.Minute * 100000).Unix()

type Order struct {
	Customer    string  `json:"customer"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	Instruction string  `json:"instruction"`
	Amount      float64 `json:"amount"`
	Item        []Item  `json:"item" gorm:"-"`
	Usercode    string  `json:"usercode"`
	Ordercode   string  `json:"code"`
	gorm.Model
}
type Item struct {
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Prodductcode string  `json:"prodductcode"`
	Quantity     int64   `json:"quantity"`
	Total        float64 `json:"total"`
	Ordercode    string  `json:"ordercode"`
	gorm.Model
}

//other micorservice structs
type Product struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BuyPrice    float64 `json:"buy_price"`
	SellPrice   float64 `json:"sell_price"`
	Quantity    int64   `json:"quantity"`
	Picture     string  `json:"picture"`
	Available   bool    `json:"available"`
	Usercode    string  `json:"usercode"`
	Productcode string  `json:"code"`
	gorm.Model
}

//Validate ..
func (order Order) Validate() httperors.HttpErr {
	if order.Customer == "" {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if order.Address == "" {
		return httperors.NewNotFoundError("Invalid Address")
	}
	if order.Amount == 0 {
		return httperors.NewNotFoundError("Invalid Amount")
	}
	return nil
}
func (item Item) Validate() httperors.HttpErr {
	if item.Name == "" {
		return httperors.NewNotFoundError("Invalid  Name")
	}
	if item.Price == 0 {
		return httperors.NewNotFoundError("Invalid Price")
	}
	if item.Quantity == 0 {
		return httperors.NewNotFoundError("Invalid Quantity")
	}
	return nil
}
