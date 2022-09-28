// Credit
// Thanks to TUD NetSoc for this brilliant entry point.
// https://github.com/netsoc/webspaced

package config

import (
	"reflect"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// Config describes the configuration for Server
type Config struct {
	LogLevel log.Level `mapstructure:"log_level"`
	Timeouts struct {
		Startup  time.Duration
		Shutdown time.Duration
	}

	HTTP struct {
		ListenAddress string `mapstructure:"listen_address"`

		CORS struct {
			AllowedOrigins []string `mapstructure:"allowed_origins"`
		}
	}

	Database struct {
		Host     string `mapstructure:"host"`
		Name     string `mapstructure:"name"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	LDAP struct {
		URL               string `mapstrcutre:"url"`
		Bind              string `mapstructure:"bind"`
		Password          string `mapstructure:"password"`
		UserSearchBase    string `mapstructure:"user_search_base"`
		SocietySearchBase string `mapstructure:"society_search_base"`
		GroupSearchBase   string `mapstructure:"group_search_base"`
		SearchBase        string `mapstructure:"search_base"`
	}

	SocsPortal struct {
		WebservicesEndpoint                      string `mapstructure:"webservices_endpoint"`
		AjaxEndpoint                             string `mapstructure:"ajax_endpoint"`
		WebservicesUsername                      string `mapstructure:"webservices_username"`
		WebservicesPassword                      string `mapstructure:"webservices_password"`
		WebservicesMemberServiceMethodIndividual string `mapstructure:"webservices_member_service_method_individual"`
		WebservicesMemberServiceMethodAll        string `mapstructure:"webservices_member_service_method_all"`
		WebservicesMemberServiceSearchByOption   string `mapstructure:"webservices_member_service_search_by_option"`
		EventService                             string `mapstructure:"event_service"`
		EventServiceMethodIndividual             string `mapstructure:"event_service_method_individual"`
		EventServiceMethodAll                    string `mapstructure:"event_service_method_all"`
		EventServiceAction                       string `mapstructure:"event_service_action"`
	}
}

// StringToLogLevelHookFunc returns a mapstructure.DecodeHookFunc which parses a logrus Level from a string
func StringToLogLevelHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t != reflect.TypeOf(log.InfoLevel) {
			return data, nil
		}

		var level log.Level
		err := level.UnmarshalText([]byte(data.(string)))
		return level, err
	}
}

// StringToTemplateHookFunc returns a mapstructure.DecodeHookFunc which parses a template.Template from a string
func StringToTemplateHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.TypeOf(template.Template{}).Kind() {
			return data, nil
		}

		return template.New("anonymous").Parse(data.(string))
	}
}

// DecoderOptions enables necessary mapstructure decode hook functions
func DecoderOptions(config *mapstructure.DecoderConfig) {
	config.ErrorUnused = true
	config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		config.DecodeHook,
		mapstructure.StringToTimeDurationHookFunc(),
		StringToLogLevelHookFunc(),
		StringToTemplateHookFunc(),
	)
}

func loadSecret(parent interface{}, field string) error {
	return nil
}

// ReadSecrets loads values for secret config options from files
func (c *Config) ReadSecrets() error {
	return nil
}
