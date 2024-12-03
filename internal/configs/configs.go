package configs

import "github.com/spf13/viper"

type (
	config struct {
		DatabaseURL string 	`mapstructure:"DATABASE_URL"`
		SecretJWT   string 	`mapstructure:"SECRET_JWT"`
		PORT        string 	`mapstructure:"PORT"`
		ENV					string	`mapstructure:"ENV"`
	}
)

func Init(
	Path,
	ConfigType,
	ConfigName string,
) (*config, error) {
	viper.AutomaticEnv()
	viper.AddConfigPath(Path)
	viper.SetConfigType(ConfigType)
	viper.SetConfigName(ConfigName)

	err := viper.ReadInConfig()
	if err != nil {
		return &config{}, err
	}

	var config config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}


	return &config, nil
}