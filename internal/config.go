package internal

type Config struct {
	ServerHost     string `json:"server_host"`
	ServerPort     int    `json:"server_port"`
	DbDsn          string `json:"db_dsn"`
	DbMaxOpenConns int    `json:"db_max_open_conns"`
	DbMaxIdleConns int    `json:"db_max_idle_conns"`
}
