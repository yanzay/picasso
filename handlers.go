package main

import (
	"fmt"
	"os"
	"time"

	"github.com/yanzay/log"
	"github.com/yanzay/tbot"

	"github.com/yanzay/picasso/game"
	"github.com/yanzay/picasso/models"
	"github.com/yanzay/picasso/storage"
	"github.com/yanzay/picasso/templates"
)

type handlers struct {
	store storage.Storage
	mids  *middlewares
	game  *game.Game
}

func setHandlers(bot *tbot.Server, mids *middlewares, store storage.Storage, game *game.Game) {
	h := handlers{store: store, game: game, mids: mids}

	// Base routes
	bot.SetAlias(tbot.RouteRefresh, templates.InfoButton)
	bot.SetAlias(tbot.RouteBack, templates.BackButton)
	bot.HandleFunc(tbot.RouteRoot, h.homeHandler)
	bot.SetAlias(tbot.RouteRoot, templates.HomeButton)

	// Common aliases
	bot.SetAlias("upgrade", templates.UpgradeButton)
	bot.SetAlias("fill", templates.FillButton)

	// Implants
	bot.HandleFunc("/implants", h.implantsHandler)
	bot.SetAlias("implants", templates.ImplantsButton)

	// Battery
	bot.HandleFunc("/implants/battery", h.batteryHandler)
	bot.SetAlias("battery", templates.BatteryButton)
	bot.HandleFunc("/implants/battery/upgrade", h.mids.inBattle(h.batteryUpgradeHandler))

	// Exoframe
	bot.HandleFunc("/implants/exoframe", h.exoframeHandler)
	bot.SetAlias("exoframe", templates.ExoframeButton)
	bot.HandleFunc("/implants/exoframe/upgrade", h.mids.inBattle(h.exoframeUpgradeHandler))
	bot.HandleFunc("/implants/exoframe/fill", h.exoframeFillHandler)

	// Shield
	bot.HandleFunc("/implants/shield", h.shieldHandler)
	bot.SetAlias("shield", templates.ShieldButton)
	bot.HandleFunc("/implants/shield/upgrade", h.mids.inBattle(h.shieldUpgradeHandler))
	bot.HandleFunc("/implants/shield/fill", h.shieldFillHandler)

	// Weapon
	bot.HandleFunc("/implants/weapon", h.weaponHandler)
	bot.SetAlias("weapon", templates.WeaponButton)
	bot.HandleFunc("/implants/weapon/upgrade", h.mids.inBattle(h.weaponUpgradeHandler))
	bot.HandleFunc("/implants/weapon/fill", h.weaponFillHandler)

	// Equipment
	bot.HandleFunc("/equipment", h.equipmentHandler)
	bot.SetAlias("equipment", templates.EquipmentButton)

	// Coin Miner
	bot.HandleFunc("/equipment/coinminer", h.coinMinerHandler)
	bot.SetAlias("coinminer", templates.CoinminerButton)
	bot.HandleFunc("/equipment/coinminer/upgrade", h.mids.inBattle(h.coinMinerUpgradeHandler))

	// Backpack
	bot.HandleFunc("/equipment/backpack", h.backpackHandler)
	bot.SetAlias("backpack", templates.BackpackButton)
	bot.HandleFunc("/equipment/backpack/upgrade", h.mids.inBattle(h.backpackUpgradeHandler))
	bot.HandleFunc("/equipment/backpack/fill", h.mids.inBattle(h.backpackFillHandler))

	// Nano Factory
	bot.HandleFunc("/equipment/nanofactory", h.nanofactoryHandler)
	bot.SetAlias("nanofactory", templates.NanofactoryButton)
	bot.HandleFunc("/equipment/nanofactory/upgrade", h.mids.inBattle(h.nanofactoryUpgradeHandler))

	// Parts Factory
	bot.HandleFunc("/equipment/partsfactory", h.partsfactoryHandler)
	bot.SetAlias("partsfactory", templates.PartsfactoryButton)
	bot.HandleFunc("/equipment/partsfactory/upgrade", h.mids.inBattle(h.partsfactoryUpgradeHandler))

	// Food Factory
	bot.HandleFunc("/equipment/foodfactory", h.foodfactoryHandler)
	bot.SetAlias("foodfactory", templates.FoodfactoryButton)
	bot.HandleFunc("/equipment/foodfactory/upgrade", h.mids.inBattle(h.foodfactoryUpgradeHandler))

	// Shop
	bot.HandleFunc("/shop", h.shopHandler)
	bot.SetAlias("shop", templates.ShopButton)
	bot.HandleFunc("/shop/buy", h.shopBuyHandler)
	bot.SetAlias("buy", templates.BuyButton)

	bot.HandleFunc("/shop/buy/nanobots", h.shopBuyNanobotsHandler)
	bot.SetAlias("nanobots", templates.NanobotsButton)
	bot.HandleFunc("/shop/buy/nanobots/{amount}", h.mids.inBattle(h.shopBuyNanobotsAmountHandler))

	bot.HandleFunc("/shop/buy/parts", h.shopBuyPartsHandler)
	bot.SetAlias("parts", templates.PartsButton)
	bot.HandleFunc("/shop/buy/parts/{amount}", h.mids.inBattle(h.shopBuyPartsAmountHandler))

	bot.HandleFunc("/shop/buy/food", h.shopBuyFoodHandler)
	bot.SetAlias("food", templates.FoodButton)
	bot.HandleFunc("/shop/buy/food/{amount}", h.mids.inBattle(h.shopBuyFoodAmountHandler))

	// Battle
	bot.HandleFunc("/battle", h.battleHandler)
	bot.SetAlias("battle", templates.BattleButton)

	bot.HandleFunc("/battle/search", h.mids.inBattle(h.battleSearchHandler))
	bot.SetAlias("search", templates.SearchButton)

	bot.HandleFunc("/battle/searchguild", h.mids.inBattle(h.battleSearchGuildHandler))
	bot.SetAlias("searchguild", templates.SearchGuildButton)

	bot.HandleFunc("/battle/attack", h.mids.inBattle(h.battleAttackHandler))
	bot.SetAlias("attack", templates.AttackButton)

	bot.HandleFunc("/battle/join", h.mids.inBattle(h.battleJoinHandler))
	bot.SetAlias("join", templates.JoinButton)

	// Guild
	bot.HandleFunc("/guild", h.guildHandler)
	bot.SetAlias("guild", templates.GuildButton)

	bot.HandleFunc("/guild/{name}", h.mids.inBattle(h.guildJoinHandler))

	// Top
	bot.HandleFunc("/top", h.topHandler)
	bot.SetAlias("top", templates.TopButton)

	bot.HandleFunc("/top/players", h.topPlayersHandler)
	bot.SetAlias("players", templates.PlayersButton)

	bot.HandleFunc("/top/players/prestige", h.topPlayersPrestigeHandler)
	bot.SetAlias("prestige", templates.PrestigeButton)

	bot.HandleFunc("/top/players/coins", h.topPlayersCoinsHandler)
	bot.SetAlias("coins", templates.CoinsButton)

	bot.HandleFunc("/top/players/wins", h.topPlayersWinsHandler)
	bot.SetAlias("wins", templates.WinsButton)

	bot.HandleFunc("/top/guilds", h.topGuildsHandler)
	bot.SetAlias("guilds", templates.GuildsButton)

	bot.HandleFunc("/top/guilds/prestige", h.topGuildsPrestigeHandler)

	bot.HandleFunc("/top/guilds/coins", h.topGuildsCoinsHandler)

	bot.HandleFunc("/top/guilds/wins", h.topGuildsWinsHandler)

	// Help
	bot.HandleFunc("/gethelp", h.helpHandler)
	bot.SetAlias("gethelp", templates.HelpButton)

	// Start
	bot.HandleFunc("/start", h.mids.inBattle(h.startHandler)) // danger
	bot.SetAlias("start", templates.YesButton, templates.NoButton)

	bot.HandleFunc("/start/yes", h.startYesHandler)
	bot.SetAlias("yes", templates.YesButton)

	bot.HandleFunc("/start/no", h.startNoHandler)
	bot.SetAlias("no", templates.NoButton)

	bot.HandleFunc("/unblock {name}", h.mids.onlyAdmin(h.unblockHandler))
	bot.HandleFunc("/block {name}", h.mids.onlyAdmin(h.blockHandler))

	bot.HandleFunc("/fix", h.mids.onlyAdmin(h.fixHandler))
	bot.HandleFunc("/backup", h.mids.onlyAdmin(h.backupHandler))

	bot.HandleDefault(h.defaultHandler)
}

