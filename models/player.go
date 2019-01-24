package models

import (
	"fmt"
	"time"
)

const (
	ImmuneStatusInClinic   = "In clinic"
	ImmuneStatusInHospital = "In hospital"
)

// Player main model
type Player struct {
	ID        int
	ChatID    int64
	FirstName string
	LastName  string
	Blocked   bool
	Starting  bool

	Guild          string
	GuildID        int
	Wins           int
	BattleID       int
	Prestige       int
	LastDefence    time.Time
	LastDefenceWin bool

	Energy        int
	BattleEnergy  int
	DefenceEnergy int
	AttackEnergy  int
	Coins         int
	Resource1     int
	Resource2     int
	Resource3     int

	Implants  Implants
	Equipment Equipment
	Enemy     *Player

	Bot bool
}

type ImmuneStatus struct {
	Title        string
	DurationLeft time.Duration
}

func NewPlayer(id int) *Player {
	return &Player{
		ID:        id,
		Coins:     50000,
		Resource1: 1000,
		Resource2: 1000,
		Resource3: 1000,
		Energy:    10,
		Equipment: NewEquipment(),
		Implants:  NewImplants(),
	}
}

// IsFull indicates if player's profile has all fields filled up
func (p *Player) IsFull() bool {
	return p.FirstName != "" && p.LastName != ""
}

// FullName returs first and last name of player with guild
func (p *Player) FullName() string {
	var name string
	if p.Guild != "" {
		name = fmt.Sprintf("[%s]", p.Guild)
	}
	return name + fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}

func (p *Player) ShortName() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}

func (p *Player) AddEnergy() {
	if p.Resource3 == 0 {
		return
	}
	p.Energy += p.Implants.Battery
	cap := p.Implants.BatteryCapacity()
	if p.Energy > cap {
		p.AddBattleEnergy(p.Energy - cap)
		p.Energy = cap
	}
}

func (p *Player) AddBattleEnergy(amount int) {
	p.BattleEnergy += amount
	cap := p.Implants.ExoframeCapacity()
	if p.BattleEnergy > cap {
		p.AddDefenceEnergy(p.BattleEnergy - cap)
		p.BattleEnergy = cap
	}
}

func (p *Player) AddDefenceEnergy(amount int) {
	p.DefenceEnergy += amount
	cap := p.Implants.ShieldCapacity()
	if p.DefenceEnergy > cap {
		p.AddAttackEnergy(p.DefenceEnergy - cap)
		p.DefenceEnergy = cap
	}
}

func (p *Player) AddAttackEnergy(amount int) {
	p.AttackEnergy += amount
	cap := p.Implants.WeaponCapacity()
	if p.AttackEnergy > cap {
		p.AttackEnergy = cap
	}
}

func (p *Player) AddCoins() {
	diff := p.Equipment.CoinMinerSpeed(p.Energy)
	cap := p.Equipment.CoinMinerCapacity()
	p.Coins += diff
	if p.Coins > cap {
		p.Coins = cap
	}
}

func (p *Player) ProduceResource1() {
	diff := p.Equipment.ResourceMiner1Speed()
	cap := p.Equipment.BackpackCapacity()
	p.Resource1 += diff
	if p.Resource1 > cap {
		p.Resource1 = cap
	}
}

func (p *Player) ProduceResource2() {
	diff := p.Equipment.ResourceMiner2Speed()
	cap := p.Equipment.BackpackCapacity()
	p.Resource2 += diff
	if p.Resource2 > cap {
		p.Resource2 = cap
	}
}

func (p *Player) ProduceResource3() {
	diff := p.Equipment.ResourceMiner3Speed() - p.ConsumingSpeed()
	cap := p.Equipment.BackpackCapacity()
	p.Resource3 += diff
	if p.Resource3 > cap {
		p.Resource3 = cap
	}
	if p.Resource3 < 0 {
		p.Resource3 = 0
	}
}

func (p *Player) ConsumingSpeed() int {
	return p.Energy / 2
}

func (p *Player) UpgradeCoinMiner() error {
	up := GetUpgrade(p.Equipment.CoinMiner, CoinMiner)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Equipment.CoinMiner++
	return nil
}

func (p *Player) UpgradeBackpack() error {
	up := GetUpgrade(p.Equipment.Backpack, Backpack)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Equipment.Backpack++
	return nil
}

func (p *Player) UpgradeResourceMiner1() error {
	up := GetUpgrade(p.Equipment.ResourceMiner1, ResourceMiner1)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Equipment.ResourceMiner1++
	return nil
}

func (p *Player) UpgradeResourceMiner2() error {
	up := GetUpgrade(p.Equipment.ResourceMiner2, ResourceMiner2)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Equipment.ResourceMiner2++
	return nil
}

func (p *Player) UpgradeResourceMiner3() error {
	up := GetUpgrade(p.Equipment.ResourceMiner3, ResourceMiner3)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Equipment.ResourceMiner3++
	return nil
}

func (p *Player) UpgradeBattery() error {
	up := GetUpgrade(p.Implants.Battery, Battery)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Implants.Battery++
	return nil
}

func (p *Player) UpgradeExoframe() error {
	up := GetUpgrade(p.Implants.Exoframe, Exoframe)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Implants.Exoframe++
	return nil
}

