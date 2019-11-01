package source

// LyricSource source
type LyricSource interface {
	FetchLyric(file string, artlist string, title string, duration int, size int) error
	FindSongID(name string, artlist string, title string, duration int, size int) int
	DownloadLyric(id int) (string, string, error)
}

// LyricSrc lyric source
var LyricSrc = make(map[string]LyricSource)
