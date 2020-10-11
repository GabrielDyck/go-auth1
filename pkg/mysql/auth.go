package mysql

type Auth interface {
	IsAuthenticated(token string) (bool, error)
}


func (c *client) IsAuthenticated(token string) (bool, error) {
	row, err := c.db.Query("SELECT COUNT(1) FROM SESSION_TOKENS WHERE TOKEN = ?;", token)

	if err != nil {
		return false, err
	}
	var count int
	if !row.Next() {
		return false, nil
	}
	err = row.Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 1, nil
}