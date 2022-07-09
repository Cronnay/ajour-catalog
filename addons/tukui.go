package addons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type TukuiProvider struct{}

type tukuiAddonResponse struct {
	Name          string      `json:"name"`
	Author        string      `json:"author"`
	URL           string      `json:"url"`
	Version       string      `json:"version"`
	Changelog     string      `json:"changelog"`
	Ticket        string      `json:"ticket"`
	Git           string      `json:"git"`
	ID            interface{} `json:"id"`
	Patch         string      `json:"patch"`
	Lastupdate    string      `json:"lastupdate"`
	WebURL        string      `json:"web_url"`
	Lastdownload  string      `json:"lastdownload"`
	DonateURL     string      `json:"donate_url"`
	SmallDesc     string      `json:"small_desc"`
	ScreenshotURL string      `json:"screenshot_url"`
	Downloads     interface{} `json:"downloads"`
	Category      string      `json:"category"`
	flavor        string
}

func NewTukui() Backend {
	return &TukuiProvider{}
}

const baseURL string = "https://www.tukui.org/api.php"

func tukuiURL() string {
	return fmt.Sprintf("%s?ui=tukui", baseURL)
}

func elvuiURL() string {
	return fmt.Sprintf("%s?ui=elvui", baseURL)
}

func addonsURL(flavor GameVersion) string {
	switch flavor {
	case ClassicEra:
		return fmt.Sprintf("%s?classic-addons=all", baseURL)
	case ClassicTBC:
		return fmt.Sprintf("%s?classic-tbc-addons=all", baseURL)
	case Retail:
		return fmt.Sprintf("%s?addons=all", baseURL)
	default:
		panic("Dont recognize flavor: " + flavor)
	}
}

func (t TukuiProvider) GetAddons(c chan []Addon) {
	addonsChannel := make(chan []Addon, 0)

	flavors := []GameVersion{Retail, ClassicEra, ClassicTBC}
	for _, flavor := range flavors {
		if flavor == Retail {
			go fetchFromURL(true, flavor, addonsURL(flavor), addonsChannel)
			go fetchFromURL(false, flavor, elvuiURL(), addonsChannel)
			go fetchFromURL(false, flavor, tukuiURL(), addonsChannel)
		} else {
			go fetchFromURL(true, flavor, addonsURL(flavor), addonsChannel)
		}
	}

	addons := make([]Addon, 0)
	for i := 0; i < 5; i++ {
		addonsFromChannel := <-addonsChannel
		for _, a := range addonsFromChannel {
			addons = append(addons, a)
		}
	}
	c <- addons
}

func (t tukuiAddonResponse) toAddon(flavor GameVersion) Addon {
	addonID := 0
	if s, isCorrectType := t.ID.(string); isCorrectType {
		var err error
		addonID, err = strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
	}

	if f, isCorrectType := t.ID.(float64); isCorrectType {
		addonID = int(f)
	}

	var downloads int64 = 0
	if s, isCorrectType := t.Downloads.(string); isCorrectType {
		var err error
		downloads, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	if f, isCorrectType := t.Downloads.(float64); isCorrectType {
		downloads = int64(f)
	}

	return Addon{
		ID:                int32(addonID),
		Name:              t.Name,
		URL:               t.WebURL,
		NumberOfDownloads: uint64(downloads),
		Summary:           t.SmallDesc,
		Versions: []Version{
			{
				Flavor:      flavor,
				GameVersion: t.Patch,
				Date:        t.Lastupdate,
			},
		},
		Categories: []string{t.Category},
		Source:     Tukui,
	}
}

func fetchFromURL(multiple bool, flavor GameVersion, url string, c chan []Addon) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	addons := make([]Addon, 0)
	if multiple {
		var data []tukuiAddonResponse

		json.Unmarshal(body, &data)
		for _, t := range data {
			addons = append(addons, t.toAddon(flavor))
		}
		c <- addons
	} else {
		var data tukuiAddonResponse
		json.Unmarshal(body, &data)
		addons = append(addons, data.toAddon(flavor))
		c <- addons
	}
}
