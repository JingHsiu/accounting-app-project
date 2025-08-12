package model

import (
	"errors"
	"strings"
)

type CategoryName struct {
	Value string
}

func NewCategoryName(name string) (*CategoryName, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}
	if len(name) > 50 {
		return nil, errors.New("category name cannot exceed 50 characters")
	}
	
	return &CategoryName{Value: name}, nil
}

func (cn CategoryName) String() string {
	return cn.Value
}

func (cn CategoryName) Equals(other CategoryName) bool {
	return cn.Value == other.Value
}