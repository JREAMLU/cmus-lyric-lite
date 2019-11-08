package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/JREAMLU/cmus-lyric-plus/source"
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
				err := fetchLyric(song.File, song.Artist, song.Title, song.Duration, 0)
				if err == nil {
					curLyric = loadLyrics(song.File)
				}
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

		if cmus.CurLyric == nil {
			DrawEmpty()
			return
		}

		var tmpPos int
		for i, n := range cmus.Pkeys {
			if song.Position < n {
				tmpPos = cmus.Pkeys[i-1]
				break
			}
		}
		// the same line needn't move
		if cmus.CurPos == tmpPos && cmus.CurPos != 0 {
			return
		}

		cmus.CurPos = tmpPos

		list := make([]string, 2*len(cmus.Pkeys))
		idx, cline := 0, 0

		for _, pos := range cmus.Pkeys {
			data := cmus.CurLyric[pos]
			datas := make([]string, len(data))
			copy(datas, data)
			if cmus.CurPos == pos && pos != 0 {
				text := datas[0]
				if len(text) < 1 {
					text = "..."
				}
				datas[0] = "[" + text + "](fg:cyan)"
				if len(datas[1]) > 0 {
					datas[1] = "[" + datas[1] + "](fg:cyan)"
				}
				cline = idx

			}
			list[idx] = datas[0]
			idx++
			if len(datas[1]) > 0 {
				list[idx] = datas[1]
				idx++
			}
		}

		DrawList(list, cline)
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

func fetchLyric(file string, artlist string, title string, duration int, size int) error {
	var err error
	for _, lyric := range source.LyricSrc {
		err = lyric.FetchLyric(file, artlist, title, duration, size)
		if err == nil {
			return nil
		}
	}

	return err
}
