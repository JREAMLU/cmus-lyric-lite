package source

import (
	"io"
	"log"
	"os"
	"strings"
)

// Ntes 163
type Ntes struct {
}

func init() {
	LyricSrc["ntes"] = NewNtes()
}

// FetchLyric fetch lyric and save to local
func (n *Ntes) FetchLyric(file string, artlist string, title string, duration int, size int) error {
	pathIdx := strings.LastIndexAny(file, ".")
	titleIdx := strings.LastIndexAny(file, "/")
	dir := file[:titleIdx]
	name := file[titleIdx+1 : pathIdx]

	sid := n.FindSongID(name, "", "", duration, 0)

	if sid > 0 {
		lyrc, tlyrc, err := n.DownloadLyric(sid)
		if err != nil {
			return err
		}

		if len(lyrc) > 0 {
			path := dir + "/" + name + ".lyric"
			err := save(path, strings.NewReader(lyrc))
			if err != nil {
				return err
			}
		}

		if len(tlyrc) > 0 {
			path := dir + "/" + name + ".t.lyric"
			err := save(path, strings.NewReader(tlyrc))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// FindSongID find song id
// TODO:
func (n *Ntes) FindSongID(name string, artlist string, title string, duration int, size int) int {
	return 0
}

// DownloadLyric download ly
// TODO:
func (n *Ntes) DownloadLyric(id int) (string, string, error) {
	return "", "", nil
}

// NewNtes new
func NewNtes() LyricSource {
	return &Ntes{}
}

func getLyric(id int) (string, string) {
	return "", ""
}

func save(path string, src io.Reader) error {
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		log.Printf("Write eror: %v \n", err)
		return err
	}
	n, err := io.Copy(out, src)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("<- %v, size: %v\n", path, n)

	return err
}

// Song song
type Song struct {
	Result struct {
		Songs []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				ID        int           `json:"id"`
				Name      string        `json:"name"`
				PicURL    interface{}   `json:"picUrl"`
				Alias     []interface{} `json:"alias"`
				AlbumSize int           `json:"albumSize"`
				PicID     int           `json:"picId"`
				Img1V1URL string        `json:"img1v1Url"`
				Img1V1    int           `json:"img1v1"`
				Trans     interface{}   `json:"trans"`
			} `json:"artists"`
			Album struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Artist struct {
					ID        int           `json:"id"`
					Name      string        `json:"name"`
					PicURL    interface{}   `json:"picUrl"`
					Alias     []interface{} `json:"alias"`
					AlbumSize int           `json:"albumSize"`
					PicID     int           `json:"picId"`
					Img1V1URL string        `json:"img1v1Url"`
					Img1V1    int           `json:"img1v1"`
					Trans     interface{}   `json:"trans"`
				} `json:"artist"`
				PublishTime int64 `json:"publishTime"`
				Size        int   `json:"size"`
				CopyrightID int   `json:"copyrightId"`
				Status      int   `json:"status"`
				PicID       int64 `json:"picId"`
				Mark        int   `json:"mark"`
			} `json:"album"`
			Duration    int           `json:"duration"`
			CopyrightID int           `json:"copyrightId"`
			Status      int           `json:"status"`
			Alias       []interface{} `json:"alias"`
			Rtype       int           `json:"rtype"`
			Ftype       int           `json:"ftype"`
			Mvid        int           `json:"mvid"`
			Fee         int           `json:"fee"`
			RURL        interface{}   `json:"rUrl"`
			Mark        int           `json:"mark"`
		} `json:"songs"`
		SongCount int `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}
