package types

type PostGameCarnageReportResponse struct {
	Response        PostGameCarnageReport `json:"Response"`
	ErrorCode       int                   `json:"ErrorCode"`
	ErrorStatus     string                `json:"ErrorStatus"`
	ThrottleSeconds int                   `json:"ThrottleSeconds"`
}

type PostGameCarnageReport struct {
	ActivityDetails                 ActivityDetails              `json:"activityDetails"`
	Period                          string                       `json:"period"`
	ActivityWasStartedFromBeginning bool                         `json:"activityWasStartedFromBeginning"`
	StartingPhaseIndex              int                          `json:"startingPhaseIndex"`
	Entries                         []PostGameCarnageReportEntry `json:"entries"`
}

type ActivityDetails struct {
	ReferenceId    int64  `json:"referenceId"`
	ActivityHash   int64  `json:"directorActivityHash"`
	InstanceId     string `json:"instanceId"`
	Mode           int    `json:"mode"`
	Modes          []int  `json:"modes"`
	IsPrivate      bool   `json:"isPrivate"`
	MembershipType int    `json:"membershipType"`
}

type PostGameCarnageReportEntry struct {
	Player      PlayerInformation            `json:"player"`
	CharacterId string                       `json:"characterId"`
	Values      map[string]Metric            `json:"values"`
	Extended    *WeaponAndAbilityInformation `json:"extended"`
}

type PlayerInformation struct {
	DestinyUserInfo DestinyUserInfo `json:"destinyUserInfo"`
	CharacterClass  string          `json:"characterClass"`
	ClassHash       int64           `json:"classHash"`
	RaceHash        int64           `json:"raceHash"`
	GenderHash      int64           `json:"genderHash"`
	LightLevel      int             `json:"lightLevel"`
	EmblemHash      int64           `json:"emblemHash"`
}

type DestinyUserInfo struct {
	IconPath                    string `json:"iconPath"`
	IsPublic                    bool   `json:"isPublic"`
	MembershipId                string `json:"membershipId"`
	MembershipType              int    `json:"membershipType"`
	DisplayName                 string `json:"displayName"`
	BungieGlobalDisplayName     string `json:"bungieGlobalDisplayName"`
	BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
}

type Basic struct {
	Value        float32 `json:"value"`
	DisplayValue string  `json:"displayValue"`
}

type Metric struct {
	Basic Basic `json:"basic"`
}

type WeaponAndAbilityInformation struct {
	Weapons   []WeaponInformation `json:"weapons"`
	Abilities map[string]Metric   `json:"values"`
}

type WeaponInformation struct {
	ReferenceId int64             `json:"referenceId"`
	Values      map[string]Metric `json:"values"`
}