func (h *handlers) backupHandler(m *tbot.Message) {
	f, err := os.Create(fmt.Sprintf("/data/_backup_%d.db", time.Now().Unix()))
	if err != nil {
		m.Reply(err.Error())
		return
	}
	defer f.Close()
	err = h.store.Backup(f)
	if err != nil {
		m.Reply(err.Error())
		return
	}
	m.Reply("OK")
}

func (h *handlers) fixHandler(m *tbot.Message) {
	m.Reply(h.game.Fix())
}

func (h *handlers) unblockHandler(m *tbot.Message) {
	players, err := h.store.GetAllPlayers()
	if err != nil {
		log.Errorf("can't get players: %v", err)
		return
	}
	for _, player := range players {
		if player.FullName() == m.Vars["name"] {
			player.Blocked = false
			h.store.SetPlayer(player)
		}
	}
	m.Reply("OK")
}

func (h *handlers) blockHandler(m *tbot.Message) {
	players, err := h.store.GetAllPlayers()
	if err != nil {
		log.Errorf("can't get players: %v", err)
		return
	}
	for _, player := range players {
		if player.FullName() == m.Vars["name"] {
			player.Blocked = true
			h.store.SetPlayer(player)
		}
	}
	m.Reply("OK")
}

func (h *handlers) startHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.PlayerSetStarting(player)
	if err != nil {
		log.Errorf("unable to set starting player: %v", err)
		return
	}
	m.ReplyKeyboard("Are you sure?", [][]string{[]string{templates.YesButton, templates.NoButton}}, tbot.OneTimeKeyboard)
}

