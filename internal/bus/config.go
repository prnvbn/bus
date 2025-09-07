package bus

type Config struct {
	Arrivals []Arrival `yaml:"arrivals"`
}

type Arrival struct {
	Route     string `yaml:"route"`
	StopPoint string `yaml:"stopPoint"`
	Letter    string `yaml:"letter"`
	TflID     string `yaml:"tflID"`
}

func (a Arrival) Title() string       { return a.StopPoint }
func (a Arrival) Description() string { return "Stop " + a.Letter }
func (a Arrival) FilterValue() string { return a.StopPoint }
