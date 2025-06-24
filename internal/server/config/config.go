package config

import (
	"github.com/Painkiller675/gophkeeper/pkg/cipher/aes/gcm"
	"time"
)

// Config server settings
type Config struct {
	GRPC       GRPCConfig       `mapstructure:"grpc"`
	DB         StorageConfig    `mapstructure:"db"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Hash       HashConfig       `mapstructure:"hasher"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
}

var BlockC *gcm.Cipher

// EncryptionConfig cypher settings
type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}

// GRPCConfig  some gRPC settings
type GRPCConfig struct {
	Address string `mapstructure:"address"`
}

// StorageConfig server database settings
type StorageConfig struct {
	URL string `mapstructure:"url"`
}

// AuthConfig authentication settings
type AuthConfig struct {
	Key            string        `mapstructure:"key"`
	ExpirationTime time.Duration `mapstructure:"expiration_time"`
}

// HashConfig hashing settings
type HashConfig struct {
	Key string `mapstructure:"key"`
}
