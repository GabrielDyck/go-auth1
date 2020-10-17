package mysql


type Logout interface {
	Logout (token string)(bool,error)
}



func (c *client) Logout(token string)(bool,error) {
	r, err := c.db.Exec("DELETE FROM SESSION_TOKENS WHERE TOKEN = ?;", token)

	if err != nil {
		return false,err
	}
	quantity, err := r.RowsAffected()
	if err != nil {
		return false,err
	}

	return quantity ==1 , nil
}