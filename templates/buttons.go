package templates

import "github.com/yanzay/picasso/models"

var (
	// generic
	YesButton     = YesEmoji + " Yes"
	NoButton      = NoEmoji + " No"
	HomeButton    = HomeEmoji + " Home"
	UpgradeButton = UpgradeEmoji + " Upgrade"
	FillButton    = FillEmoji + " Fill"
	BackButton    = BackEmoji + " Back"

	// main menu
	InfoButton      = InfoEmoji + " Info"
	ImplantsButton  = PagerEmoji + " Implants"
	EquipmentButton = HammerEmoji + " Equipment"
	BattleButton    = BattleEmoji + " Battle"
	GuildButton     = GuildEmoji + " Guild"
	ShopButton      = ShopEmoji + " Shop"
	TopButton       = TopEmoji + " Top"
	HelpButton      = HelpEmoji + " Help"

	// implants menu
	BatteryButton  = BatteryEmoji + " Battery"
	ExoframeButton = FrameEmoji + " Exoframe"
	ShieldButton   = ShieldEmoji + " Shield"
	WeaponButton   = WeaponEmoji + " Weapon"

	// equipment menu
	CoinminerButton    = MinerEmoji + " Coin Miner"
	BackpackButton     = BackpackEmoji + " Backpack"
	NanofactoryButton  = MicroscopeEmoji + " Nano Factory"
	PartsfactoryButton = BoltEmoji + " Parts Factory"
	FoodfactoryButton  = FoodEmoji + " Food Factory"

	// shop
	BuyButton      = BuyEmoji + " Buy"
	NanobotsButton = MicroscopeEmoji + " Nanobots"
	PartsButton    = BoltEmoji + " Parts"
	FoodButton     = FoodEmoji + " Food"

	// battle
	AttackButton      = BattleEmoji + " Attack"
	SearchButton      = SearchEmoji + " Search"
	SearchGuildButton = SearchEmoji + " Guild"
	JoinButton        = GuildEmoji + " Join"

	// top
	PlayersButton  = PlayerEmoji + " Players"
	GuildsButton   = GuildEmoji + " Guilds"
	PrestigeButton = PrestigeEmoji + " Prestige"
	CoinsButton    = CoinEmoji + " Coins"
	WinsButton     = WinEmoji + " Wins"
)

var (
	mainMenuButtons = [][]string{
		[]string{InfoButton, EquipmentButton, ImplantsButton},
		[]string{BattleButton, GuildButton, ShopButton},
		[]string{TopButton, HelpButton},
	}

	implantsMenuButtons = [][]string{
		[]string{BatteryButton, ExoframeButton},
		[]string{ShieldButton, WeaponButton},
		[]string{BackButton, HomeButton},
	}

	equipmentMenuButtons = [][]string{
		[]string{CoinminerButton, BackpackButton},
		[]string{NanofactoryButton, PartsfactoryButton, FoodfactoryButton},
		[]string{BackButton, HomeButton},
	}

	coinminerMenuButtons = [][]string{
		[]string{InfoButton, UpgradeButton},
		[]string{BackButton, HomeButton},
	}

	backpackMenuButtons = [][]string{
		[]string{InfoButton, FillButton, UpgradeButton},
		[]string{BackButton, HomeButton},
	}

	resourceMinerMenuButtons = [][]string{
		[]string{InfoButton, UpgradeButton},
		[]string{BackButton, HomeButton},
	}

	batteryMenuButtons = [][]string{
		[]string{InfoButton, UpgradeButton},
		[]string{BackButton, HomeButton},
	}

	exoframeMenuButtons = [][]string{
		[]string{InfoButton, FillButton, UpgradeButton},
		[]string{BackButton, HomeButton},
	}

	shieldMenuButtons = [][]string{
		[]string{InfoButton, FillButton, UpgradeButton},
		[]string{BackButton, HomeButton},
	}

	weaponMenuButtons = [][]string{
		[]string{InfoButton, FillButton, UpgradeButton},
		[]string{BackButton, HomeButton},
	}

	shopMenuButtons = [][]string{
		[]string{InfoButton, BuyButton},
		[]string{BackButton, HomeButton},
	}

	shopBuyMenuButtons = [][]string{
		[]string{NanobotsButton, PartsButton, FoodButton},
		[]string{BackButton, HomeButton},
	}

	battleMenuButtons = [][]string{
		[]string{InfoButton, AttackButton, JoinButton},
		[]string{SearchGuildButton, SearchButton},
		[]string{BackButton, HomeButton},
	}

	guildViewMenuButtons = [][]string{
		[]string{InfoButton},
		[]string{BackButton, HomeButton},
	}

	topMenuButtons = [][]string{
		[]string{PlayersButton, GuildsButton},
		[]string{BackButton, HomeButton},
	}

	topPlayersMenuButtons = [][]string{
		[]string{PrestigeButton, CoinsButton, WinsButton},
		[]string{BackButton, HomeButton},
	}

	topGuildsMenuButtons = [][]string{
		[]string{PrestigeButton, CoinsButton, WinsButton},
		[]string{BackButton, HomeButton},
	}
)

func amountButtons(emoji string) [][]string {
	return [][]string{
		[]string{"100" + emoji, "500" + emoji, "1000" + emoji},
		[]string{"5000" + emoji, "10000" + emoji, "50000" + emoji},
		[]string{"100000" + emoji, "500000" + emoji, "1000000" + emoji},
		[]string{BackButton, HomeButton},
	}
}

func guildMenuButtons() [][]string {
	buttons := make([][]string, 0)
	for _, guild := range models.Guilds {
		buttons = append(buttons, []string{guild.Title()})
	}
	buttons = append(buttons, []string{BackButton, HomeButton})
	return buttons
}
