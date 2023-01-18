package pixiv_api_go

const (
	IllustRankModeDaily      IllustRankMode = "daily"    // 今日
	IllustRankModeWeekly     IllustRankMode = "weekly"   // 本周
	IllustRankModeMonthly    IllustRankMode = "monthly"  // 本月
	IllustRankModeRookie     IllustRankMode = "rookie"   // 新人
	IllustRankModeDailyAi    IllustRankMode = "daily_ai" // AI 生成
	IllustRankModeMale       IllustRankMode = "male"     // 受男性欢迎
	IllustRankModeFemale     IllustRankMode = "female"   // 受女性欢迎
	IllustRankModeDailyR18   IllustRankMode = "daily_r18"
	IllustRankModeWeeklyR18  IllustRankMode = "weekly_r18"
	IllustRankModeDailyR18Ai IllustRankMode = "daily_r18_ai"
	IllustRankModeMaleR18    IllustRankMode = "male_r18"
	IllustRankModeFemaleR18  IllustRankMode = "female_r18"
	IllustRankModeR18g       IllustRankMode = "r18g"
)

const (
	IllustRankContentAll    IllustRankContent = "all"    // 综合
	IllustRankContentIllust IllustRankContent = "illust" // 插画
	IllustRankContentUgoira IllustRankContent = "ugoira" // 动图
	IllustRankContentManga  IllustRankContent = "manga"  // 漫画
)

//================================================================

type IllustTypeCode int

const (
	IllustTypeIllust IllustTypeCode = 0
	IllustTypeManga  IllustTypeCode = 1
	IllustTypeUgoira IllustTypeCode = 2
)

var illustName = map[IllustTypeCode]string{
	IllustTypeIllust: "Illust",
	IllustTypeManga:  "Manga",
	IllustTypeUgoira: "Ugoira",
}

func IllustTypeName(code IllustTypeCode) string {
	if v, ok := illustName[code]; ok {
		return v
	}
	return "UNKNOWN"
}

func (i IllustTypeCode) MarshalJSON() ([]byte, error) {
	name := IllustTypeName(i)
	return []byte(`"` + name + `"`), nil
}

//================================================================

type AITypeCode int

const (
	AITypeUndefined     AITypeCode = 0
	AITypeNotAiGenerate AITypeCode = 1
	AITypeAiGenerate    AITypeCode = 2
)

var aiTypeCodeName = map[AITypeCode]string{
	AITypeUndefined:     "Undefined",
	AITypeNotAiGenerate: "NotAiGenerate",
	AITypeAiGenerate:    "AiGenerate",
}

func AITypeCodeName(code AITypeCode) string {
	if v, ok := aiTypeCodeName[code]; ok {
		return v
	}
	return "UNKNOWN"
}

func (a AITypeCode) MarshalJSON() ([]byte, error) {
	name := AITypeCodeName(a)
	return []byte(`"` + name + `"`), nil
}

//================================================================

type RestrictLevel int

const (
	RestrictLevelPublic  RestrictLevel = 0
	RestrictLevelMypixiv RestrictLevel = 1 // illust will only be visible to people who are added to your My pixiv
	RestrictLevelPrivate RestrictLevel = 2
)

var restrictLevelName = map[RestrictLevel]string{
	RestrictLevelPublic:  "Public",
	RestrictLevelMypixiv: "Mypixiv",
	RestrictLevelPrivate: "Private",
}

func RestrictName(level RestrictLevel) string {
	if v, ok := restrictLevelName[level]; ok {
		return v
	}
	return "UNKNOWN"
}

func (r RestrictLevel) MarshalJSON() ([]byte, error) {
	name := RestrictName(r)
	return []byte(`"` + name + `"`), nil
}

//================================================================

type XRestrictLevel int

const (
	XRestrictLevelSafe XRestrictLevel = 0
	XRestrictLevelR18  XRestrictLevel = 1
	XRestrictLevelR18G XRestrictLevel = 2
)

var xRestrictLevelName = map[XRestrictLevel]string{
	XRestrictLevelSafe: "Safe",
	XRestrictLevelR18:  "R18",
	XRestrictLevelR18G: "R18G",
}

func XRestrictName(level XRestrictLevel) string {
	if v, ok := xRestrictLevelName[level]; ok {
		return v
	}
	return "UNKNOWN"
}

func (xr XRestrictLevel) MarshalJSON() ([]byte, error) {
	name := XRestrictName(xr)
	return []byte(`"` + name + `"`), nil
}

//================================================================

type SanityLevelCode int

const (
	SanityLevelUnchecked SanityLevelCode = 0
	SanityLevelGray      SanityLevelCode = 1
	SanityLevelWhite     SanityLevelCode = 2
	SanityLevelSemiBlack SanityLevelCode = 4
	SanityLevelBlack     SanityLevelCode = 6
	SanityLevelIllegal   SanityLevelCode = 7
)

var sanityLevelCodeName = map[SanityLevelCode]string{
	SanityLevelUnchecked: "Unchecked",
	SanityLevelGray:      "Gray",
	SanityLevelWhite:     "White",
	SanityLevelSemiBlack: "SemiBlack",
	SanityLevelBlack:     "Black",
	SanityLevelIllegal:   "Illegal",
}

func SanityLevelName(code SanityLevelCode) string {
	if v, ok := sanityLevelCodeName[code]; ok {
		return v
	}
	return "UNKNOWN"
}

func (c SanityLevelCode) MarshalJSON() ([]byte, error) {
	name := SanityLevelName(c)
	return []byte(`"` + name + `"`), nil
}
