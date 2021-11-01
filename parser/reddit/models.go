package reddit

type Post struct {
	Id          string
	Title       string
	Type        string
	Src         string
	Description string
	Permalink   string
	IsVideo     bool
}

type Update struct {
	Kind string `json:"kind"`
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Data struct {
				Id          string `json:"id"`
				Title       string `json:"title"`
				Type        string
				Src         string `json:"url_overridden_by_dest"`
				Description string `json:"description"`
				Permalink   string `json:"permalink"`
				IsVideo     bool   `json:"is_video"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`

	// Children []Post `json:"data>children"`
}
