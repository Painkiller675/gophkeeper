package config

// Config client settings
type Config struct {
	GRPC GRPCConfig `mapstructure:"grpc"`
}

// GRPCConfig  gRPC settings
type GRPCConfig struct {
	Address string `mapstructure:"address"`
}
