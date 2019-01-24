package templates

import (
	"fmt"
	"strings"
	"time"

	"github.com/yanzay/picasso/models"
	"github.com/yanzay/tbot"
)

const extended = 23

// Main screen
func RenderPlayer(m *tbot.Message, p *models.Player) {
	s := newScreen(p.FullName())
	s.add("Prestige", p.Prestige, PrestigeEmoji)
	s.addLine()
	s.add("Energy", p.Energy, BatteryEmoji)
	s.add("Battle", p.BattleEnergy, FrameEmoji)
	s.add("Defence", p.DefenceEnergy, ShieldEmoji)
	s.add("Attack", p.AttackEnergy, WeaponEmoji)
	s.addLine()
	s.add("Coins", p.Coins, CoinEmoji)
	s.add("Nanobots", p.Resource1, MicroscopeEmoji)
	s.add("Parts", p.Resource2, BoltEmoji)
	s.add("Food", p.Resource3, FoodEmoji)
	m.ReplyKeyboard(wrapMarkdown(s.render()), mainMenuButtons, tbot.WithMarkdown)
}

// Implants
func RenderImplants(m *tbot.Message, p *models.Player) {
	im := p.Implants
	s := newScreen(ImplantsButton)
	batteryLevel := padN(BatteryEmoji, fmt.Sprintf("%d%s", im.Battery, renderAvailable(p.UpgradeAvailable(im.Battery, models.Battery))), 7)
	exoframeLevel := padN(FrameEmoji, fmt.Sprintf("%d%s", im.Exoframe, renderAvailable(p.UpgradeAvailable(im.Exoframe, models.Exoframe))), 7)
	shieldLevel := padN(ShieldEmoji, fmt.Sprintf("%d%s", im.Shield, renderAvailable(p.UpgradeAvailable(im.Shield, models.Shield))), 7)
	weaponLevel := padN(WeaponEmoji, fmt.Sprintf("%d%s", im.Weapon, renderAvailable(p.UpgradeAvailable(im.Weapon, models.Weapon))), 7)
	s.add(batteryLevel, fmt.Sprintf("%d/%d", p.Energy, im.BatteryCapacity()), PowerEmoji)
	s.add(exoframeLevel, fmt.Sprintf("%d/%d", p.BattleEnergy, im.ExoframeCapacity()), PowerEmoji)
	s.addLine()
	s.add(shieldLevel, fmt.Sprintf("%d/%d", p.DefenceEnergy, im.ShieldCapacity()), PowerEmoji)
	s.add(weaponLevel, fmt.Sprintf("%d/%d", p.AttackEnergy, im.WeaponCapacity()), PowerEmoji)
	m.ReplyKeyboard(wrapMarkdown(s.render()), implantsMenuButtons, tbot.WithMarkdown)
}

func RenderBattery(m *tbot.Message, p *models.Player) {
	s := newScreen(BatteryButton)
	s.add("Level", p.Implants.Battery, "")
	s.add("Energy", fmt.Sprintf("%d/%d", p.Energy, p.Implants.BatteryCapacity()), PowerEmoji)
	s.add("Food", fmt.Sprintf("%+d", p.Equipment.ResourceMiner3Speed()-p.ConsumingSpeed()), FoodEmoji+"/min")
	s.addLine()
	up := models.GetUpgrade(p.Implants.Battery, models.Battery)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), batteryMenuButtons, tbot.WithMarkdown)
}

func RenderExoframe(m *tbot.Message, p *models.Player) {
	s := newScreen(ExoframeButton)
	s.add("Level", p.Implants.Exoframe, "")
	s.add("Power", fmt.Sprintf("%d/%d", p.BattleEnergy, p.Implants.ExoframeCapacity()), PowerEmoji)
	s.addLine()
	s.add("Energy", p.Energy, PowerEmoji)
	s.addLine()
	up := models.GetUpgrade(p.Implants.Exoframe, models.Exoframe)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), exoframeMenuButtons, tbot.WithMarkdown)
}

func RenderShield(m *tbot.Message, p *models.Player) {
	s := newScreen(ShieldButton)
	s.add("Level", p.Implants.Shield, "")
	s.add("Power", fmt.Sprintf("%d/%d", p.DefenceEnergy, p.Implants.ShieldCapacity()), PowerEmoji)
	s.addLine()
	s.add("Energy", p.Energy, PowerEmoji)
	s.addLine()
	up := models.GetUpgrade(p.Implants.Shield, models.Shield)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), shieldMenuButtons, tbot.WithMarkdown)
}

