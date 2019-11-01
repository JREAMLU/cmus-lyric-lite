package source

// Ntes 163
type Ntes struct {
}

func init() {
	LyricSrc["ntes"] = NewNtes()
}

// FetchLyric fetch lyric and save to local
func (n *Ntes) FetchLyric(file string, artlist string, title string, duration int, size int) error {
	return nil
}

// FindSongID find song id
func (n *Ntes) FindSongID(file string, artlist string, title string, duration int, size int) int {
	return 0
}

// DownloadLyric download ly
func (n *Ntes) DownloadLyric(id int) error {
	return nil
}

// NewNtes new
func NewNtes() LyricSource {
	return &Ntes{}
}
