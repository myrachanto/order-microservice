package service

import (
	"testing"

	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/microservice/order/src/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

var (
	Order = &model.Order{
		Customer: "Walking in",
		Address:  "Lake view",
		Phone:    "1234567g",
		Usercode: "user1234567",
		Amount:   1000,
	}
)

type ProductMockInterface interface {
	Create(Product *model.Order) (*model.Order, httperors.HttpErr)
	GetOne(id string) (*model.Order, httperors.HttpErr)
	GetAll() ([]*model.Order, httperors.HttpErr)
	Update(code string, Product *model.Order) httperors.HttpErr
}

func (mock MockRepository) Create(Product *model.Order) (*model.Order, httperors.HttpErr) {
	args := mock.Called()
	result := args.Get(0)
	Product, err := result.(*model.Order), args.Error(1)
	if err != nil {
		return nil, httperors.NewNotFoundError("Something went wrong creating the resourse")
	}
	return Product, nil
}
func (mock MockRepository) GetOne() (*model.Order, httperors.HttpErr) {
	args := mock.Called()
	result := args.Get(0)
	Product, err := result.(*model.Order), args.Error(1)
	if err != nil {
		return nil, httperors.NewNotFoundError("Something went wrong getting the resourse")
	}
	return Product, nil
}
func (mock MockRepository) GetAll() ([]*model.Order, httperors.HttpErr) {
	args := mock.Called()
	result := args.Get(0)
	Products, err := result.([]*model.Order), args.Error(1)
	if err != nil {
		return nil, httperors.NewNotFoundError("Something went wrong getting the resourses")
	}
	return Products, nil
}
func (mock MockRepository) Update(code string, Product *model.Order) httperors.HttpErr {
	args := mock.Called()
	result := args.Get(0)
	_, err := result.(*model.Order), args.Error(1)
	if err != nil {
		return httperors.NewNotFoundError("Something went wrong updating the resourse")
	}
	return nil
}
func TestGetAll(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Create").Return(Order, nil)
	_, _ = mockRepo.Create(Order)
	mockRepo.On("GetAll").Return([]*model.Order{Order}, nil)
	results, _ := mockRepo.GetAll()
	//mock assertion: behavioral
	mockRepo.AssertExpectations(t)
	//data assertion
	assert.Equal(t, Order.Customer, results[0].Customer)
	assert.Equal(t, Order.Address, results[0].Address)
	assert.Equal(t, Order.Phone, results[0].Phone)
	assert.Equal(t, Order.Usercode, results[0].Usercode)
	assert.Equal(t, Order.Amount, results[0].Amount)

}
func TestGetOne(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Create").Return(Order, nil)
	_, _ = mockRepo.Create(Order)
	mockRepo.On("GetOne").Return(Order, nil)
	results, _ := mockRepo.GetOne()
	//mock assertion: behavioral
	mockRepo.AssertExpectations(t)
	//data assertion
	assert.Equal(t, Order.Customer, results.Customer)
	assert.Equal(t, Order.Address, results.Address)
	assert.Equal(t, Order.Amount, results.Amount)

}
func TestCreate(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Create").Return(Order, nil)
	results, err := mockRepo.Create(Order)
	mockRepo.AssertExpectations(t)
	assert.Equal(t, Order.Customer, results.Customer)
	assert.Equal(t, Order.Address, results.Address)
	assert.Equal(t, Order.Amount, results.Amount)
	assert.Nil(t, err)

}
