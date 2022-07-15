package openweathermap

// Exclude defines the API response parts that can be excluded.
type Exclude string

// fields that can be excluded from an API response.
const (
	Current  Exclude = "current"
	Minutely         = "minutely"
	Hourly           = "hourly"
	Daily            = "daily"
	Alerts           = "alerts"
)

// Units defines the units used in the API response.
type Units string

// supported units.
const (
	Standard Units = "standard"
	Metric         = "metric"
	Imperial       = "imperial"
)

// Lang defines the language used in the API response.
type Lang string

// supported languages
const (
	AF    Lang = "af"    // Afrikaans
	AL    Lang = "al"    // Albanian
	AR    Lang = "ar"    // Arabic
	AZ    Lang = "az"    // Azerbaijani
	BG    Lang = "bg"    // Bulgarian
	CA    Lang = "ca"    // Catalan
	CZ    Lang = "cz"    // Czech
	DA    Lang = "da"    // Danish
	DE    Lang = "de"    // German
	EL    Lang = "el"    // Greek
	EN    Lang = "en"    // English
	EU    Lang = "eu"    // Basque
	FA    Lang = "fa"    // Persian (Farsi)
	FI    Lang = "fi"    // Finnish
	FR    Lang = "fr"    // French
	GL    Lang = "gl"    // Galician
	HE    Lang = "he"    // Hebrew
	HI    Lang = "hi"    // Hindi
	HR    Lang = "hr"    // Croatian
	HU    Lang = "hu"    // Hungarian
	ID    Lang = "id"    // Indonesian
	IT    Lang = "it"    // Italian
	JA    Lang = "ja"    // Japanese
	KR    Lang = "kr"    // Korean
	LA    Lang = "la"    // Latvian
	LT    Lang = "lt"    // Lithuanian
	MK    Lang = "mk"    // Macedonian
	NO    Lang = "no"    // Norwegian
	NL    Lang = "nl"    // Dutch
	PL    Lang = "pl"    // Polish
	PT    Lang = "pt"    // Portuguese
	PT_BR Lang = "pt_br" // Português Brasil
	RO    Lang = "ro"    // Romanian
	RU    Lang = "ru"    // Russian
	SV    Lang = "sv"    // Swedish
	SE    Lang = "se"    // Swedish
	SK    Lang = "sk"    // Slovak
	SL    Lang = "sl"    // Slovenian
	SP    Lang = "sp"    // Spanish
	ES    Lang = "es"    // Spanish
	SR    Lang = "sr"    // Serbian
	TH    Lang = "th"    // Thai
	TR    Lang = "tr"    // Turkish
	UA    Lang = "ua"    // Ukrainian
	UK    Lang = "uk"    // Ukrainian
	VI    Lang = "vi"    // Vietnamese
	ZH_CN Lang = "zh_cn" // Chinese Simplified
	ZH_TW Lang = "zh_tw" // Chinese Traditional
	ZU    Lang = "zu"    // Zulu

)
