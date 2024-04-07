package config

// ServerProperties defines global config properties
type ServerProperties struct {
	// for Public configuration
	Bind              string `cfg:"bind"`
	Port              int    `cfg:"port"`
	Databases         int    `cfg:"databases"`
}

// Properties holds global config properties
var Properties *ServerProperties
