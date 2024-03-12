package common

type Hero int32

const (
	HeroCatcher Hero = iota
	HeroFishguard
	HeroImp
	HeroKnight
	HeroMonkeydong
	HeroNosedman
	HeroPitboy
	HeroSpike
	HeroTreestor
	HeroWedger
)

var (
	HeroToName = map[Hero]string{
		HeroCatcher:    "Catcher",
		HeroFishguard:  "Fishguard",
		HeroImp:        "Imp",
		HeroKnight:     "Knight",
		HeroMonkeydong: "Monkeydong",
		HeroNosedman:   "Nosedman",
		HeroPitboy:     "Pitboy",
		HeroSpike:      "Spike",
		HeroTreestor:   "Treestor",
		HeroWedger:     "Wedger",
	}
	NameToHero = map[string]Hero{
		"Catcher":    HeroCatcher,
		"Fishguard":  HeroFishguard,
		"Imp":        HeroImp,
		"Knight":     HeroKnight,
		"Monkeydong": HeroMonkeydong,
		"Nosedman":   HeroNosedman,
		"Pitboy":     HeroPitboy,
		"Spike":      HeroSpike,
		"Treestor":   HeroTreestor,
		"Wedger":     HeroWedger,
	}
)