func (h *handlers) startYesHandler(m *tbot.Message) {
	{
		player, err := h.store.GetPlayer(m.From.ID)
		if err != nil {
			log.Errorf("can't get player %d: %v", m.From.ID, err)
			return
		}
		if player.Starting {
			err = h.game.PlayerStart(player)
			if err != nil {
				log.Errorf("player can't /start %v", err)
			}
		}
	}
	h.store.Reset(m.ChatID)
	h.mids.login(h.homeHandler)(m)
}

func (h *handlers) startNoHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	h.store.Reset(m.ChatID)
	templates.RenderPlayer(m, player)
}

func (h *handlers) homeHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderPlayer(m, player)
}

func (h *handlers) implantsHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderImplants(m, player)
}

func (h *handlers) equipmentHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderEquipment(m, player)
}

func (h *handlers) coinMinerHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderCoinMiner(m, player)
}

func (h *handlers) coinMinerUpgradeHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.UpgradeCoinMiner(player)
	if err != nil {
		up, ok := err.(*models.Upgrade)
		if !ok {
			log.Errorf("unexpected error upgrading coin miner: %v", err)
			return
		}
		templates.RenderNotEnoughResources(m, up)
		return
	}
	templates.RenderCoinMiner(m, player)
}

func (h *handlers) backpackHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderBackpack(m, player)
}

func (h *handlers) backpackUpgradeHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.UpgradeBackpack(player)
	if err != nil {
		up, ok := err.(*models.Upgrade)
		if !ok {
			log.Errorf("unexpected error upgrading backpack: %v", err)
			return
		}
		templates.RenderNotEnoughResources(m, up)
		return
	}
	templates.RenderBackpack(m, player)
}

func (h *handlers) backpackFillHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.FillBackpack(player)
	if err != nil {
		up, ok := err.(*models.Upgrade)
		if !ok {
			log.Errorf("unexpected error filling backpack: %v", err)
			return
		}
		templates.RenderNotEnoughResources(m, up)
		return
	}
	templates.RenderBackpack(m, player)
}

func (h *handlers) nanofactoryHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderNanofactory(m, player)
}

func (h *handlers) nanofactoryUpgradeHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.UpgradeResourceMiner1(player)
	if err != nil {
		up, ok := err.(*models.Upgrade)
		if !ok {
			log.Errorf("unexpected error upgrading resource miner 1: %v", err)
			return
		}
		templates.RenderNotEnoughResources(m, up)
		return
	}
	templates.RenderNanofactory(m, player)
}

func (h *handlers) partsfactoryHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderPartsfactory(m, player)
}

func (h *handlers) partsfactoryUpgradeHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.UpgradeResourceMiner2(player)
	if err != nil {
		up, ok := err.(*models.Upgrade)
		if !ok {
			log.Errorf("unexpected error upgrading resource miner 2: %v", err)
			return
		}
		templates.RenderNotEnoughResources(m, up)
		return
	}
	templates.RenderPartsfactory(m, player)
}

func (h *handlers) foodfactoryHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderFoodfactory(m, player)
}

func (h *handlers) foodfactoryUpgradeHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.UpgradeResourceMiner3(player)
	if err != nil {
		up, ok := err.(*models.Upgrade)
		if !ok {
			log.Errorf("unexpected error upgrading resource miner 3: %v", err)
			return
		}
		templates.RenderNotEnoughResources(m, up)
		return
	}
	templates.RenderFoodfactory(m, player)
}

func (h *handlers) batteryHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	templates.RenderBattery(m, player)
}

