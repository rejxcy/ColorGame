package main

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestNewPlayerDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error opening test database: %v", err)
	}
	defer db.Close()

	err = createPlayerTable(db)
	if err != nil {
		t.Fatalf("error creating player table: %v", err)
	}

	// Query sqlite_master 看 table 有沒有建立
	query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", TableName)
	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("error checking if table was created: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatalf("table '%s' was not created", TableName)
	}
}

func TestPlayerDBInsert(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error opening test database: %v", err)
	}
	defer db.Close()

	playerDB := &playerDB{db:db}
	err = createPlayerTable(db)
	if err != nil {
		t.Fatalf("error creating player table: %v", err)
	}

	player := Player{
		name:      "Test",
		timeRecord:  10.5,
	}
	playerDB.insert(player)

	retrievedPlayer, err := playerDB.selectPlayerByName(player.name)
	if err != nil {
		t.Fatalf("error selecting player: %v", err)
	}

	if retrievedPlayer.name != player.name || retrievedPlayer.timeRecord != player.timeRecord {
		t.Errorf("Expected player name: %s, best time: %f, got name: %s, best time: %f",
			player.name, player.timeRecord, retrievedPlayer.name, retrievedPlayer.timeRecord)
	}
}

func TestPlayerDBUpdate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error opening test database: %v", err)
	}
	defer db.Close()

	playerDB := &playerDB{db:db}
	err = createPlayerTable(db)
	if err != nil {
		t.Fatalf("error creating player table: %v", err)
	}

	player := Player{
		name:      "Test",
		timeRecord:  10.5,
	}
	playerDB.insert(player)

	newBestTime := 9.8
	player.timeRecord = newBestTime
	playerDB.update(player)

	updatedPlayer, err := playerDB.selectPlayerByName(player.name)
	if err != nil {
		t.Fatalf("error selecting player: %v", err)
	}

	if updatedPlayer.timeRecord != newBestTime {
		t.Errorf("Expected updated best time: %f, got: %f", newBestTime, updatedPlayer.timeRecord)
	}
}

func TestPlayerDBDelete(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error opening test database: %v", err)
	}
	defer db.Close()

	playerDB := &playerDB{db:db}
	err = createPlayerTable(db)
	if err != nil {
		t.Fatalf("error creating player table: %v", err)
	}

	player := Player{
		name:      "Test",
		timeRecord:  10.5,
	}
	playerDB.insert(player)

	playerDB.deleteByName(player.name)
	deletedPlayer, err := playerDB.selectPlayerByName(player.name)
	if err == nil || deletedPlayer != nil {
		t.Errorf("Expected deleted player to be nil, got: %v", deletedPlayer)
	}
}
