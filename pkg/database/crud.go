package database

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	internals "github.com/marcomaiermm/word-guess/internals"
)

type Game struct {
	ID      int64
	UUID    string
	Word    string
	Rows    int
	Symbols []string
}

type GameInDb struct {
	ID        int64
	UUID      string
	Rows      int
	Symbols   string
	WordID    int
	Word      string
	CreatedAt string
	UpdatedAt string
}

type WordsResponse struct {
	Data []string `json:"data"`
}

func GetGameByUUID(uuid string) (*GameInDb, error) {
	row := DB.QueryRow("SELECT game.id, game.uuid, game.rows, game.symbols, game.word_id, game.created_at, game.updated_at, word.word FROM game LEFT JOIN word ON game.word_id = word.id WHERE game.uuid = ?", uuid)

	var game GameInDb

	err := row.Scan(
		&game.ID,
		&game.UUID,
		&game.Rows,
		&game.Symbols,
		&game.WordID,
		&game.CreatedAt,
		&game.UpdatedAt,
		&game.Word,
	)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func CreateGame(rows int) (*GameInDb, error) {
	var word string
	var wordId int
	err := DB.QueryRow("SELECT id, word FROM Word ORDER BY RANDOM() LIMIT 1").Scan(&wordId, &word)
	if err != nil {
		return nil, err
	}

	newUUID := uuid.New().String()
	symbols := strings.Repeat(",", rows)

	row := DB.QueryRow("INSERT INTO game (uuid, rows, symbols, word_id) VALUES (?, ?, ?, ?) RETURNING id, uuid, rows, symbols, created_at, updated_at", newUUID, rows, symbols, wordId)

	game := GameInDb{
		Word: word,
	}

	err = row.Scan(
		&game.ID,
		&game.UUID,
		&game.Rows,
		&game.Symbols,
		&game.CreatedAt,
		&game.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

type UpdateGameParams struct {
	UUID    string
	Symbols string
}

func UpdateGame(params UpdateGameParams) (*GameInDb, error) {
	existingGame, err := GetGameByUUID(params.UUID)
	if err != nil {
		return nil, err
	}

	if len(existingGame.Word) != len(params.Symbols) {
		return nil, errors.New("symbols length does not match word length")
	}

	symbols := strings.Split(existingGame.Symbols, ",")[:existingGame.Rows]

	// cant contain special characters
	if internals.HasSpecialChar(params.Symbols) {
		return nil, errors.New("symbols cannot contain special characters")
	}

	for i, symbol := range symbols {
		if symbol == "" {
			symbols[i] = string(params.Symbols)
			break
		}
	}
	rowSymbols := strings.Join(symbols, ",")

	game := GameInDb{
		Word: existingGame.Word,
	}

	err = DB.QueryRow("UPDATE game SET symbols = ? WHERE uuid = ? RETURNING id, uuid, rows, symbols, created_at, updated_at", rowSymbols, params.UUID).Scan(
		&game.ID,
		&game.UUID,
		&game.Rows,
		&game.Symbols,
		&game.CreatedAt,
		&game.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func DeleteGame(uuid string) error {
	_, err := DB.Exec("DELETE FROM game WHERE uuid = ?", uuid)
	if err != nil {
		return err
	}

	return nil
}
