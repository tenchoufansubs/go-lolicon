package nyaa

// Torrent contains torrent information.
type Torrent struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Downloads int    `json:"downloads"`
}