func (p *Player) UpgradeShield() error {
	up := GetUpgrade(p.Implants.Shield, Shield)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Implants.Shield++
	return nil
}

func (p *Player) UpgradeWeapon() error {
	up := GetUpgrade(p.Implants.Weapon, Weapon)
	err := p.doUpgrade(up)
	if err != nil {
		return err
	}
	p.Implants.Weapon++
	return nil
}

func (p *Player) doUpgrade(up *Upgrade) error {
	diffCoins := up.Coins - p.Coins
	diffResource1 := up.Resource1 - p.Resource1
	diffResource2 := up.Resource2 - p.Resource2
	if diffCoins > 0 || diffResource1 > 0 || diffResource2 > 0 {
		return &Upgrade{
			Coins:     diffCoins,
			Resource1: diffResource1,
			Resource2: diffResource2,
		}
	}
	p.Coins -= up.Coins
	p.Resource1 -= up.Resource1
	p.Resource2 -= up.Resource2
	p.Prestige += 10
	return nil
}

func (p *Player) UpgradeAvailable(level int, item string) bool {
	up := GetUpgrade(level, item)
	if up.Coins > p.Coins {
		return false
	}
	if up.Resource1 > p.Resource1 {
		return false
	}
	if up.Resource2 > p.Resource2 {
		return false
	}
	return true
}

func (p *Player) FillBackpack() error {
	price := p.FillBackpackPrice()
	if p.Coins < price {
		return &Upgrade{
			Coins: price - p.Coins,
		}
	}
	cap := p.Equipment.BackpackCapacity()
	p.Resource1 = cap
	p.Resource2 = cap
	p.Resource3 = cap
	p.Coins -= price
	return nil
}

func (p *Player) FillBackpackPrice() int {
	cap := p.Equipment.BackpackCapacity()
	r1 := (cap - p.Resource1) * PriceBuyResource1
	r2 := (cap - p.Resource2) * PriceBuyResource2
	r3 := (cap - p.Resource3) * PriceBuyResource3
	return r1 + r2 + r3
}

func (p *Player) FillExoframe() error {
	toFull := p.Implants.ExoframeCapacity() - p.BattleEnergy
	available := p.Energy
	transfer := min(toFull, available)
	p.Energy -= transfer
	p.BattleEnergy += transfer
	return nil
}

func (p *Player) FillShield() error {
	toFull := p.Implants.ShieldCapacity() - p.DefenceEnergy
	available := p.Energy
	transfer := min(toFull, available)
	p.Energy -= transfer
	p.DefenceEnergy += transfer
	return nil
}

func (p *Player) FillWeapon() error {
	toFull := p.Implants.WeaponCapacity() - p.AttackEnergy
	available := p.Energy
	transfer := min(toFull, available)
	p.Energy -= transfer
	p.AttackEnergy += transfer
	return nil
}

func (p *Player) BuyResource1(amount int) error {
	price := PriceBuyResource1 * amount
	if p.Coins < price {
		return fmt.Errorf("Not enough money")
	}
	if p.Resource1+amount > p.Equipment.BackpackCapacity() {
		return fmt.Errorf("Not enough space in Backpack")
	}
	p.Coins -= price
	p.Resource1 += amount
	return nil
}

func (p *Player) BuyResource2(amount int) error {
	price := PriceBuyResource2 * amount
	if p.Coins < price {
		return fmt.Errorf("Not enough money")
	}
	if p.Resource2+amount > p.Equipment.BackpackCapacity() {
		return fmt.Errorf("Not enough space in Backpack")
	}
	p.Coins -= price
	p.Resource2 += amount
	return nil
}

func (p *Player) BuyResource3(amount int) error {
	price := PriceBuyResource3 * amount
	if p.Coins < price {
		return fmt.Errorf("Not enough money")
	}
	if p.Resource3+amount > p.Equipment.BackpackCapacity() {
		return fmt.Errorf("Not enough space in Backpack")
	}
	p.Coins -= price
	p.Resource3 += amount
	return nil
}

func (p *Player) ApplyBattleResult(result *BattleResult) {
	p.BattleID = 0
	p.Enemy = nil
	if result.Attack {
		p.AttackEnergy = 0
	} else {
		p.DefenceEnergy = 0
		p.LastDefence = time.Now()
		p.LastDefenceWin = result.Win
	}
	if result.Win {
		p.Wins++
		p.Coins += result.Coins
		p.Prestige += result.Prestige
	} else {
		p.Coins -= result.Coins
		p.Prestige -= result.Prestige
		p.BattleEnergy = 0
	}
}

func (p *Player) NoGuild() bool {
	return p.Guild == ""
}

func (p *Player) InBattle() bool {
	return p.BattleID != 0
}

func (p *Player) Status() *ImmuneStatus {
	if p.LastDefenceWin {
		if time.Since(p.LastDefence) < ClinicDuration {
			return &ImmuneStatus{
				Title:        ImmuneStatusInClinic,
				DurationLeft: ClinicDuration - time.Since(p.LastDefence),
			}
		}
	} else {
		if time.Since(p.LastDefence) < HospitalDuration {
			return &ImmuneStatus{
				Title:        ImmuneStatusInHospital,
				DurationLeft: HospitalDuration - time.Since(p.LastDefence),
			}
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