func (h *handlers) batteryUpgradeHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.UpgradeBattery(player)
	if err != nil {
		up, ok := err.(*models.Upgrade)
		if !ok {
			log.Errorf("unexpected error upgrading battery: %v", err)
			return
		}
		templates.RenderNotEnoughResources(m, up)
		return
	}
	templates.RenderBattery(m, player)
}

func (h *handlers) exoframeHandler(m *tbot.Message) {
	h.viewHandler(m, templates.RenderExoframe)
}

func (h *handlers) exoframeUpgradeHandler(m *tbot.Message) {
	h.changeHandler(m, h.game.UpgradeExoframe, templates.RenderExoframe)
}

func (h *handlers) exoframeFillHandler(m *tbot.Message) {
	h.changeHandler(m, h.game.FillExoframe, templates.RenderExoframe)
}

func (h *handlers) shieldHandler(m *tbot.Message) {
	h.viewHandler(m, templates.RenderShield)
}

func (h *handlers) shieldUpgradeHandler(m *tbot.Message) {
	h.changeHandler(m, h.game.UpgradeShield, templates.RenderShield)
}

func (h *handlers) shieldFillHandler(m *tbot.Message) {
	h.changeHandler(m, h.game.FillShield, templates.RenderShield)
}

func (h *handlers) weaponHandler(m *tbot.Message) {
	h.viewHandler(m, templates.RenderWeapon)
}

func (h *handlers) weaponUpgradeHandler(m *tbot.Message) {
	h.changeHandler(m, h.game.UpgradeWeapon, templates.RenderWeapon)
}

func (h *handlers) weaponFillHandler(m *tbot.Message) {
	h.changeHandler(m, h.game.FillWeapon, templates.RenderWeapon)
}

func (h *handlers) shopHandler(m *tbot.Message) {
	h.viewHandler(m, templates.RenderShop)
}

func (h *handlers) shopBuyHandler(m *tbot.Message) {
	templates.RenderShopMenu(m)
}

func (h *handlers) shopBuyNanobotsHandler(m *tbot.Message) {
	templates.RenderChooseAmount(m, templates.MicroscopeEmoji)
}

func (h *handlers) shopBuyNanobotsAmountHandler(m *tbot.Message) {
	h.buyHandler(m, h.game.BuyResource1, templates.MicroscopeEmoji)
}

func (h *handlers) shopBuyPartsHandler(m *tbot.Message) {
	templates.RenderChooseAmount(m, templates.BoltEmoji)
}

func (h *handlers) shopBuyPartsAmountHandler(m *tbot.Message) {
	h.buyHandler(m, h.game.BuyResource2, templates.BoltEmoji)
}

func (h *handlers) shopBuyFoodHandler(m *tbot.Message) {
	templates.RenderChooseAmount(m, templates.FoodEmoji)
}

func (h *handlers) shopBuyFoodAmountHandler(m *tbot.Message) {
	h.buyHandler(m, h.game.BuyResource3, templates.FoodEmoji)
}

func (h *handlers) battleHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	var battle *models.Battle
	if player.InBattle() {
		battle, err = h.store.GetBattle(player.BattleID)
		if err != nil {
			log.Errorf("unable to get battle: %v", err)
			return
		}
	}
	templates.RenderBattle(m, player, battle)
}

func (h *handlers) battleSearchHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.FindSingleEnemy(player)
	if err != nil {
		_, ok := err.(game.EmptySearch)
		if !ok {
			log.Errorf("unexpected error: %v", err)
			return
		}
		templates.RenderEmptySearch(m, player)
		return
	}
	templates.RenderSearchResult(m, player)
}

func (h *handlers) battleSearchGuildHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = h.game.FindGuildEnemy(player)
	if err != nil {
		_, ok := err.(game.EmptySearch)
		if !ok {
			log.Errorf("unexpected error: %v", err)
			return
		}
		templates.RenderEmptySearch(m, player)
		return
	}
	templates.RenderSearchResult(m, player)
}

func (h *handlers) battleAttackHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	if player.BattleEnergy == 0 {
		m.Reply("You can't attack without Battle Energy")
		return
	}
	battle, err := h.game.AttackEnemy(player.ID)
	if err != nil {
		_, ok := err.(game.NoTarget)
		if ok {
			templates.RenderNoTarget(m)
			return
		}
		_, ok = err.(game.AlreadyInBattle)
		if ok {
			templates.RenderAlreadyInBattle(m)
			return
		}
		_, ok = err.(game.YouAlreadyInBattle)
		if ok {
			m.Reply("You're already in battle")
			return
		}
		_, ok = err.(game.HasImmune)
		if ok {
			templates.RenderHasImmune(m)
			return
		}
		log.Errorf("unexpected error: %v", err)
		return
	}
	templates.RenderBattle(m, player, battle)
}

