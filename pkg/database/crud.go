package database

func GetItem(id int64) (*Item, error) {
	row := DB.QueryRow("SELECT * FROM item WHERE id = ?", id)

	var item Item

	err := row.Scan(
		&item.ID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
