package source

// LyricSource source
type LyricSource interface {
	FetchLyric(file string, artlist string, title string, duration int, size int) error
	FindSongID(file string, artlist string, title string, duration int, size int) int
	DownloadLyric(id int) error
}

var LyricSrc = make(map[string]LyricSource)
