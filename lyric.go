package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Listen listen cmus info
func Listen(cmus *Cmus) {
	song := cmus.Remote()
	if song.Position > 0 {
		if cmus.CurFile != song.File {
			cmus.CurFile = song.File
			curLyric := loadLyrics(song.File)
			if curLyric == nil {
				log.Println("fetching..")
				fetchLyric(song.File, song.Artist, song.Title, song.Duration)
				curLyric = loadLyrics(song.File)
			}

			// pos keys
			pkeys := make([]int, 0, len(curLyric))
			for k := range curLyric {
				pkeys = append(pkeys, k)
			}
			sort.Ints(pkeys)
			cmus.Pkeys = pkeys
			cmus.CurLyric = curLyric
		}

		return
	}

	DrawEmpty()
}

func loadLyrics(path string) map[int][]string {

	pathIdx := strings.LastIndexAny(path, ".")

	lpath := path[:pathIdx] + ".lyric"
	tlpath := path[:pathIdx] + ".t.lyric"

	titleIdx := strings.LastIndexAny(path, "/")
	title := path[titleIdx+1 : pathIdx]

	content, err := ioutil.ReadFile(lpath)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(content), "\n")

	var tlines []string
	tcontent, err := ioutil.ReadFile(tlpath)
	if err == nil {
		tlines = strings.Split(string(tcontent), "\n")
	}

	m := make(map[int][]string)

	lyricMap := buildLyricMap(lines)
	tlyricMap := buildLyricMap(tlines)

	for k, v := range lyricMap {
		t1 := v
		t2 := tlyricMap[k]
		m[k] = []string{t1, t2}
	}
	m[0] = []string{title, ""}
	return m
}
func buildLyricMap(lyric []string) map[int]string {
	m := make(map[int]string)
	re := regexp.MustCompile("^\\[([0-9]+):([0-9]+).*](.*)")
	for _, v := range lyric {
		ar := re.FindStringSubmatch(v)
		if len(ar) > 3 {
			mi, _ := strconv.Atoi(ar[1])
			sec, _ := strconv.Atoi(ar[2])

			pos := 60*mi + sec

			m[pos] = ar[3]
		}

	}
	return m
}

// TODO:
func fetchLyric(file string, artlist string, title string, duration int) {

}
