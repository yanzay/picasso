package models

// Implants of player
type Implants struct {
	Battery  int
	Exoframe int
	Shield   int
	Weapon   int
}

func NewImplants() Implants {
	return Implants{
		Battery:  1,
		Exoframe: 1,
		Shield:   1,
		Weapon:   1,
	}
}

func (im Implants) BatteryCapacity() int {
	return im.Battery * 20
}

func (im Implants) ExoframeCapacity() int {
	return im.Exoframe * 40
}

func (im Implants) ShieldCapacity() int {
	return im.Shield * 10
}

func (im Implants) WeaponCapacity() int {
	return im.Weapon * 10
}