func RenderWeapon(m *tbot.Message, p *models.Player) {
	s := newScreen(WeaponButton)
	s.add("Level", p.Implants.Weapon, "")
	s.add("Power", fmt.Sprintf("%d/%d", p.AttackEnergy, p.Implants.WeaponCapacity()), PowerEmoji)
	s.addLine()
	s.add("Energy", p.Energy, PowerEmoji)
	s.addLine()
	up := models.GetUpgrade(p.Implants.Weapon, models.Weapon)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), weaponMenuButtons, tbot.WithMarkdown)
}

// Equipment
func RenderEquipment(m *tbot.Message, p *models.Player) {
	eq := p.Equipment
	s := newScreen(EquipmentButton)
	s.add(CoinminerButton, eq.CoinMiner, renderAvailable(p.UpgradeAvailable(eq.CoinMiner, models.CoinMiner)))
	s.add(BackpackButton, eq.Backpack, renderAvailable(p.UpgradeAvailable(eq.Backpack, models.Backpack)))
	s.add(NanofactoryButton, eq.ResourceMiner1, renderAvailable(p.UpgradeAvailable(eq.ResourceMiner1, models.ResourceMiner1)))
	s.add(PartsfactoryButton, eq.ResourceMiner2, renderAvailable(p.UpgradeAvailable(eq.ResourceMiner2, models.ResourceMiner2)))
	s.add(FoodfactoryButton, eq.ResourceMiner3, renderAvailable(p.UpgradeAvailable(eq.ResourceMiner3, models.ResourceMiner3)))
	m.ReplyKeyboard(wrapMarkdown(s.render()), equipmentMenuButtons, tbot.WithMarkdown)
}

func RenderCoinMiner(m *tbot.Message, p *models.Player) {
	s := newScreen(CoinminerButton)
	s.add("Level", p.Equipment.CoinMiner, "")
	s.add("Coins", p.Coins, CoinEmoji)
	s.add("Max", p.Equipment.CoinMinerCapacity(), CoinEmoji)
	s.add("Energy", p.Energy, PowerEmoji)
	s.add("Speed", p.Equipment.CoinMinerSpeed(p.Energy), CoinEmoji+"/min")
	s.addLine()
	up := models.GetUpgrade(p.Equipment.CoinMiner, models.CoinMiner)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), coinminerMenuButtons, tbot.WithMarkdown)
}

func RenderBackpack(m *tbot.Message, p *models.Player) {
	cap := p.Equipment.BackpackCapacity()
	s := newScreen(BackpackButton)
	s.add("Level", p.Equipment.Backpack, "")
	s.addLine()
	s.add("Resources", "", "")
	s.add("", fmt.Sprintf("%d/%d", p.Resource1, cap), MicroscopeEmoji)
	s.add("", fmt.Sprintf("%d/%d", p.Resource2, cap), BoltEmoji)
	s.add("", fmt.Sprintf("%d/%d", p.Resource3, cap), FoodEmoji)
	s.addLine()
	s.add("Coins", p.Coins, CoinEmoji)
	s.addLine()
	s.add("Fill", p.FillBackpackPrice(), CoinEmoji)
	s.addLine()
	up := models.GetUpgrade(p.Equipment.Backpack, models.Backpack)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), backpackMenuButtons, tbot.WithMarkdown)
}

func RenderNanofactory(m *tbot.Message, p *models.Player) {
	s := newScreen(NanofactoryButton)
	s.add("Level", p.Equipment.ResourceMiner1, "")
	s.add("Production", p.Equipment.ResourceMiner1Speed(), MicroscopeEmoji+"/min")
	s.addLine()
	s.add("Energy", p.Energy, PowerEmoji)
	s.add("Coins", p.Coins, CoinEmoji)
	s.addLine()
	up := models.GetUpgrade(p.Equipment.ResourceMiner1, models.ResourceMiner1)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), resourceMinerMenuButtons, tbot.WithMarkdown)
}

func RenderPartsfactory(m *tbot.Message, p *models.Player) {
	s := newScreen(PartsfactoryButton)
	s.add("Level", p.Equipment.ResourceMiner2, "")
	s.add("Production", p.Equipment.ResourceMiner2Speed(), BoltEmoji+"/min")
	s.addLine()
	s.add("Energy", p.Energy, PowerEmoji)
	s.add("Coins", p.Coins, CoinEmoji)
	s.addLine()
	up := models.GetUpgrade(p.Equipment.ResourceMiner2, models.ResourceMiner2)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), resourceMinerMenuButtons, tbot.WithMarkdown)
}

