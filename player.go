package main

type Player struct {
	name   string
	bestTime float32
}


func NewPlayer(name string) *Player {
	// TODO: check player from DB
	// TODO: if not exit, create new player

    return &Player{
		name: name,
		bestTime: 0,
	}
}

func (p *Player) UpdateRecord(time float32) {
    if p.bestTime > time {
        p.bestTime = time
    }
}

