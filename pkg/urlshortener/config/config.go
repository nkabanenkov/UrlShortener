package config

type Config struct {
	Alphabet []rune
	Width    uint

	InMemory   bool
	DbHost     string
	DbPort     uint
	DbUser     string
	DbPassword string
	DbName     string
}
