package game

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/yanzay/log"

	"github.com/yanzay/picasso/models"
	"github.com/yanzay/picasso/storage"
)

type BattleRender interface {
	RenderWinResult(*models.Player, *models.BattleResult) string
	RenderLoseResult(*models.Player, *models.BattleResult) string
	RenderAttackedBy(*models.Player) string
	RenderGuildAttack(*models.Guild, *models.Guild) string
	RenderGuildDefence(*models.Guild) string
	RenderJoinedGuildBattle(*models.Player, bool) string
}

type Game struct {
	sync.Mutex
	battleResults chan *models.BattleResult
	sendMessage   func(int, string)
	store         storage.Storage
	battleRender  BattleRender
}

func New(store storage.Storage, sendMessage func(int, string), battleRender BattleRender) *Game {
	return &Game{
		battleResults: make(chan *models.BattleResult),
		sendMessage:   sendMessage,
		store:         store,
		battleRender:  battleRender,
	}
}

func (g *Game) Fix() string {
	g.Lock()
	defer g.Unlock()
	players, err := g.store.GetAllPlayers()
	if err != nil {
		return err.Error()
	}
	members := make(map[int][]int)
	for _, player := range players {
		if player.GuildID > 0 {
			members[player.GuildID] = append(members[player.GuildID], player.ID)
		}
	}
	for guildID, mems := range members {
		guild, err := g.store.GetGuild(guildID)
		if err != nil {
			return err.Error()
		}
		guild.PlayerIDs = mems
		g.store.SetGuild(guild)
	}
	guilds, _ := g.store.GetAllGuilds()
	res := ""
	for _, guild := range guilds {
		res = fmt.Sprintf("%s\n%s: %v", res, guild.Emoji, guild.PlayerIDs)
	}
	return res
}

func (g *Game) Start() {
	g.continueBattles()
	go g.energyLoop()
	go g.battleLoop()
}

func (g *Game) continueBattles() {
	g.Lock()
	defer g.Unlock()
	battles, err := g.store.GetAllBattles()
	if err != nil {
		log.Fatalf("cant continue battles: %q", err)
	}
	for _, battle := range battles {
		go g.startBattle(battle.ID, battle.Start)
	}
}

func (g *Game) PlayerSetStarting(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	player, err := g.store.GetPlayer(p.ID)
	if err != nil {
		return fmt.Errorf("unable to get player %d", p.ID)
	}
	player.Starting = true
	return g.store.SetPlayer(player)
}

func (g *Game) PlayerStart(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	player, err := g.store.GetPlayer(p.ID)
	if err != nil {
		return fmt.Errorf("unable to get player %d", p.ID)
	}
	if player.Starting {
		if player.GuildID > 0 {
			g.removePlayerFromGuild(player.ID, player.GuildID)
		}
		return g.store.DeletePlayer(player.ID)
	}
	return fmt.Errorf("player is not starting! can't proceed")
}

func (g *Game) SendNotification(p *models.Player, result *models.BattleResult) {
	if p.Bot {
		return
	}
	var msg string
	if result.Win {
		msg = g.battleRender.RenderWinResult(p, result)
	} else {
		msg = g.battleRender.RenderLoseResult(p, result)
	}
	g.sendMessage(p.ID, msg)
}

func (g *Game) energyLoop() {
	for range time.Tick(models.GameMinute) {
		start := time.Now()
		g.Lock()
		g.energyForAll()
		g.Unlock()
		log.Infof("Game cycle duration: %s", time.Since(start))
	}
}

func (g *Game) energyForAll() {
	players, err := g.store.GetAllPlayers()
	if err != nil {
		log.Errorf("unable to get all players: %v", err)
		return
	}
	for _, player := range players {
		player.AddCoins()
		player.ProduceResource1()
		player.ProduceResource2()
		player.ProduceResource3()
		if player.InBattle() {
			continue
		}
		player.AddEnergy()
	}
	err = g.store.SetAllPlayers(players)
	if err != nil {
		log.Errorf("unable to set all players: %v", err)
		return
	}
}

