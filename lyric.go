package main

// Listen listen cmus info
func Listen(cmus *Cmus) {
	song := cmus.Remote()
	if song.Position > 0 {
		return
	}

	DrawEmpty()
}
