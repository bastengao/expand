# expand

Inspired by https://stripe.com/docs/expand . Resolve N+1 query for HTTP API with gorm, preload data what client needs.

## Usage

```go
import (
	"gorm.io/gorm"
	expand "github.com/bastengao/expand/gormadapter"
)

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

// array whitelist
var whitelistArray = []string["CreditCards", "Addresses"]

var users []User
preloads, err := expand.Expand([]string{"CreditCards"}, whitelistArray)
db.Scopes(preloads).Find(&users)


// map whitelist
var whitelistMap = map[string]interface{}{
	"CreditCards": nil,
	"Addresses": nil,
}
preloads, err = expand.Expand([]string{"Addresses"}, whitelistArray)
db.Scopes(preloads).Find(&users)

// nested whitelist
var whitelistNested = map[string]interface{}{
	"CreditCards": []string{"Address"},
	"Addresses": nil,
}
preloads, err = expand.Expand([]string{"CreditCards.Address"}, whitelistArray)
db.Scopes(preloads).Find(&users)
```

Client pass `expand` what data need to preload to backend, then preload data with `expand` by gorm.

Or use `Expander`

```go
expander := expand.New([]string["CreditCards", "Addresses"])

var users []User
preloads, err := expander.Expand([]string{"CreditCards"})
db.Scopes(preloads).Find(&users)
```
