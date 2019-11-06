package source

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/JREAMLU/cmus-ly/encrypt"
)

// Ntes 163
type Ntes struct {
}

const (
	// LyricAPI get ly
	LyricAPI = "http://music.163.com/weapi/song/lyric?csrf_token="
	// SearchAPI search song
	SearchAPI = "http://music.163.com/weapi/search/get?csrf_token="
	// CommentAPI comment api
	CommentAPI = "http://music.163.com/weapi/v1/resource/comments/R_SO_4_%v/?csrf_token="
	// Cookie cookie
	Cookie = "os=pc; osver=Microsoft-Windows-10-Professional-build-10586-64bit; appver=2.0.3.131777; channel=netease; __remember_me=true"
	// UserAgent ua
	UserAgent = "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"
)

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
		err := n.DownloadLyric(sid, dir, name)
		if err != nil {
			return err
		}
	}

	return nil
}

// FindSongID find song id
func (n *Ntes) FindSongID(name string, artlist string, title string, duration int, size int) int {
	m := make(map[string]interface{})

	m["s"] = name
	m["type"] = 1
	m["limit"] = 10
	m["offset"] = 0
	m["total"] = true
	m["csrf_token"] = ""

	req, _ := json.Marshal(m)
	params, encSecKey, err := encrypt.EncParams(string(req))
	if err != nil {
		log.Printf("error: %v \n", err)
		return 0
	}

	resp, err := post(SearchAPI, params, encSecKey)
	if err != nil {
		return 0
	}

	ret := &Songs{}
	err = json.Unmarshal(resp, ret)
	if err != nil {
		log.Println(err)
		return 0
	}
	code := ret.Code

	if 200 != code {
		log.Printf("code: %v, msg: %v \n", code, ret.Result)
		return 0
	}

	if len(ret.Result.Songs) > 0 {
		for _, v := range ret.Result.Songs {
			dt := v.Duration / 1000
			if dt == duration {
				return v.ID
			}
		}
	}

	return 0
}

// DownloadLyric download ly
func (n *Ntes) DownloadLyric(id int, dir string, name string) error {
	lyrc, tlyrc, err := n.GetLyric(id)
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

	return nil
}

// GetLyric get lyric
func (n *Ntes) GetLyric(id int) (string, string, error) {
	var lyrc, tlyrc string
	m := make(map[string]interface{})

	m["id"] = id
	m["os"] = "pc"
	m["lv"] = -1
	m["kv"] = -1
	m["tv"] = -1
	m["os"] = "pc"
	m["csrf_token"] = ""

	req, _ := json.Marshal(m)
	params, encSecKey, err := encrypt.EncParams(string(req))
	if err != nil {
		log.Printf("error: %v \n", err)
		return "", "", err
	}

	resp, err := post(LyricAPI, params, encSecKey)
	if err != nil {
		log.Println(err)
		return lyrc, tlyrc, err
	}

	result := &Lyric{}

	err = json.Unmarshal(resp, result)

	if err != nil {
		log.Println(err)
		return lyrc, tlyrc, err
	}

	if 200 != result.Code {
		log.Printf("code: %v, result: %v \n", result.Code, result)
		return lyrc, tlyrc, nil
	}

	lyrc, tlyrc = result.Lrc.Lyric, result.Tlyric.Lyric
	return lyrc, tlyrc, nil
}

// NewNtes new
func NewNtes() LyricSource {
	return &Ntes{}
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

func post(_url, params, encSecKey string) ([]byte, error) {
	client := &http.Client{}
	form := url.Values{}
	form.Set("params", params)
	form.Set("encSecKey", encSecKey)

	request, err := http.NewRequest("POST", _url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Host", "music.163.com")
	request.Header.Set("Origin", "http://music.163.com")
	request.Header.Set("User-Agent", UserAgent)

	request.Header.Set("Cookie", Cookie)

	resp, err := client.Do(request)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	resBody, resErr := ioutil.ReadAll(resp.Body)
	if resErr != nil {
		log.Println(err)
		return nil, resErr
	}
	return resBody, nil
}

// Songs search songs
type Songs struct {
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

// Lyric lyric
type Lyric struct {
	Sgc bool `json:"sgc"`
	Sfy bool `json:"sfy"`
	Qfy bool `json:"qfy"`
	Lrc struct {
		Version int    `json:"version"`
		Lyric   string `json:"lyric"`
	} `json:"lrc"`
	Klyric struct {
		Version int    `json:"version"`
		Lyric   string `json:"lyric"`
	} `json:"klyric"`
	Tlyric struct {
		Version int    `json:"version"`
		Lyric   string `json:"lyric"`
	} `json:"tlyric"`
	Code int `json:"code"`
}
