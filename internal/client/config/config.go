package config

// Config client settings
type Config struct {
	GRPC       GRPCConfig       `mapstructure:"grpc"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
}

// GRPCConfig  gRPC settings
type GRPCConfig struct {
	Address string `mapstructure:"address"`
}

// EncryptionConfig cypher settings
type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}