func RenderFoodfactory(m *tbot.Message, p *models.Player) {
	s := newScreen(FoodfactoryButton)
	s.add("Level", p.Equipment.ResourceMiner3, "")
	s.add("Production", p.Equipment.ResourceMiner3Speed(), FoodEmoji+"/min")
	s.add("Food", fmt.Sprintf("%+d", p.Equipment.ResourceMiner3Speed()-p.ConsumingSpeed()), FoodEmoji+"/min")
	s.addLine()
	s.add("Energy", p.Energy, PowerEmoji)
	s.add("Coins", p.Coins, CoinEmoji)
	s.addLine()
	up := models.GetUpgrade(p.Equipment.ResourceMiner3, models.ResourceMiner3)
	renderUpgrade(s, p, up)
	m.ReplyKeyboard(wrapMarkdown(s.render()), resourceMinerMenuButtons, tbot.WithMarkdown)
}

// Shop
func RenderShop(m *tbot.Message, p *models.Player) {
	s := newScreen(ShopButton)
	s.add("Resources", "", "")
	s.add("Coins", p.Coins, CoinEmoji)
	s.add("Nanobots", p.Resource1, MicroscopeEmoji)
	s.add("Parts", p.Resource2, BoltEmoji)
	s.add("Food", p.Resource3, FoodEmoji)
	m.ReplyKeyboard(wrapMarkdown(s.render()), shopMenuButtons, tbot.WithMarkdown)
}

func RenderShopMenu(m *tbot.Message) {
	reply := fmt.Sprintf("Price for each item - 2%s\nWhat do you want to buy?", CoinEmoji)
	m.ReplyKeyboard(reply, shopBuyMenuButtons)
}

func RenderChooseAmount(m *tbot.Message, emoji string) {
	m.ReplyKeyboard("Choose amount:", amountButtons(emoji))
}

func RenderBuyState(m *tbot.Message, p *models.Player, footer string) {
	s := newScreen("")
	s.add("Coins", p.Coins, CoinEmoji)
	s.add("Nanobots", p.Resource1, MicroscopeEmoji)
	s.add("Parts", p.Resource2, BoltEmoji)
	s.add("Food", p.Resource3, FoodEmoji)
	s.footer = footer
	m.Reply(wrapMarkdown(s.render()), tbot.WithMarkdown)
}

func RenderBattle(m *tbot.Message, p *models.Player, b *models.Battle) {
	im := p.Implants
	s := newScreen(BattleButton)
	s.add("Wins", p.Wins, WinEmoji)
	s.add("Prestige", p.Prestige, PrestigeEmoji)
	s.addLine()
	s.add("Battle", fmt.Sprintf("%d/%d", p.BattleEnergy, im.ExoframeCapacity()), PowerEmoji)
	s.add("Defence", fmt.Sprintf("%d/%d", p.DefenceEnergy, im.ShieldCapacity()), PowerEmoji)
	s.add("Attack", fmt.Sprintf("%d/%d", p.AttackEnergy, im.WeaponCapacity()), PowerEmoji)
	var battle string
	if b != nil {
		attackers := []string{}
		for _, attacker := range b.Attack {
			attackers = append(attackers, fmt.Sprintf("`%s`", attacker.FullName()))
		}
		defenders := []string{}
		for _, defender := range b.Defence {
			defenders = append(defenders, fmt.Sprintf("`%s`", defender.FullName()))
		}
		attack := fmt.Sprintf("Attack: %s", strings.Join(attackers, ", "))
		defence := fmt.Sprintf("Defence: %s", strings.Join(defenders, ", "))
		battle = fmt.Sprintf("\n\nBattle status\n%s\n%s", attack, defence)
	}
	status := p.Status()
	if status != nil {
		s.addLine()
		s.add("Status:", status.Title, "")
		s.add("Time left:", roundDuration(status.DurationLeft), "")
	}
	m.ReplyKeyboard(wrapMarkdown(s.render())+battle, battleMenuButtons, tbot.WithMarkdown)
}

func RenderSearchResult(m *tbot.Message, p *models.Player) {
	m.Replyf("You found %s with prestige %d%s", p.Enemy.FullName(), p.Enemy.Prestige, PrestigeEmoji)
}

func RenderEmptySearch(m *tbot.Message, p *models.Player) {
	m.Reply("Sorry, but you can't find any suitable competitor.")
}

func RenderNoTarget(m *tbot.Message) {
	m.Reply("You have no target to attack")
}

func RenderAlreadyInBattle(m *tbot.Message) {
	m.Reply("Player is already in battle, you can't attack now")
}

func RenderHasImmune(m *tbot.Message) {
	m.Reply("Player was attacked recently, you can't attack now")
}

func RenderGuilds(m *tbot.Message) {
	s := newScreen(GuildButton)
	s.add("Choose your destiny:", "", "")
	for _, guild := range models.Guilds {
		s.add(guild.Title(), "", "")
	}
	m.ReplyKeyboard(wrapMarkdown(s.render()), guildMenuButtons(), tbot.WithMarkdown)
}

