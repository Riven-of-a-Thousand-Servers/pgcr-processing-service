package utils

import (
  "strings"
  "errors"
  "rivenbot/types/entity"
)

func ValidateRaidName(name string) (entity.RaidName, error) {
	switch entity.RaidName(strings.ToUpper(strings.TrimSpace(name))) {
	case entity.SalvationsEdge:
		return entity.SalvationsEdge, nil
	case entity.CrotasEnd:
		return entity.CrotasEnd, nil
	case entity.RootOfNightmares:
		return entity.RootOfNightmares, nil
	case entity.KingsFall:
		return entity.KingsFall, nil
	case entity.VowOfTheDisciple:
		return entity.VowOfTheDisciple, nil
	case entity.VaultOfGlass:
		return entity.VaultOfGlass, nil
	case entity.DeepStoneCrypt:
		return entity.DeepStoneCrypt, nil
	case entity.GardenOfSalvation:
		return entity.GardenOfSalvation, nil
	case entity.LeviathanCrownOfSorrow:
		return entity.LeviathanCrownOfSorrow, nil
	case entity.LastWish:
		return entity.LastWish, nil
	case entity.LeviathanSpireOfStars:
		return entity.LeviathanSpireOfStars, nil
	case entity.LeviathanEaterOfWorlds:
		return entity.LeviathanEaterOfWorlds, nil
	case entity.Leviathan:
		return entity.Leviathan, nil
	case entity.ScourgeOfThePast:
		return entity.ScourgeOfThePast, nil
	default:
		return "", errors.New("invalid raid name")
	}
}

func ValidateRaidDifficulty(name string) (entity.RaidDifficulty, error) {
	switch entity.RaidDifficulty(strings.ToUpper(strings.TrimSpace(name))) {
	case entity.Normal:
		return entity.Normal, nil
	case entity.Prestige:
		return entity.Prestige, nil
	case entity.Master:
		return entity.Master, nil
	default:
		return "", errors.New("invalid raid difficulty")
	}
}

func MapCharacterClass(name string) entity.CharacterClass {
  switch entity.CharacterClass(strings.ToUpper(strings.TrimSpace(name))) {
  case entity.Hunter:
    return entity.Hunter
  case entity.Titan:
    return entity.Titan
  case entity.Warlock:
    return entity.Warlock
  default:
    return "" 
  } 
}

func MapCharacterRace(name string) entity.CharacterRace {
  switch entity.CharacterRace(strings.ToUpper(strings.TrimSpace(name))) {
  case entity.Awoken:
    return entity.Awoken
  case entity.Exo:
    return entity.Exo
  case entity.Human:
    return entity.Human
  default:
    return ""
  }
}

func MapCharacterGender(name string) entity.CharacterGender {
  switch entity.CharacterGender(strings.ToUpper(strings.TrimSpace(name))) {
  case entity.Male:
    return entity.Male
  case entity.Female:
    return entity.Female
  default:
    return ""
  }
}
