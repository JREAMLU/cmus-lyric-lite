package main

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	cmd  = "cmus-remote"
	args = "-Q"
)

// Cmus cmus info
type Cmus struct {
	CurFile  string
	CurLyric map[int][]string
	CurPos   int
	Pkeys    []int
}

// Song info
type Song struct {
	Position int
	File     string
	Artist   string
	Title    string
	Duration int
}

// NewCmus new cmus
func NewCmus() *Cmus {
	return &Cmus{}
}

// Remote remote cmus info [position, file, artist, duration]
func (c *Cmus) Remote() *Song {
	cmd := exec.Command(cmd, args)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.Fatalf("\n\n> cmus not running.\n\n")
	}

	info := strings.Split(stdout.String(), "\n")

	if len(info) < 1 || len(info[0]) < 1 {
		log.Fatalf("\n\n> cmus not running.\n\n")
	}

	//status stopped
	status := strings.Split(info[0], " ")[1]
	if status != "playing" {
		return &Song{}
	}

	if status == "pause" {
		return &Song{Position: 1}
	}

	idx := strings.Index(info[1], " ") + 1

	duration := strings.Split(info[2], " ")[1]
	position := strings.Split(info[3], " ")[1]
	pos, _ := strconv.Atoi(position)
	dt, _ := strconv.Atoi(duration)

	return &Song{
		Position: pos,
		File:     info[1][idx:],
		Artist:   strings.Split(info[6], " ")[1],
		Title:    strings.Split(info[5], " ")[1],
		Duration: dt,
	}
}
