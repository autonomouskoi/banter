package banter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/autonomouskoi/datastruct/mapset"
)

func (cfg *Config) Validate() error {
	if cfg.IntervalSeconds < 30 {
		return errors.New("interval must be at least 30")
	}
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
