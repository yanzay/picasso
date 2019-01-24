package models

import "time"

var (
	GameMinute          = 59 * time.Second
	ClinicDuration      = 20 * GameMinute
	HospitalDuration    = 40 * GameMinute
	BattleDuration      = 5 * GameMinute
	PrestigeToJoinGuild = 200
)

const (
	PriceBuyResource1 = 2
	PriceBuyResource2 = 2
	PriceBuyResource3 = 2

	ResourceMiningSpeed = 15
)
