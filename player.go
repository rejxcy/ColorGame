package main

type Player struct {
	name       string
	timeRecord float64
}

func NewPlayer(name string) *Player {
	// TODO: check player from DB
	// TODO: if not exit, create new player

	return &Player{
		name:       name,
		timeRecord: 0,
	}
}

func (p *Player) UpdateRecord(time float64) {
	if p.timeRecord > time {
		p.timeRecord = time
	}
}
