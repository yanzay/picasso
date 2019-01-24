package game

import (
	"testing"
)

type EmptyWriter struct{}

func (EmptyWriter) Write([]byte) (int, error) { return 0, nil }

func TestGameUpgrade(t *testing.T) {
	p, err := store.GetPlayer(0)
	if err != nil {
		t.Fatalf("can't get player: %v", err)
	}
	err = game.UpgradeBattery(p)
	if err != nil {
		t.Fatalf("can't upgrade battery: %v", err)
	}
	err = game.UpgradeCoinMiner(p)
	if err != nil {
		t.Fatalf("can't upgrade battery: %v", err)
	}
	err = game.UpgradeBackpack(p)
	if err == nil {
		t.Fatalf("should be not enough resources: %v", err)
	}
}
