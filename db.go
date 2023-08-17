package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)


type playerDB struct {
	db *sql.DB
}

const DBPath = "./players.db"
const TableName = "players"

const (
	PlayerName string = "name"
	TimeRecord string = "time_record"
)

func NewPlayerDB() (*playerDB, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return nil, err
	}

	err = createPlayerTable(db)
	if err != nil {
		return nil, err
	}

	return &playerDB{
		db: db,
	}, nil
}

func createPlayerTable(db *sql.DB) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		%s TEXT,
		%s FLOAT);`,
		TableName, PlayerName, TimeRecord)

	_, err := db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}


func (playerDB *playerDB) selectPlayerByName(name string) (*Player, error) {
	var player Player

	query := fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s = ?;`, PlayerName, TimeRecord, TableName, PlayerName)

	row := playerDB.db.QueryRow(query, name)
	err := row.Scan(&player.name, &player.timeRecord)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (playerDB *playerDB) insert(player Player) {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES (?, ?);`, TableName, PlayerName, TimeRecord)

	_, err := playerDB.db.Exec(query, player.name, player.timeRecord)
	if err != nil {
		log.Println(err.Error())
	}
}

func (playerDB *playerDB) update(player Player) {
	query := fmt.Sprintf(`UPDATE  %s SET %s=? WHERE %s=?;`, TableName, TimeRecord, PlayerName)

	_, err := playerDB.db.Exec(query, player.timeRecord, player.name)
	if err != nil {
		log.Println(err.Error())
	}
}

func (playerDB *playerDB) deleteByName(name string) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s=?;`, TableName, PlayerName)

	_, err := playerDB.db.Exec(query, name)
	if err != nil {
		log.Println(err.Error())
	}
}