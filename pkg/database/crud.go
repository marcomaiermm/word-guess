package database

import (
	"errors"
	"math/rand"
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
}

type WordsResponse struct {
	Data []string `json:"data"`
}

func GetGameByUUID(uuid string) (*GameInDb, error) {
	row := DB.QueryRow("SELECT * FROM game WHERE uuid = ?", uuid)

	var game GameInDb

	err := row.Scan(
		&game.ID,
		&game.UUID,
		&game.Rows,
		&game.Symbols,
		&game.Word,
		&game.Won,
	)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func CreateGame(rows int) (*GameInDb, error) {
	words, err := internals.GetWordsList()
	if err != nil {
		return nil, err
	}

	randomIndex := rand.Intn(len(words.Data))
	randomWord := words.Data[randomIndex]
	newUUID := uuid.New().String()
	symbols := strings.Repeat(",", rows)

	row := DB.QueryRow("INSERT INTO game (uuid, word, rows, symbols) VALUES (?, ?, ?, ?) RETURNING *", newUUID, randomWord, rows, symbols)

	var game GameInDb

	err = row.Scan(
		&game.ID,
		&game.UUID,
		&game.Rows,
		&game.Symbols,
		&game.Word,
		&game.Won,
		&game.CreatedAt,
		&game.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

type UpdateGameParams struct {
	UUID     string
	PlayerId string
	Symbols  string
}

func UpdateGame(params UpdateGameParams) (*Game, error) {
	existingGame, err := GetGameByUUID(params.UUID)
	if err != nil {
		return nil, err
	}

	if len(existingGame.Word) != len(params.Symbols) {
		return nil, errors.New("symbols length does not match word length")
	}

	if len(strings.Split(existingGame.Symbols, ",")) == existingGame.Rows {
		return nil, errors.New("no more rows available")
	}

	// cant contain special characters
	if internals.HasSpecialChar(params.Symbols) {
		return nil, errors.New("symbols cannot contain special characters")
	}

	won := existingGame.Word == params.Symbols
	rowSymbols := existingGame.Symbols + params.Symbols + ","

	row := DB.QueryRow("UPDATE game SET symbols = ?, won = ? WHERE uuid = ? RETURNING id, uuid, word, won", rowSymbols, won, params.UUID)

	updatedGame := Game{
		Symbols: strings.Split(rowSymbols, ","),
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

func DeleteGame(uuid string) error {
	_, err := DB.Exec("DELETE FROM game WHERE uuid = ?", uuid)
	if err != nil {
		return err
	}

	return nil
}
