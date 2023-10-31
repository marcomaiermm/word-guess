package database

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/google/uuid"
	internals "github.com/marcomaiermm/word-guess/internals"
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
	ID      int64
	UUID    string
	Word    string
	Rows    int
	Symbols []string
	Won     bool
}

type GameInDb struct {
	ID        int64
	UUID      string
	Rows      int
	Symbols   string
	Word      string
	Won       bool
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

	_, err = DB.Exec("DELETE FROM game WHERE player_id = ?", id)
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
		&game.Symbols,
		&game.Rows,
		&game.Won,
		&game.Word,
	)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func CreateGame(rows int, player_id int64) (*GameInDb, error) {
	words, err := internals.GetWordsList()
	if err != nil {
		return nil, err
	}

	randomIndex := rand.Intn(len(words.Data))
	randomWord := words.Data[randomIndex]
	newUUID := uuid.New().String()

	row := DB.QueryRow("INSERT INTO game (uuid, word, rows, player_id) VALUES (?, ?, ?, ?) RETURNING *", newUUID, randomWord, rows, player_id)

	var game GameInDb

	err = row.Scan(
		&game.ID,
		&game.UUID,
		&game.Word,
		&game.Rows,
		&game.Symbols,
		&game.Won,
		&game.CreatedAt,
		&game.UpdatedAt,
		&game.PlayerId,
	)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

type UpdateGameParams struct {
	PlayerId int64
	UUID     string
	Symbols  string
}

func UpdateGame(params UpdateGameParams) (*Game, error) {
	existingGame, err := GetGameByUUID(params.UUID, params.PlayerId)
	if err != nil {
		return nil, err
	}

	if len(existingGame.Word) != len(params.Symbols) {
		return nil, errors.New("symbols length does not match word length")
	}

	if len(strings.Split(existingGame.Symbols, ";")) == existingGame.Rows {
		return nil, errors.New("no more rows available")
	}

	// cant contain special characters
	if internals.HasSpecialChar(params.Symbols) {
		return nil, errors.New("symbols cannot contain special characters")
	}

	won := existingGame.Word == params.Symbols
	rowSymbols := existingGame.Symbols + params.Symbols + ";"

	row := DB.QueryRow("UPDATE game SET symbols = ?, won = ? WHERE uuid = ? AND player_id = ? RETURNING id, uuid, word, won", rowSymbols, won, params.UUID, params.PlayerId)

	updatedGame := Game{
		Symbols: strings.Split(rowSymbols, ";"),
	}

	err = row.Scan(
		&updatedGame.ID,
		&updatedGame.UUID,
		&updatedGame.Word,
		&updatedGame.Won,
	)

	if err != nil {
		return nil, err
	}

	return &updatedGame, nil
}

func DeleteGame(uuid string, player_id int64) error {
	_, err := DB.Exec("DELETE FROM game WHERE uuid = ? AND player_id = ?", uuid, player_id)
	if err != nil {
		return err
	}

	return nil
}