func (g *Game) UpgradeCoinMiner(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeCoinMiner()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) UpgradeBackpack(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeBackpack()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) FillBackpack(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.FillBackpack()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) UpgradeResourceMiner1(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeResourceMiner1()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) UpgradeResourceMiner2(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeResourceMiner2()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) UpgradeResourceMiner3(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeResourceMiner3()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) UpgradeBattery(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeBattery()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) UpgradeExoframe(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeExoframe()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) FillExoframe(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.FillExoframe()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) UpgradeShield(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeShield()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) FillShield(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.FillShield()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) UpgradeWeapon(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.UpgradeWeapon()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) FillWeapon(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	err := p.FillWeapon()
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) buyResource(p *models.Player, buy func(int) error, amount int) error {
	g.Lock()
	defer g.Unlock()
	err := buy(amount)
	if err != nil {
		return err
	}
	return g.store.SetPlayer(p)
}

func (g *Game) BuyResource1(p *models.Player, amount int) error {
	return g.buyResource(p, p.BuyResource1, amount)
}

func (g *Game) BuyResource2(p *models.Player, amount int) error {
	return g.buyResource(p, p.BuyResource2, amount)
}

func (g *Game) BuyResource3(p *models.Player, amount int) error {
	return g.buyResource(p, p.BuyResource3, amount)
}

type EmptySearch struct{}

func (es EmptySearch) Error() string { return "Can't find any enemy" }

func (g *Game) FindSingleEnemy(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	players, err := g.store.GetAllPlayers()
	if err != nil {
		return err
	}
	enemies := findSingle(p, players)
	log.Infof("Found %d enemies", len(enemies))
	if len(enemies) < 1 {
		return EmptySearch{}
	}
	enemy := enemies[rand.Intn(len(enemies))]
	p.Enemy = enemy
	return g.store.SetPlayer(p)
}

func findSingle(p *models.Player, players []*models.Player) []*models.Player {
	for d := 15; d <= 90; d += 5 {
		enemies := onlySingleEnemies(filterEnemies(p, players, d))
		if len(enemies) > 0 {
			log.Infof("Search deviation: %d", d)
			return enemies
		}
	}
	return []*models.Player{}
}

func (g *Game) FindGuildEnemy(p *models.Player) error {
	g.Lock()
	defer g.Unlock()
	players, err := g.store.GetAllPlayers()
	if err != nil {
		return err
	}
	enemies := findGuild(p, players, p.Guild)
	log.Infof("Found %d enemies", len(enemies))
	if len(enemies) < 1 {
		return EmptySearch{}
	}
	enemy := enemies[rand.Intn(len(enemies))]
	p.Enemy = enemy
	return g.store.SetPlayer(p)
}

func findGuild(p *models.Player, players []*models.Player, guild string) []*models.Player {
	for d := 15; d <= 90; d += 5 {
		enemies := onlyGuildEnemies(filterEnemies(p, players, d), p.Guild)
		if len(enemies) > 0 {
			log.Infof("Search deviation: %d", d)
			return enemies
		}
	}
	return []*models.Player{}
}

type NoTarget struct{}

func (nt NoTarget) Error() string { return "" }

type AlreadyInBattle struct{}

func (AlreadyInBattle) Error() string { return "" }

type YouAlreadyInBattle struct{}

func (YouAlreadyInBattle) Error() string { return "" }

type HasImmune struct{}

func (HasImmune) Error() string { return "" }

func (g *Game) notifyGuildMemebersAttack(attacker, defender *models.Guild, exceptID int) {
	if attacker == nil {
		return
	}
	log.Infof("[notifyGuildMemebersAttack] %s %#v", attacker.Name, attacker.PlayerIDs)
	for _, playerID := range attacker.PlayerIDs {
		if playerID != exceptID {
			g.sendMessage(playerID, g.battleRender.RenderGuildAttack(attacker, defender))
		}
	}
}

