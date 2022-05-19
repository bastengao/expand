package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string
	CreditCards []CreditCard
	Addresses   []Address
}

type CreditCard struct {
	gorm.Model
	Number    string
	UserID    uint
	AddressID uint
	Address   Address
}

type Address struct {
	gorm.Model
	Street string
	UserID *uint
}
