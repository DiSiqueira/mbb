package config

type (
	Specification interface {
		Key() string
		Secret() string
		Exchange() string
	}

	specs struct {
		ExchangeKey    string `envconfig:"key" required:"true"`
		ExchangeSecret string `envconfig:"secret" required:"true"`
		ExchangeName   string `envconfig:"exchange" required:"true"`
	}
)

func NewSpecification() Specification {
	return &specs{}
}

func (s *specs) Key() string {
	return s.ExchangeKey
}

func (s *specs) Secret() string {
	return s.ExchangeSecret
}

func (s *specs) Exchange() string {
	return s.ExchangeName
}
