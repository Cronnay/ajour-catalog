package addons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//CurseProvider implements Backend, requires API key to fetch data
type CurseProvider struct {
	APIKey string
}

//NewCurse creates new instance of CurseProvider
func NewCurse(apiKey string) Backend {
	return &CurseProvider{
		APIKey: apiKey,
	}
}

func getFlavorFromGameID(gameID int) GameVersion {
	switch gameID {
	case 73246:
		return ClassicTBC
	case 67408:
		return ClassicEra
	case 517:
		return Retail
	default:
		fmt.Printf("Unknown version: %d", gameID)
		return ""
	}
}

//GetAddons will fetch all available addons from Curse for the three flavors
func (p CurseProvider) GetAddons() ([]Addon, error) {
	addons := make([]Addon, 0, 0)
	client := &http.Client{}

	flavors := []int{73246, 67408, 517}
	for _, flavor := range flavors {
		index := 0
		pageSize := 50
		numberOfAddons := pageSize

		for pageSize == numberOfAddons {
			endpoint := fmt.Sprintf("https://api.curseforge.com/v1/mods/search?gameId=1&gameVersionTypeId=%d&pageSize=%d&index=%d", flavor, pageSize, index)
			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				panic(err.Error())
			}
			req.Header.Add("x-api-key", p.APIKey)
			resp, err := client.Do(req)
			if err != nil {
				panic(err.Error())
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err.Error())
			}

			var data Mods
			json.Unmarshal(body, &data)

			for _, mod := range data.Data {
				addon := createAddonFromMod(mod, flavor)
				addons = append(addons, addon)
			}
			numberOfAddons = len(data.Data)
			index += pageSize

		}
	}

	return addons, nil
}

func createAddonFromMod(mod Data, flavor int) Addon {
	files := []File{}
	for _, file := range mod.LatestFiles {
		if file.ReleaseType == 1 || file.ReleaseType == 2 {
			files = append(files, file)
		}
	}

	versions := []Version{}
	for _, file := range files {
		shouldFilter := false
		for _, f := range files {
			shouldFilter = f.GameID == file.GameID && f.ID > file.ID
		}

		if !shouldFilter {
			versions = append(versions, createVersionFromFile(file, flavor))
		}
	}

	categories := make([]string, 0, 0)
	for _, cat := range mod.Categories {
		categories = append(categories, cat.Name)
	}

	return Addon{
		ID:                int32(mod.ID),
		Name:              mod.Name,
		URL:               mod.Links.WebsiteURL,
		NumberOfDownloads: uint64(mod.DownloadCount),
		Summary:           mod.Summary,
		Versions:          versions,
		Categories:        categories,
		Source:            Curse,
	}
}

func createVersionFromFile(file File, flavor int) Version {
	return Version{
		GameVersion: file.GameVersions[0],
		Flavor:      getFlavorFromGameID(flavor),
		Date:        fmt.Sprintf("%d", file.FileDate.UTC().Unix()),
	}
}

// Mods is the top level struct we get from Curse
type Mods struct {
	Data       []Data     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Links has multiple links to use if wanted
type Links struct {
	WebsiteURL string      `json:"websiteUrl"`
	WikiURL    string      `json:"wikiUrl"`
	IssuesURL  interface{} `json:"issuesUrl"`
	SourceURL  string      `json:"sourceUrl"`
}

// Category describes a category
type Category struct {
	ID               int       `json:"id"`
	GameID           int       `json:"gameId"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	URL              string    `json:"url"`
	IconURL          string    `json:"iconUrl"`
	DateModified     time.Time `json:"dateModified"`
	IsClass          bool      `json:"isClass"`
	ClassID          int       `json:"classId"`
	ParentCategoryID int       `json:"parentCategoryId"`
}

//Author is the ones who developed
type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

//Logo is the image on Curse
type Logo struct {
	ID           int    `json:"id"`
	ModID        int    `json:"modId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnailUrl"`
	URL          string `json:"url"`
}

//Hashes is unclear
type Hashes struct {
	Value string `json:"value"`
	Algo  int    `json:"algo"`
}

//Modules has some fingerprints
type Modules struct {
	Name        string `json:"name"`
	Fingerprint int64  `json:"fingerprint"`
}

//File describes each file in an addon
type File struct {
	ID                   int           `json:"id"`
	GameID               int           `json:"gameId"`
	ModID                int           `json:"modId"`
	IsAvailable          bool          `json:"isAvailable"`
	DisplayName          string        `json:"displayName"`
	FileName             string        `json:"fileName"`
	ReleaseType          int           `json:"releaseType"`
	FileStatus           int           `json:"fileStatus"`
	Hashes               []Hashes      `json:"hashes"`
	FileDate             time.Time     `json:"fileDate"`
	FileLength           int           `json:"fileLength"`
	DownloadCount        int           `json:"downloadCount"`
	DownloadURL          string        `json:"downloadUrl"`
	GameVersions         []string      `json:"gameVersions"`
	SortableGameVersions []interface{} `json:"sortableGameVersions"`
	Dependencies         []interface{} `json:"dependencies"`
	AlternateFileID      int           `json:"alternateFileId"`
	IsServerPack         bool          `json:"isServerPack"`
	FileFingerprint      int64         `json:"fileFingerprint"`
	Modules              []Modules     `json:"modules"`
}

// LatestFilesIndexes is unclear
type LatestFilesIndexes struct {
	GameVersion       string      `json:"gameVersion"`
	FileID            int         `json:"fileId"`
	Filename          string      `json:"filename"`
	ReleaseType       int         `json:"releaseType"`
	GameVersionTypeID int         `json:"gameVersionTypeId"`
	ModLoader         interface{} `json:"modLoader"`
}

//Data is 2nd level
type Data struct {
	ID                   int                  `json:"id"`
	GameID               int                  `json:"gameId"`
	Name                 string               `json:"name"`
	Slug                 string               `json:"slug"`
	Links                Links                `json:"links"`
	Summary              string               `json:"summary"`
	Status               int                  `json:"status"`
	DownloadCount        float64              `json:"downloadCount"`
	IsFeatured           bool                 `json:"isFeatured"`
	PrimaryCategoryID    int                  `json:"primaryCategoryId"`
	Categories           []Category           `json:"categories"`
	ClassID              int                  `json:"classId"`
	Authors              []Author             `json:"authors"`
	Logo                 Logo                 `json:"logo"`
	Screenshots          []interface{}        `json:"screenshots"`
	MainFileID           int                  `json:"mainFileId"`
	LatestFiles          []File               `json:"latestFiles"`
	LatestFilesIndexes   []LatestFilesIndexes `json:"latestFilesIndexes"`
	DateCreated          time.Time            `json:"dateCreated"`
	DateModified         time.Time            `json:"dateModified"`
	DateReleased         time.Time            `json:"dateReleased"`
	AllowModDistribution bool                 `json:"allowModDistribution"`
	GamePopularityRank   int                  `json:"gamePopularityRank"`
	IsAvailable          bool                 `json:"isAvailable"`
	ThumbsUpCount        int                  `json:"thumbsUpCount"`
}

// Pagination is to know how many pages we can fetch
type Pagination struct {
	Index       int `json:"index"`
	PageSize    int `json:"pageSize"`
	ResultCount int `json:"resultCount"`
	TotalCount  int `json:"totalCount"`
}
