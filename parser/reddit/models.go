package reddit

type Post struct {
	Title       string
	Id          string
	Src         string
	Description string
}

type Update struct {
	Kind string `json:"kind"`
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Data struct {
				Title       string `json:"title"`
				Id          string `json:"id"`
				Src         string `json:"url_overridden_by_dest"`
				Description string `json:"description"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`

	// Children []Post `json:"data>children"`
}
