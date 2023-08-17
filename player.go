package main

type Player struct {
	name       string
	timeRecord float64
}

func (p *Player) UpdateRecord(time float64) {
	if p.timeRecord > time {
		p.timeRecord = time
		UpdatePlayerRecord(p)
	}
}
