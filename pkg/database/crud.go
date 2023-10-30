package database

import (
	"github.com/google/uuid"
)

type Player struct {
	ID    int64
	UUID  string
	Name  string
	Score int
}

type PlayerInDb struct {
	Player
	CreatedAt string
	UpdatedAt string
}

type Game struct {
	ID   int64
	UUID string
	Rows int
	Word string
}

type GameInDb struct {
	ID        int64
	UUID      string
	Rows      int
	Word      string
	CreatedAt string
	UpdatedAt string
	PlayerId  int64
}

type WordsResponse struct {
	Data []string `json:"data"`
}

func GetPlayer(id int64) (*Player, error) {
	row := DB.QueryRow("SELECT player.id, player.name, player.score FROM player WHERE id = ?", id)

	var player Player

	err := row.Scan(
		&player.ID,
		&player.Name,
		&player.Score,
	)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func CreatePlayer(name string) (*PlayerInDb, error) {
	newUUID := uuid.New().String()
	row := DB.QueryRow("INSERT INTO player (uuid, name, score) VALUES (?, ?) RETURNING id, name, score", newUUID, name, 0)

	var player PlayerInDb

	err := row.Scan(
		&player.ID,
		&player.UUID,
		&player.Name,
		&player.Score,
	)

	if err != nil {
		return nil, err
	}

	return &player, nil
}

type UpdatePlayerParams struct {
	ID    int64
	Name  string
	Score int
}

func UpdatePlayer(params UpdatePlayerParams) (*PlayerInDb, error) {
	// build args
	args := make([]interface{}, 0)

	// build query
	query := "UPDATE player SET "
	if params.Name != "" {
		query += "name = ?, "
		args = append(args, params.Name)
	}
	if params.Score != 0 {
		query += "score = ?, "
		args = append(args, params.Score)
	}

	query = query[:len(query)-2]
	query += " WHERE id = ? RETURNING id, uuid, name, score"

	args = append(args, params.ID)

	row := DB.QueryRow(query, args...)

	var player PlayerInDb

	err := row.Scan(
		&player.ID,
		&player.UUID,
		&player.Name,
		&player.Score,
	)

	if err != nil {
		return nil, err
	}

	return &player, nil
}

func DeletePlayer(id int64) error {
	_, err := DB.Exec("DELETE FROM player WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func GetGameByUUID(uuid string, player_id int64) (*GameInDb, error) {

	row := DB.QueryRow("SELECT * FROM game WHERE uuid = ? AND player_id = ?", uuid, player_id)

	var game GameInDb

	err := row.Scan(
		&game.ID,
		&game.UUID,
		&game.Rows,
		&game.Word,
	)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func CreateGame(rows int, player_id int64) (*GameInDb, error) {
	return nil, nil
}

type UpdateGameParams struct {
	UUID string
	Rows int
	Word string
}

func UpdateGame(params UpdateGameParams) (*GameInDb, error) {
	// build args
	args := make([]interface{}, 0)

	// build query
	query := "UPDATE game SET "
	if params.Rows != 0 {
		query += "rows = ?, "
		args = append(args, params.Rows)
	}
	if params.Word != "" {
		query += "word = ?, "
		args = append(args, params.Word)
	}

	query = query[:len(query)-2]
	query += " WHERE uuid = ? RETURNING id, uuid, rows, word"

	args = append(args, params.UUID)

	row := DB.QueryRow(query, args...)

	var game GameInDb

	err := row.Scan(
		&game.ID,
		&game.UUID,
		&game.Rows,
		&game.Word,
	)

	if err != nil {
		return nil, err
	}

	return &game, nil
}
