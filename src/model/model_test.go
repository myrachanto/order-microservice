package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateOrderInputRequiredFields(t *testing.T) {
	jsondata := `{"Customer":"Customer Name","Address":"lake view","Amount":10}`
	order := &Order{}
	if err := json.Unmarshal([]byte(jsondata), &order); err != nil {
		t.Errorf("failed to unmarshal product data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	expected := ""
	if err := order.Validate(); err != nil {
		fmt.Println("------------------", err.Message())
		expected = "Invalid Name"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Customer name")
		}
		expected = "Invalid Address"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Address")
		}
		expected = "Invalid Amount"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Amount")
		}
	}

}
func TestValidateItemInputRequiredFields(t *testing.T) {
	jsondata := `{"Name":"Product Name","Price": 20.0,"Quantity":50}`
	item := &Item{}
	if err := json.Unmarshal([]byte(jsondata), &item); err != nil {
		t.Errorf("failed to unmarshal product data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	expected := ""
	if err := item.Validate(); err != nil {
		fmt.Println("------------------", err.Message())
		expected = "Invalid Name"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Product name")
		}
		expected = "Invalid Price"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Price")
		}
		expected = "Invalid Quantity"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Quantity")
		}
	}

}
