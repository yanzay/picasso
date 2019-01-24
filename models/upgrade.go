package models

import "fmt"

const (
	CoinMiner      = "CoinMiner"
	Backpack       = "Backpack"
	ResourceMiner1 = "ResourceMiner1"
	ResourceMiner2 = "ResourceMiner2"
	ResourceMiner3 = "ResourceMiner3"

	Battery  = "Battery"
	Exoframe = "Exoframe"
	Shield   = "Shield"
	Weapon   = "Weapon"
)

type Coef struct {
	Coins     int
	Resource1 int
	Resource2 int
}

var coefs = map[string]Coef{
	CoinMiner:      Coef{Coins: 500, Resource1: 200, Resource2: 200},
	Backpack:       Coef{Coins: 200, Resource1: 100, Resource2: 100},
	ResourceMiner1: Coef{Coins: 100, Resource1: 50, Resource2: 50},
	ResourceMiner2: Coef{Coins: 100, Resource1: 50, Resource2: 50},
	ResourceMiner3: Coef{Coins: 100, Resource1: 50, Resource2: 50},

	Battery:  Coef{Coins: 200, Resource1: 100, Resource2: 100},
	Exoframe: Coef{Coins: 200, Resource1: 100, Resource2: 100},
	Shield:   Coef{Coins: 5000, Resource1: 500, Resource2: 1500},
	Weapon:   Coef{Coins: 5000, Resource1: 1500, Resource2: 500},
}

type Upgrade struct {
	Coins     int
	Resource1 int
	Resource2 int
}

func (u Upgrade) Error() string {
	return fmt.Sprintf("Not enough resources")
}
