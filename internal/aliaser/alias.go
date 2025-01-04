package aliaser

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed aliases.yaml
var aliases []byte

type item struct {
	UserId int64  `yaml:"user_id"`
	Alias  string `yaml:"alias"`
}

type config struct {
	Aliases []item `yaml:"aliases"`
}

type aliaser struct {
	mem map[int64]string
}

func New() *aliaser {
	cfg := config{}
	err := yaml.Unmarshal(aliases, &cfg)
	if err != nil {
		return &aliaser{}
	}

	mem := make(map[int64]string, len(cfg.Aliases))
	for _, alias := range cfg.Aliases {
		mem[alias.UserId] = alias.Alias
	}

	return &aliaser{
		mem: mem,
	}
}

func (a *aliaser) Alias(userID int64) string {
	return a.mem[userID]
}
