package gag

type Post struct {
	Id        string
	Title     string
	Type      string
	Src       string
	Permalink string
}

type Update struct {
	Data struct {
		NextCursor string `json:"nextCursor"`
		Children   []struct {
			Id        string `json:"id"`
			Title     string `json:"title"`
			Type      string `json:"type"` // Animated, Photo
			Permalink string `json:"url"`
			Images    struct {
				Image struct {
					Src string `json:"url"`
				} `json:"image700"`
				Animated struct {
					Src string `json:"url"`
				} `json:"image460sv"`
			} `json:"images"`
		} `json:"posts"`
	} `json:"data"`

	// Children []Post `json:"data>children"`
}
