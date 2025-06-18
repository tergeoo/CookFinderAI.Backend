package db

type Datasource struct {
	Host   string `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
	Port   int    `yaml:"port" env:"DB_PORT" envDefault:"6434"`
	User   string `yaml:"user" env:"DB_USER" envDefault:"test"`
	Pass   string `yaml:"pass" env:"DB_PASS" envDefault:"test"`
	Name   string `yaml:"name" env:"DB_NAME" envDefault:"test"`
	Schema string `yaml:"schema" env:"DB_SCHEMA" envDefault:"public"`
}