func (g *Game) notifyGuildMemebersDefence(defender, attacker *models.Guild, exceptID int) {
	log.Infof("[notifyGuildMemebersDefence] %s %#v", defender.Name, defender.PlayerIDs)
	for _, playerID := range defender.PlayerIDs {
		if playerID != exceptID {
			g.sendMessage(playerID, g.battleRender.RenderGuildDefence(attacker))
		}
	}
}

type NoBattle struct{}

func (NoBattle) Error() string { return "no battle" }

func (g *Game) PlayerJoinGuild(p *models.Player, title string) (*models.Guild, error) {
	g.Lock()
	defer g.Unlock()
	p, err := g.store.GetPlayer(p.ID)
	if err != nil {
		return nil, err
	}
	guild, err := models.GuildByTitle(title)
	if err != nil {
		return nil, err
	}
	p.Guild = guild.Emoji
	p.GuildID = guild.ID
	err = g.addPlayerToGuild(p.ID, guild.ID)
	if err != nil {
		return nil, err
	}
	return &guild, g.store.SetPlayer(p)
}

func (g *Game) addPlayerToGuild(playerID, guildID int) error {
	guild, err := g.store.GetGuild(guildID)
	if err != nil {
		return err
	}
	guild.PlayerIDs = append(guild.PlayerIDs, playerID)
	return g.store.SetGuild(guild)
}

func (g *Game) removePlayerFromGuild(playerID, guildID int) error {
	guild, err := g.store.GetGuild(guildID)
	if err != nil {
		return nil
	}
	for i, id := range guild.PlayerIDs {
		if playerID == id {
			guild.PlayerIDs = append(guild.PlayerIDs[:i], guild.PlayerIDs[i+1:]...)
			return g.store.SetGuild(guild)
		}
	}
	return nil
}

func filterEnemies(p *models.Player, players []*models.Player, deviation int) []*models.Player {
	enemies := make([]*models.Player, 0)
	for _, enemy := range players {
		if p.ID == enemy.ID {
			continue
		}
		playerLevel := p.Implants.Exoframe
		enemyLevel := enemy.Implants.Exoframe
		if playerLevel >= enemyLevel-deviation && playerLevel < enemyLevel+deviation && playerCanBeAttacked(enemy) {
			enemies = append(enemies, enemy)
		}
	}
	return enemies
}

func onlySingleEnemies(players []*models.Player) []*models.Player {
	enemies := make([]*models.Player, 0)
	for _, enemy := range players {
		if enemy.NoGuild() {
			enemies = append(enemies, enemy)
		}
	}
	return enemies
}

func onlyGuildEnemies(players []*models.Player, guild string) []*models.Player {
	enemies := make([]*models.Player, 0)
	for _, enemy := range players {
		if !enemy.NoGuild() && enemy.Guild != guild {
			enemies = append(enemies, enemy)
		}
	}
	return enemies
}

func playerCanBeAttacked(p *models.Player) bool {
	return !p.InBattle() && p.Status() == nil
}

func (g *Game) TopPlayersPrestige() []*models.Player {
	g.Lock()
	defer g.Unlock()
	top := []*models.Player{}
	players, err := g.store.GetAllPlayers()
	if err != nil {
		log.Errorf("can't get players: %v", err)
		return top
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].Prestige > players[j].Prestige
	})
	if len(players) < 5 {
		top = players
	} else {
		top = players[:5]
	}
	return top
}

func (g *Game) TopPlayersCoins() []*models.Player {
	g.Lock()
	defer g.Unlock()
	top := []*models.Player{}
	players, err := g.store.GetAllPlayers()
	if err != nil {
		log.Errorf("can't get players: %v", err)
		return top
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].Coins > players[j].Coins
	})
	if len(players) < 5 {
		top = players
	} else {
		top = players[:5]
	}
	return top
}

