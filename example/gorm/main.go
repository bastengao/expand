package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/bastengao/expand/gormadapter"
)

// define expandable fields
var whitelist = map[string]interface{}{
	"CreditCards": []string{"Address"},
	"Addresses":   nil,
}

func main() {
	db := initDB()
	seedData(db)

	scopes, err := gormadapter.Expand([]string{"CreditCards.Address", "Addresses"}, whitelist)
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	err = db.Scopes(scopes).Find(&users).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Print(users)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&User{}, &CreditCard{}, &Address{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func seedData(db *gorm.DB) {
	for _, model := range []interface{}{&User{}, &CreditCard{}, &Address{}} {
		err := db.Unscoped().Where("1 = 1").Delete(model).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	db.Create([]*User{
		{
			Name: "John",
			CreditCards: []CreditCard{
				{
					Number: "123",
					Address: Address{
						Street: "Short street",
					},
				},
			},
			Addresses: []Address{
				{
					Street: "Long street",
				},
			},
		},
	})
}
