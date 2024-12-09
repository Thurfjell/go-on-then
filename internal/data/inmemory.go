package data

import (
	"goonthen/internal/server/game"
	"sync"
)

type state struct {
	source  *game.State
	mapSync sync.RWMutex
}

const sad string = "(╯︵╰,)"
const littleBitExcited string = "∩( ・ω・)∩"
const reallyExcited string = "ᕕ(⌐■_■)ᕗ ♪♬"
const yolo string = "Yᵒᵘ Oᶰˡʸ Lᶤᵛᵉ Oᶰᶜᵉ"

func whatsMood(clicks, likes int) string {
	mood := sad
	diff := likes - clicks
	if diff >= 5 && diff < 10 {
		mood = littleBitExcited
	} else if diff >= 10 && diff < 15 {
		mood = reallyExcited
	} else if diff >= 15 {
		mood = yolo
	}

	return mood
}

func (m *state) Get() *game.State {
	return m.source
}

func (m *state) Like() {
	m.mapSync.Lock()
	defer m.mapSync.Unlock()

	likes := m.source.Likes + 1
	m.source.Likes = likes
	m.source.Mood = whatsMood(m.source.Clicks, likes)
}

func (m *state) Click() {
	m.mapSync.Lock()
	defer m.mapSync.Unlock()

	clicks := m.source.Clicks + 1
	m.source.Clicks = clicks
	m.source.Mood = whatsMood(clicks, m.source.Likes)
}

func NewInMemoryState() *state {
	return &state{
		source: &game.State{
			Clicks: 0,
			Likes:  0,
			Mood:   sad,
		},
	}
}
