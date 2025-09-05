package bq

type Config struct {
	// map of route name to route
	Routes map[string]Route `yaml:"-"`
}

type Route struct {
	// map of stop point name to stop point
	StopPoints map[string]StopPoint `yaml:"-"`
}

type StopPoint struct {
	// map of stop letter to stop config
	Letters map[string]StopConfig `yaml:"-"`
}

type StopConfig struct {
	TflID   string `yaml:"tfl_id"`
	UserTag string `yaml:"user_tag"`
}
