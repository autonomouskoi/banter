package banter

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/autonomouskoi/datastruct/mapset"
	"github.com/autonomouskoi/twitch"
)

var (
	TestUser = &twitch.User{
		Id:              "id",
		Login:           "login",
		DisplayName:     "display name",
		Type:            "type",
		BroadcasterType: "broadcaster type",
		Description:     "description",
		ProfileImageUrl: "profile image URL",
		OfflineImageUrl: "offline image URL",
		ViewCount:       123,
		Email:           "email",
		CreatedAt:       1740162543,
	}
)

// Validate Banter's config
func (cfg *Config) Validate() error {
	if cfg.IntervalSeconds < 30 {
		return errors.New("interval must be at least 30")
	}
	if err := cfg.validateBanters(); err != nil {
		return fmt.Errorf("validating banters: %w", err)
	}
	if err := cfg.validateGuestListCommands(); err != nil {
		return fmt.Errorf("validating guest list commands: %w", err)
	}
	return nil
}

func (cfg *Config) validateBanters() error {
	seen := mapset.MapSet[string]{}
	for _, banter := range cfg.Banters {
		if !strings.HasPrefix(banter.Command, "!") {
			return fmt.Errorf("command %q must have ! prefix", banter.Command)
		}
		if banter.Command == "!" {
			return errors.New(`command can't be just "!"`)
		}
		if banter.Text == "" {
			return fmt.Errorf("command %q can't have empty text", banter.Command)
		}
		banter.Command = strings.ToLower(banter.Command)
		if seen.Has(banter.Command) {
			return fmt.Errorf("duplicate command %q", banter.Command)
		}
		seen.Add(banter.Command)
	}
	return nil
}

func (cfg *Config) validateGuestListCommands() error {
	for name, glc := range cfg.GuestListCommands {
		if err := renderGuestListCommand(glc, TestUser, io.Discard); err != nil {
			return fmt.Errorf("command %s: %w", name, err)
		}
	}
	return nil
}
