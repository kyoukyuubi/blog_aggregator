package config

type Config struct {
    DBURL         string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}