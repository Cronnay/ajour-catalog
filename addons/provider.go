package addons

// Addon is the struct and definition on how it should look in the json-file
type Addon struct {
	ID                int32     `json:"id"`
	Name              string    `json:"name"`
	URL               string    `json:"url"`
	NumberOfDownloads uint64    `json:"number_of_downloads"`
	Summary           string    `json:"summary"`
	Versions          []Version `json:"versions"`
	Categories        []string  `json:"categories"`
	Source            Source    `json:"source"`
}

// Version is the struct and definition on how it should look in the json-file
type Version struct {
	Flavor      GameVersion `json:"flavor"`
	GameVersion string      `json:"game_version"`
	Date        string      `json:"date"`
	fileID      int
}

// Source is an enum of available sources to download from
type Source string

const (
	// Curse is biggest repo
	Curse Source = "Curse"
	// Tukui is a source for Elvui & Tukui
	Tukui Source = "Tukui"
	// WowI is similar to Curse
	WowI Source = "WowI"
	// Hub is unknown repo
	Hub Source = "Hub"
)

//GameVersion is an enum of available versions to download addons for
type GameVersion string

const (
	//Retail is current latest expansion
	Retail GameVersion = "Retail"
	//ClassicTBC is the current Classic version. May change to ClassicWotlk
	ClassicTBC GameVersion = "ClassicTbc"
	//ClassicEra is the current "stay forever" expansion
	ClassicEra GameVersion = "ClassicEra"
)

// Backend is an interface to implement getting addons
type Backend interface {
	GetAddons() ([]Addon, error)
}
