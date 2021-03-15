package vonage

type Config struct {
	APIKey    string `env:"API_KEY,required"`
	APISecret string `env:"API_SECRET,required"`
	Sender    string `env:"SENDER,required"`
}

func ConfigSkeleton() Config { return Config{} }
