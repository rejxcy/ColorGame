package player

type Player struct {
	name       string
	timeRecord float64
}

func (p *Player) UpdateRecord(time float64) {
	
	// timeRecord < 1 代表尚未有紀錄，需進行更新
	if p.timeRecord < 1.0 ||p.timeRecord > time {
		p.timeRecord = time
		UpdatePlayerRecord(p)
	}
}