func (h *handlers) battleJoinHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	if player.NoGuild() {
		m.Reply("You are not a part of a Guild yet.")
		return
	}
	battle, err := h.game.JoinGuildBattle(player.ID)
	if err != nil {
		if _, ok := err.(game.NoBattle); ok {
			m.Reply("There is no Guild Battle at the time.")
			return
		}
		if _, ok := err.(game.AlreadyInBattle); ok {
			m.Reply("You're already in battle!")
			return
		}
		log.Errorf("unable to join battle: %v", err)
	}
	templates.RenderBattle(m, player, battle)
}

func (h *handlers) guildHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	if player.NoGuild() {
		templates.RenderGuilds(m)
		return
	}
	stat := h.game.GuildStat(player.ID)
	templates.RenderGuild(m, stat)
}

func (h *handlers) guildJoinHandler(m *tbot.Message) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	if !player.NoGuild() {
		m.Reply("hm?")
		return
	}
	if player.Prestige < models.PrestigeToJoinGuild {
		m.Replyf("You need %d%s prestige to join a Guild", models.PrestigeToJoinGuild, templates.PrestigeEmoji)
		return
	}
	name, ok := m.Vars["name"]
	if !ok {
		m.Reply("hm?")
		return
	}
	_, err = h.game.PlayerJoinGuild(player, name)
	if err != nil {
		log.Errorf("player %s unable to join guild: %s", player.FullName(), name)
		m.Reply("hm?")
		return
	}
	stat := h.game.GuildStat(player.ID)
	templates.RenderGuild(m, stat)
}

func (h *handlers) topHandler(m *tbot.Message) {
	templates.RenderTop(m)
}

func (h *handlers) topPlayersHandler(m *tbot.Message) {
	h.topPlayersPrestigeHandler(m)
}

func (h *handlers) topPlayersPrestigeHandler(m *tbot.Message) {
	templates.RenderTopPlayersPrestige(m, h.game.TopPlayersPrestige())
}

func (h *handlers) topPlayersCoinsHandler(m *tbot.Message) {
	templates.RenderTopPlayersCoins(m, h.game.TopPlayersCoins())
}

func (h *handlers) topPlayersWinsHandler(m *tbot.Message) {
	templates.RenderTopPlayersWins(m, h.game.TopPlayersWins())
}

func (h *handlers) topGuildsHandler(m *tbot.Message) {
	h.topGuildsPrestigeHandler(m)
}

func (h *handlers) topGuildsPrestigeHandler(m *tbot.Message) {
	templates.RenderTopGuildsPrestige(m, h.game.TopGuildsPrestige())
}

func (h *handlers) topGuildsCoinsHandler(m *tbot.Message) {
	templates.RenderTopGuildsCoins(m, h.game.TopGuildsCoins())
}

func (h *handlers) topGuildsWinsHandler(m *tbot.Message) {
	templates.RenderTopGuildsWins(m, h.game.TopGuildsWins())
}

func (h *handlers) helpHandler(m *tbot.Message) {
	templates.RenderHelp(m, supportChatURL, wikiURL)
}

func (h *handlers) buyHandler(m *tbot.Message, buy func(p *models.Player, amount int) error, emoji string) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	amountStr, ok := m.Vars["amount"]
	if !ok {
		log.Println("not an amount")
		m.Reply("hm?")
		return
	}
	amount := 0
	_, err = fmt.Sscanf(amountStr, "%d", &amount)
	if err != nil || amount == 0 {
		m.Reply("hm?")
		return
	}
	err = buy(player, amount)
	if err != nil {
		m.Reply(err.Error())
		return
	}
	templates.RenderBuyState(m, player, fmt.Sprintf("Thank you for purchasing %d%s!", amount, emoji))
}

func (h *handlers) changeHandler(m *tbot.Message, change func(p *models.Player) error, render func(m *tbot.Message, player *models.Player)) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	err = change(player)
	if err != nil {
		up, ok := err.(*models.Upgrade)
		if !ok {
			log.Errorf("unexpected error: %v", err)
			return
		}
		templates.RenderNotEnoughResources(m, up)
		return
	}
	render(m, player)
}

func (h *handlers) viewHandler(m *tbot.Message, render func(m *tbot.Message, player *models.Player)) {
	player, err := h.store.GetPlayer(m.From.ID)
	if err != nil {
		log.Errorf("can't get player %d: %v", m.From.ID, err)
		return
	}
	render(m, player)
}

func (h *handlers) defaultHandler(m *tbot.Message) {
	m.Reply("hm?")
}
