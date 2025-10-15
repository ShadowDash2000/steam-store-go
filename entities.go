package steamstore

type App struct {
	AppID             uint   `json:"appid"`
	Name              string `json:"name"`
	LastModified      uint   `json:"last_modified"`
	PriceChangeNumber uint   `json:"price_change_number"`
}

type AppListQuery struct {
	IfModifiedSince         uint   `url:"if_modified_since,omitempty"`
	HaveDescriptionLanguage string `url:"have_description_language,omitempty"`
	IncludeGames            bool   `url:"include_games,omitempty"`
	IncludeDlc              bool   `url:"include_dlc,omitempty"`
	IncludeSoftware         bool   `url:"include_software,omitempty"`
	IncludeVideos           bool   `url:"include_videos,omitempty"`
	IncludeHardware         bool   `url:"include_hardware,omitempty"`
	LastAppId               uint   `url:"last_appid,omitempty"`
	MaxResults              uint   `url:"max_results,omitempty"`
}

type AppListResponse struct {
	Response struct {
		Apps            []App `json:"apps"`
		HaveMoreResults bool  `json:"have_more_results"`
		LastAppId       uint  `json:"last_appid"`
	} `json:"response"`
}

type AppDetail struct {
	Type                 string            `json:"type"`
	Name                 string            `json:"name"`
	SteamAppId           uint              `json:"steam_appid"`
	RequiredAge          int               `json:"required_age"`
	IsFree               bool              `json:"is_free"`
	ControllerSupport    string            `json:"controller_support"`
	DetailedDescription  string            `json:"detailed_description"`
	AboutTheGame         string            `json:"about_the_game"`
	ShortDescription     string            `json:"short_description"`
	SupportedLanguages   string            `json:"supported_languages"`
	Reviews              string            `json:"reviews"`
	HeaderImage          string            `json:"header_image"`
	CapsuleImage         string            `json:"capsule_image"`
	CapsuleImageV5       string            `json:"capsule_imagev5"`
	Website              string            `json:"website"`
	PcRequirements       map[string]string `json:"pc_requirements"`
	MacRequirements      map[string]string `json:"mac_requirements"`
	LinuxRequirements    map[string]string `json:"linux_requirements"`
	LegalNotice          string            `json:"legal_notice"`
	ExtUserAccountNotice string            `json:"ext_user_account_notice"`
	Developers           []string          `json:"developers"`
	Publishers           []string          `json:"publishers"`
	PackageGroups        []PackageGroup    `json:"package_groups"`
	Platforms            map[string]bool   `json:"platforms"`
	Metacritic           struct {
		Score int    `json:"score"`
		Url   string `json:"url"`
	} `json:"metacritic"`
}

type PackageGroup struct {
	Name                    string `json:"name"`
	Title                   string `json:"title"`
	Description             string `json:"description"`
	SelectionText           string `json:"selection_text"`
	SaveText                string `json:"save_text"`
	DisplayType             int    `json:"display_type"`
	IsRecurringSubscription string `json:"is_recurring_subscription"`
	Subs                    []Sub  `json:"subs"`
}

type Sub struct {
	PackageId                uint    `json:"packageid"`
	PercentSavingsText       string  `json:"percent_savings_text"`
	PercentSavings           float64 `json:"percent_savings"`
	OptionText               string  `json:"option_text"`
	OptionDescription        string  `json:"option_description"`
	CanGetFreeLicense        string  `json:"can_get_free_license"`
	IsFreeLicense            bool    `json:"is_free_license"`
	PriceInCentsWithDiscount uint    `json:"price_in_cents_with_discount"`
}

type AppDetailsResponse map[string]AppDetailEntry

type AppDetailEntry struct {
	Success bool      `json:"success"`
	Data    AppDetail `json:"data"`
}

type TagListQuery struct {
	Language        string `url:"language,omitempty"`
	HaveVersionHash string `url:"have_version_hash,omitempty"`
}

type TagListResponse struct {
	Response struct {
		VersionHash string `json:"version_hash"`
		Tags        []Tag  `json:"tags"`
	} `json:"response"`
}

type Tag struct {
	TagId uint   `json:"tagid"`
	Name  string `json:"name"`
}

type SteamSpyQuery struct {
	Request string `url:"request"`
	AppId   uint   `url:"appid"`
}

type SteamSpyAppDetailsResponseRaw struct {
	AppId          uint            `json:"appid"`
	Name           string          `json:"name"`
	Developer      string          `json:"developer"`
	Publisher      string          `json:"publisher"`
	ScoreRank      string          `json:"score_rank"`
	Positive       uint            `json:"positive"`
	Negative       uint            `json:"negative"`
	UserScore      uint            `json:"user_score"`
	Owners         string          `json:"owners"`
	AverageForever uint            `json:"average_forever"`
	Average2Weeks  uint            `json:"average_2weeks"`
	MedianForever  uint            `json:"median_forever"`
	Median2Weeks   uint            `json:"median_2weeks"`
	Price          string          `json:"price"`
	InitialPrice   string          `json:"initialprice"`
	Discount       string          `json:"discount"`
	CCU            uint            `json:"ccu"`
	Languages      string          `json:"languages"`
	Genre          string          `json:"genre"`
	Tags           map[string]uint `json:"tags"`
}

type SteamSpyAppDetailsResponse struct {
	AppId          uint            `json:"appid"`
	Name           string          `json:"name"`
	Developer      string          `json:"developer"`
	Publisher      string          `json:"publisher"`
	ScoreRank      string          `json:"score_rank"`
	Positive       uint            `json:"positive"`
	Negative       uint            `json:"negative"`
	UserScore      uint            `json:"user_score"`
	Owners         string          `json:"owners"`
	AverageForever uint            `json:"average_forever"`
	Average2Weeks  uint            `json:"average_2weeks"`
	MedianForever  uint            `json:"median_forever"`
	Median2Weeks   uint            `json:"median_2weeks"`
	Price          uint64          `json:"price"`
	InitialPrice   uint64          `json:"initialprice"`
	Discount       uint64          `json:"discount"`
	CCU            uint            `json:"ccu"`
	Languages      string          `json:"languages"`
	Genre          string          `json:"genre"`
	Tags           map[string]uint `json:"tags"`
}
