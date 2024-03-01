package main

type StateType int

const (
	Failed StateType = iota
	Success
	Conflict
	Overflowed
)

type ClientConfig struct {
	Delay    int  `toml:"delay"`
	Parallel int  `toml:"parallel"`
	Keep     bool `toml:"keep"`
	Verbose  bool `toml:"verbose"`
}

type SummaryConfig struct {
	Enabled  bool `toml:"enabled"`
	Interval int  `toml:"interval"`
}

type SessionConfig struct {
	BaseURI    string   `toml:"base_uri"`
	UserAgent  string   `toml:"user_agent"`
	Credential string   `toml:"credential"`
	Targets    []string `toml:"targets"`
}

type Config struct {
	Client  ClientConfig  `toml:"client"`
	Summary SummaryConfig `toml:"summary"`
	Session SessionConfig `toml:"session"`
}

type Result struct {
	CourseID string
	TeachID  string
	State    StateType
}
