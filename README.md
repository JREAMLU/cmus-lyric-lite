# cmus-lyric-plus

[cmus](https://cmus.github.io/) lyrics viewer

# Update
Fork of [cmus-lyric](https://github.com/rockagen/cmus-lyric)

-   update termui to V3
-   update fetch from 163
-   remove comments view
-   add go mod

![](https://api.travis-ci.org/rockagen/cmus-lyric.svg?branch=master)

# Like
![](https://i.imgur.com/WNxuUZ7.png)

With tmux 
![](https://i.imgur.com/wL3tPZa.png)

help
```bash
usage:

 q or <C-c>: quit
 y         : view lyrics
 ?         : help

```

# Install
Download release binary [lyrics](https://github.com/JREAMLU/cmus-lyric-plus/releases)

`chmod u+x lyrics`

# How
Check cmus current file exist lyric,fetch from music.163.com if not found

# Requirements
`go` compile 

`termui` term ui

# Build
Install lyrics
```bash
make install
```

# Run
`./lyrics`

type `q` to quit

happy enjoy!
