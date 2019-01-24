package game

import (
	"fmt"
	"sort"
	"time"

	"github.com/yanzay/log"
	"github.com/yanzay/picasso/models"
)

func (g *Game) battleLoop() {
	log.Infof("Starting battle loop")
	for result := range g.battleResults {
		err := g.processBattleResult(result)
		if err != nil {
			log.Errorf("unable to process result: %v", err)
		}
	}
}

func (g *Game) processBattleResult(result *models.BattleResult) error {
	g.Lock()
	defer g.Unlock()
	log.Infof("Processing battle result for player %d", result.PlayerID)
	p, err := g.store.GetPlayer(result.PlayerID)
	if err != nil {
		return fmt.Errorf("unable to get player %d", result.PlayerID)
	}
	p.ApplyBattleResult(result)
	g.SendNotification(p, result)
	return g.store.SetPlayer(p)
}

func (g *Game) setGuildBattleID(guildID, battleID int) (*models.Guild, error) {
	guild, err := g.store.GetGuild(guildID)
	if err != nil {
		return nil, err
	}
	guild.BattleID = battleID
	return guild, g.store.SetGuild(guild)
}

func (g *Game) AttackEnemy(playerID int) (*models.Battle, error) {
	g.Lock()
	defer g.Unlock()
	p, err := g.store.GetPlayer(playerID)
	if err != nil {
		return nil, err
	}
	if p.InBattle() {
		return nil, YouAlreadyInBattle{}
	}
	if p.Enemy == nil {
		return nil, NoTarget{}
	}
	enemy, err := g.store.GetPlayer(p.Enemy.ID)
	if err != nil {
		return nil, err
	}
	if enemy.InBattle() {
		return nil, AlreadyInBattle{}
	}
	if enemy.Status() != nil {
		return nil, HasImmune{}
	}
	battle := models.NewBattle(p, enemy)
	err = g.store.SetBattle(battle)
	if err != nil {
		return nil, err
	}
	p.BattleID = battle.ID
	enemy.BattleID = battle.ID
	err = g.store.SetPlayer(p)
	if err != nil {
		return nil, err
	}
	err = g.store.SetPlayer(enemy)
	if err != nil {
		return nil, err
	}
	if !enemy.NoGuild() {
		var attackGuild *models.Guild
		if !p.NoGuild() {
			attackGuild, err = g.setGuildBattleID(p.GuildID, battle.ID)
			if err != nil {
				return nil, err
			}
		}
		defenceGuild, err := g.setGuildBattleID(enemy.GuildID, battle.ID)
		if err != nil {
			return nil, err
		}
		g.notifyGuildMemebersDefence(defenceGuild, attackGuild, enemy.ID)
		g.notifyGuildMemebersAttack(attackGuild, defenceGuild, p.ID)
	}
	g.sendMessage(enemy.ID, g.battleRender.RenderAttackedBy(p))
	go g.startBattle(battle.ID, battle.Start)
	return battle, nil
}

func (g *Game) JoinGuildBattle(playerID int) (*models.Battle, error) {
	g.Lock()
	defer g.Unlock()
	p, err := g.store.GetPlayer(playerID)
	if err != nil {
		return nil, err
	}
	if p.InBattle() {
		return nil, AlreadyInBattle{}
	}
	guild, err := g.store.GetGuild(p.GuildID)
	if err != nil {
		return nil, err
	}
	if !guild.InBattle() {
		return nil, NoBattle{}
	}
	battle, err := g.store.GetBattle(guild.BattleID)
	if err != nil {
		return nil, err
	}
	attack := false
	if battle.AttackGuild == p.Guild {
		battle.Attack = append(battle.Attack, p)
		attack = true
	} else {
		battle.Defence = append(battle.Defence, p)
	}
	err = g.store.SetBattle(battle)
	if err != nil {
		return nil, err
	}
	p.BattleID = battle.ID
	err = g.store.SetPlayer(p)
	if err != nil {
		return nil, err
	}
	for _, playerID := range battle.PlayerIDs() {
		g.sendMessage(playerID, g.battleRender.RenderJoinedGuildBattle(p, attack))
	}
	return battle, nil
}

func (g *Game) startBattle(battleID int, start time.Time) {
	log.Infof("Starting battle %d: start[%s]", battleID, start)
	end := start.Add(models.BattleDuration)
	<-time.After(time.Until(end))
	results, err := g.endBattle(battleID)
	if err != nil {
		log.Errorf("unable to get battle results: %v", err)
	}
	for _, result := range results {
		g.battleResults <- result
	}
}

func (g *Game) endBattle(battleID int) ([]*models.BattleResult, error) {
	g.Lock()
	defer g.Unlock()
	log.Infof("Ending battle %d", battleID)
	battle, err := g.store.GetBattle(battleID)
	if err != nil {
		return nil, err
	}
	results := battle.Results()
	battles, err := g.store.GetAllBattles()
	sort.Slice(battles, func(i, j int) bool {
		return battles[i].Start.After(battles[j].Start)
	})
	if err != nil {
		return nil, err
	}
	log.Infof("There are %d battles right now", len(battles))
	if battle.AttackGuild != "" && battle.Guild() {
		guild, err := models.GuildByEmoji(battle.AttackGuild)
		if err != nil {
			return nil, err
		}
		anotherBattleID := 0
		for _, b := range battles {
			if battle.ID != b.ID && b.Guild() && (b.AttackGuild == guild.Emoji || b.DefenceGuild == guild.Emoji) {
				log.Infof("Attack guild %s has another battle: %d", guild.Emoji, b.ID)
				anotherBattleID = b.ID
			}
		}
		_, err = g.setGuildBattleID(guild.ID, anotherBattleID)
		if err != nil {
			return nil, err
		}
	}
	if battle.DefenceGuild != "" {
		guild, err := models.GuildByEmoji(battle.DefenceGuild)
		if err != nil {
			return nil, err
		}
		anotherBattleID := 0
		for _, b := range battles {
			if battle.ID != b.ID && b.Guild() && (b.AttackGuild == guild.Emoji || b.DefenceGuild == guild.Emoji) {
				log.Infof("Defence build %s has another battle: %d", guild.Emoji, b.ID)
				anotherBattleID = b.ID
			}
		}
		_, err = g.setGuildBattleID(guild.ID, anotherBattleID)
		if err != nil {
			return nil, err
		}
	}
	err = g.store.DeleteBattle(battleID)
	return results, err
}