func (g *Game) TopPlayersWins() []*models.Player {
	g.Lock()
	defer g.Unlock()
	top := []*models.Player{}
	players, err := g.store.GetAllPlayers()
	if err != nil {
		log.Errorf("can't get players: %v", err)
		return top
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].Wins > players[j].Wins
	})
	if len(players) < 5 {
		top = players
	} else {
		top = players[:5]
	}
	return top
}

func (g *Game) TopGuildsPrestige() []*models.Guild {
	g.Lock()
	defer g.Unlock()
	top := []*models.Guild{}
	guilds, err := g.store.GetAllGuilds()
	if err != nil {
		log.Errorf("cat't get guilds: %v", err)
		return top
	}
	sort.Slice(guilds, func(i, j int) bool {
		return g.guildPrestige(guilds[i]) > g.guildPrestige(guilds[j])
	})
	if len(guilds) < 5 {
		top = guilds
	} else {
		top = guilds[:5]
	}
	return top
}

func (g *Game) TopGuildsCoins() []*models.Guild {
	g.Lock()
	defer g.Unlock()
	top := []*models.Guild{}
	guilds, err := g.store.GetAllGuilds()
	if err != nil {
		log.Errorf("cat't get guilds: %v", err)
		return top
	}
	sort.Slice(guilds, func(i, j int) bool {
		return g.guildCoins(guilds[i]) > g.guildCoins(guilds[j])
	})
	if len(guilds) < 5 {
		top = guilds
	} else {
		top = guilds[:5]
	}
	return top
}

func (g *Game) TopGuildsWins() []*models.Guild {
	g.Lock()
	defer g.Unlock()
	top := []*models.Guild{}
	guilds, err := g.store.GetAllGuilds()
	if err != nil {
		log.Errorf("cat't get guilds: %v", err)
		return top
	}
	sort.Slice(guilds, func(i, j int) bool {
		return g.guildWins(guilds[i]) > g.guildWins(guilds[j])
	})
	if len(guilds) < 5 {
		top = guilds
	} else {
		top = guilds[:5]
	}
	return top
}

func (g *Game) guildPlayers(guild *models.Guild) []*models.Player {
	players := []*models.Player{}
	for _, id := range guild.PlayerIDs {
		player, err := g.store.GetPlayer(id)
		if err != nil {
			log.Errorf("unagle to get player: %v", err)
			return players
		}
		players = append(players, player)
	}
	return players
}

func (g *Game) guildPrestige(guild *models.Guild) int {
	players := g.guildPlayers(guild)
	total := 0
	for _, player := range players {
		total += player.Prestige
	}
	guild.Prestige = total
	return total
}

func (g *Game) guildCoins(guild *models.Guild) int {
	players := g.guildPlayers(guild)
	total := 0
	for _, player := range players {
		total += player.Coins
	}
	guild.Coins = total
	return total
}

func (g *Game) guildWins(guild *models.Guild) int {
	players := g.guildPlayers(guild)
	total := 0
	for _, player := range players {
		total += player.Wins
	}
	guild.Wins = total
	return total
}

func (g *Game) GuildStat(playerID int) *models.GuildStat {
	g.Lock()
	defer g.Unlock()
	stat := &models.GuildStat{}
	player, err := g.store.GetPlayer(playerID)
	if err != nil {
		log.Errorf("unable to get player: %v", err)
		return stat
	}
	guild, err := g.store.GetGuild(player.GuildID)
	if err != nil {
		log.Errorf("unable to get guild: %v", err)
		return stat
	}
	stat.Title = guild.Title()
	stat.Coins = g.guildCoins(guild)
	stat.Prestige = g.guildPrestige(guild)
	stat.Wins = g.guildWins(guild)
	stat.PlayersNumber = len(guild.PlayerIDs)
	players := g.guildPlayers(guild)
	sort.Slice(players, func(i, j int) bool {
		return players[i].Prestige > players[j].Prestige
	})
	stat.Leader = players[0]
	return stat
}
