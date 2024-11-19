package main

type (
	SearchResponse struct {
		Took     int  `json:"took,omitempty"`
		TimedOut bool `json:"timed_out,omitempty"`
		Shards   struct {
			Total      int `json:"total,omitempty"`
			Successful int `json:"successful,omitempty"`
			Skipped    int `json:"skipped,omitempty"`
			Failed     int `json:"failed,omitempty"`
		} `json:"_shards,omitempty"`
		Hits struct {
			Total struct {
				Value    int    `json:"value,omitempty"`
				Relation string `json:"relation,omitempty"`
			} `json:"total,omitempty"`
			MaxScore float64 `json:"max_score,omitempty"`
			Hits     []struct {
				Index     string   `json:"_index,omitempty"`
				Type      string   `json:"_type,omitempty"`
				ID        string   `json:"_id,omitempty"`
				Score     float64  `json:"_score,omitempty"`
				Ignored   []string `json:"_ignored,omitempty"`
				Source    Source   `json:"_source,omitempty"`
				SourceStr string   `json:"-"`
			} `json:"hits,omitempty"`
		} `json:"hits,omitempty"`
	}

	Source struct {
		ObjectID        string `json:"objectID,omitempty"`
		UnifiedID       string `json:"UnifiedId,omitempty"`
		All             string `json:"All,omitempty"`
		Aliases         string `json:"Aliases,omitempty"`
		BroadCategory   string `json:"BroadCategory,omitempty"`
		ConsoleID       string `json:"ConsoleId,omitempty"`
		ConsoleName     string `json:"ConsoleName,omitempty"`
		ConsoleNickname string `json:"ConsoleNickname,omitempty"`
		Developer       string `json:"Developer,omitempty"`
		DiscCount       int    `json:"DiscCount,omitempty"`
		FarooName       string `json:"FarooName,omitempty"`
		Genre           string `json:"Genre,omitempty"`
		IsVariant       string `json:"IsVariant,omitempty"`
		Name            string `json:"Name,omitempty"`
		NormalizedName  string `json:"NormalizedName,omitempty"`
		NoSpacesName    string `json:"NoSpacesName,omitempty"`
		Pageviews       int    `json:"Pageviews,omitempty"`
		PhoneticName    string `json:"PhoneticName,omitempty"`
		Publisher       string `json:"Publisher,omitempty"`
		Region          string `json:"Region,omitempty"`
		RomanName       string `json:"RomanName,omitempty"`
		StemmedName     string `json:"StemmedName,omitempty"`
		TransposedName  string `json:"TransposedName,omitempty"`
	}
)