func RenderGuild(m *tbot.Message, stat *models.GuildStat) {
	s := newScreen(stat.Title)
	s.add("Leader:", "", "")
	s.add("", stat.Leader.ShortName(), PrestigeEmoji)
	s.addLine()
	s.add("Players", stat.PlayersNumber, PlayerEmoji)
	s.addLine()
	s.add("Prestige", stat.Prestige, PrestigeEmoji)
	s.add("Coins", stat.Coins, CoinEmoji)
	s.add("Wins", stat.Wins, WinEmoji)
	m.ReplyKeyboard(wrapMarkdown(s.render()), guildViewMenuButtons, tbot.WithMarkdown)
}

func RenderTop(m *tbot.Message) {
	m.ReplyKeyboard(TopButton, topMenuButtons)
}

func RenderTopPlayersPrestige(m *tbot.Message, players []*models.Player) {
	s := newScreen(PrestigeButton)
	for _, player := range players {
		s.add(player.ShortName(), player.Prestige, PrestigeEmoji)
	}
	m.ReplyKeyboard(wrapMarkdown(s.renderN(extended)), topPlayersMenuButtons, tbot.WithMarkdown)
}

func RenderTopPlayersCoins(m *tbot.Message, players []*models.Player) {
	s := newScreen(CoinsButton)
	for _, player := range players {
		s.add(player.ShortName(), player.Coins, CoinEmoji)
	}
	m.ReplyKeyboard(wrapMarkdown(s.renderN(extended)), topPlayersMenuButtons, tbot.WithMarkdown)
}

func RenderTopPlayersWins(m *tbot.Message, players []*models.Player) {
	s := newScreen(WinsButton)
	for _, player := range players {
		s.add(player.ShortName(), player.Wins, WinEmoji)
	}
	m.ReplyKeyboard(wrapMarkdown(s.renderN(extended)), topPlayersMenuButtons, tbot.WithMarkdown)
}

func RenderTopGuildsPrestige(m *tbot.Message, guilds []*models.Guild) {
	s := newScreen(PrestigeButton)
	for _, guild := range guilds {
		s.add(guild.Name, guild.Prestige, PrestigeEmoji)
	}
	m.ReplyKeyboard(wrapMarkdown(s.render()), topGuildsMenuButtons, tbot.WithMarkdown)
}

func RenderTopGuildsCoins(m *tbot.Message, guilds []*models.Guild) {
	s := newScreen(CoinsButton)
	for _, guild := range guilds {
		s.add(guild.Name, guild.Coins, CoinEmoji)
	}
	m.ReplyKeyboard(wrapMarkdown(s.render()), topGuildsMenuButtons, tbot.WithMarkdown)
}

func RenderTopGuildsWins(m *tbot.Message, guilds []*models.Guild) {
	s := newScreen(WinsButton)
	for _, guild := range guilds {
		s.add(guild.Name, guild.Wins, WinEmoji)
	}
	m.ReplyKeyboard(wrapMarkdown(s.render()), topGuildsMenuButtons, tbot.WithMarkdown)
}

func RenderHelp(m *tbot.Message, supportChatURL, wikiURL string) {
	m.Replyf("%s\n\nSupport chat: %s\nWiki: %s", HelpButton, supportChatURL, wikiURL)
}

// Util
func RenderNotEnoughResources(m *tbot.Message, up *models.Upgrade) {
	resources := []string{}
	if up.Coins > 0 {
		resources = append(resources, fmt.Sprintf("%d%s", up.Coins, CoinEmoji))
	}
	if up.Resource1 > 0 {
		resources = append(resources, fmt.Sprintf("%d%s", up.Resource1, MicroscopeEmoji))
	}
	if up.Resource2 > 0 {
		resources = append(resources, fmt.Sprintf("%d%s", up.Resource2, BoltEmoji))
	}
	m.Replyf("Not enough resources: %s", strings.Join(resources, ", "))
}

func renderUpgrade(s *screen, p *models.Player, up *models.Upgrade) {
	s.add("Upgrade", "", "")
	s.add("", up.Coins, CoinEmoji+available(p.Coins, up.Coins))
	s.add("", up.Resource1, MicroscopeEmoji+available(p.Resource1, up.Resource1))
	s.add("", up.Resource2, BoltEmoji+available(p.Resource2, up.Resource2))
}

func available(got int, need int) string {
	return renderAvailable(got >= need)
}

func renderAvailable(available bool) string {
	if available {
		return YesEmoji
	}
	return NoEmoji
}

func wrapMarkdown(msg string) string {
	return fmt.Sprintf("```\n%s```", msg)
}

func roundDuration(d time.Duration) time.Duration {
	return d - (d % time.Second)
}

func playerList(players []*models.Player) string {
	playerNames := make([]string, 0, len(players))
	for _, player := range players {
		playerNames = append(playerNames, player.FullName())
	}
	return strings.Join(playerNames, ", ")
}
