package config

import (
	"fmt"

	"github.com/hrabalvojta/micro-dvdrental/logger"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	PRODUCTION     = "production"
	DEVELOPMENT    = "development"
	STAGING        = "staging"
	DEFAULT_CONSUL = "consul"
)

// Config describes our global application configuration element.
// This structure will wrap any local or remote Viper clients that
// may be used to retrieve configuration objects from remote service providers
// such as etcd or Consul.
//
// The struct should also be composed with any additional configuration struct
// that you define in your application, e.g. "*Env"
type Config struct {
	// v is the viper instance for our environment configuration
	v *viper.Viper

	// Env is a reference to our environment configuration object
	Env *Env
}

// Env describes the gokit-base environment configuration. For this app, these
// values represent JSON data stored in Consul.
type Env struct {
	// ApplicationEnvironment provides production, development, or staging environment specification
	ApplicationEnvironment string `mapstructure:"app_env"`

	// ApplicationToken is the JWT signing token used to compare incoming authentication requests with the auth middleware
	ApplicationToken string `mapstructure:"token"`

	// Debug flags whether the application is running in debugging mode or not (increased log verbosity, no "pm", etc)
	Debug bool `mapstructure:"debug"`

	// HTTPHost defines the hostname that the HTTP server will listen on e.g. localhost for development
	HTTPHost string `mapstructure:"http_host"`

	// HTTPPort defines the port that the HTTP server will listen on e.g. 80
	HTTPPort string `mapstructure:"http_port"`

	// LogPath is the storage path for Herbert/Monolog style log output
	LogPath string `mapstructure:"log_path"`

	// LogChannel defines the channel this application's logs will be tagged with. Within our "golang" app channel
	// we have defined channels by service. This value may be "gokit-base" for this project.
	LogChannel string `mapstructure:"channel"`
}

func Init() *Config {
	C := Config{}
	C.v = viper.New()

	// Configuring defaults
	C.v.SetDefault("ApplicationEnvironment", "development")
	C.v.SetDefault("ApplicationToken", "")
	C.v.SetDefault("Debug", true)
	C.v.SetDefault("HTTPHost", "0.0.0.0")
	C.v.SetDefault("HTTPPort", "8080")
	C.v.SetDefault("LogPath", ".")
	C.v.SetDefault("LogChannel", "gokit-users")

	//C.Env = &Env{
	//	ApplicationEnvironment: "a",
	//	ApplicationToken:       "b",
	//	Debug:                  false,
	//	HTTPHost:               "c",
	//	HTTPPort:               "8080",
	//	LogPath:                "e",
	//	LogChannel:             "f",
	//}

	C.v.SetConfigName("users_config") // name of config file (without extension)
	C.v.SetConfigType("yaml")         // REQUIRED if the config file does not have the extension in the name
	//C.v.AddConfigPath("/etc/micro-dvdrental/") // path to look for the config file in
	C.v.AddConfigPath("./config/") // call multiple times to add many search paths
	//C.v.AddConfigPath(".")         // optionally look for config in the working directory

	if err := C.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(err)
		}
	}

	// Bring our configuration values into our defined struct or die
	if err := C.v.Unmarshal(&C.Env); err != nil {
		panic(err)
	}

	fmt.Printf("CONFIG: %+v\n", C.Env)

	return &C
}

// IsDevelopment returns whether the application is in dev mode
func (a *Config) IsDevelopment() bool {
	return a.Env.ApplicationEnvironment == DEVELOPMENT
}

// IsStaging returns true if the application is running in a staging environment
func (a *Config) IsStaging() bool {
	return a.Env.ApplicationEnvironment == STAGING
}

// IsProduction returns true if the application is running in a production environment
func (a *Config) IsProduction() bool {
	return a.Env.ApplicationEnvironment == PRODUCTION
}

// IsDebugEnvironment returns true if the application is in debug mode
func (a *Config) IsDebugEnvironment() bool {
	return a.Env.Debug
}

// LogLevel returns the current the application logger level
func (a *Config) LogLevel() logger.LogLevel {
	if a.Env.ApplicationEnvironment == PRODUCTION {
		return logger.ERROR
	}

	return logger.VERBOSE
}
