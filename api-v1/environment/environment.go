package environment

import "github.com/spf13/viper"

type EnvironmentConfiguration struct {
	Secret              string `mapstructure:"SECRET"`
	Database            string `mapstructure:"DATABASE"`
	Valkey              string `mapstructure:"VALKEY"`
	SupabaseSecret      string `mapstructure:"SUPABASE_SECRET"`
	SupabaseIssuer      string `mapstructure:"SUPABASE_ISSUER"`
	GoogleBucket        string `mapstructure:"GOOGLE_BUCKET"`
	GoogleProject       string `mapstructure:"GOOGLE_PROJECT"`
	GoogleCredentials   string `mapstructure:"GOOGLE_CREDENTIALS"`
	GoogleMapsAPIKey    string `mapstructure:"GOOGLE_MAPS_API_KEY"`
	BunnyVideoKey       string `mapstructure:"BUNNY_VIDEO_KEY"`
	BunnyVideoLibraryID string `mapstructure:"BUNNY_VIDEO_LIBRARY_ID"`
	BunnyStorageKey     string `mapstructure:"BUNNY_STORAGE_KEY"`
	BunnyStorageZone    string `mapstructure:"BUNNY_STORAGE_ZONE"`
	BunnyStorageRegion  string `mapstructure:"BUNNY_STORAGE_REGION"`
	BunnyStorageCDN     string `mapstructure:"BUNNY_STORAGE_CDN"`
	SENDGRID            string `mapstructure:"SENDGRID"`
	GOOGLE_AI           string `mapstructure:"GOOGLE_AI"`
	POSTHOG             string `mapstructure:"POSTHOG"`
	TwilioAccountSID    string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken     string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioFromPhone     string `mapstructure:"TWILIO_FROM_PHONE"`
}

func SetEnvironment(env string) (cfg EnvironmentConfiguration, err error) {
	if env == "docker" {
		viper.SetConfigName("docker")
	} else if env == "prod" {
		viper.SetConfigName("prod")
	} else if env == "stag" {
		viper.SetConfigName("stag")
	} else {
		viper.SetConfigName("dev")
	}

	viper.AddConfigPath("./environment")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}
