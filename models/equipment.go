package models

type Equipment struct {
	CoinMiner      int
	Backpack       int
	ResourceMiner1 int
	ResourceMiner2 int
	ResourceMiner3 int
}

func NewEquipment() Equipment {
	return Equipment{
		CoinMiner:      1,
		Backpack:       1,
		ResourceMiner1: 1,
		ResourceMiner2: 1,
		ResourceMiner3: 1,
	}
}

func (e *Equipment) CoinMinerCapacity() int {
	return e.CoinMiner * 500000
}

func (e *Equipment) CoinMinerSpeed(energy int) int {
	return (e.CoinMiner*10 + energy/4) * 42
}

func (e *Equipment) ResourceMiner1Speed() int {
	return e.ResourceMiner1 * ResourceMiningSpeed
}

func (e *Equipment) ResourceMiner2Speed() int {
	return e.ResourceMiner2 * ResourceMiningSpeed
}

func (e *Equipment) ResourceMiner3Speed() int {
	return e.ResourceMiner3 * ResourceMiningSpeed
}

func GetUpgrade(level int, item string) *Upgrade {
	k := (level + 1) * (level + 2) / 2
	return &Upgrade{
		Coins:     k * coefs[item].Coins,
		Resource1: k * coefs[item].Resource1,
		Resource2: k * coefs[item].Resource2,
	}
}

func (e *Equipment) BackpackCapacity() int {
	return (e.Backpack*50 + 1000) * e.Backpack
}
