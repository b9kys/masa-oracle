package masa

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/viper"
)

var (
	KeyFileKey           = viper.GetString("KeyFileKey")
	KeyFilePath          = viper.GetString("KeyFilePath")
	CertPem              = viper.GetString("CertPem")
	Cert                 = viper.GetString("Cert")
	Peers                = viper.GetString("Peers")
	Bootnodes            = viper.GetString("Bootnodes")
	masaPrefix           = viper.GetString("masaPrefix")
	oracleProtocol       = viper.GetString("oracleProtocol")
	NodeDataSyncProtocol = viper.GetString("NodeDataSyncProtocol")
	NodeGossipTopic      = viper.GetString("NodeGossipTopic")
	AdTopic              = viper.GetString("AdTopic")
	rendezvous           = viper.GetString("rendezvous")
	PortNbr              = viper.GetString("PortNbr")
	PageSize             = viper.GetFloat64("PageSize")
	NodeBackupFileName   = viper.GetString("NodeBackupFileName")
	NodeBackupPath       = viper.GetString("NodeBackupPath")
	Version              = viper.GetString("Version")
	DefaultRPCURL        = viper.GetString("DefaultRPCURL")
	Environment          = viper.GetString("Environment")
)

func init() {
	viper.SetDefault("KeyFilePath", "/home/masa/.masa/")
	viper.SetDefault("KeyFileKey", "masa_oracle_key")
	viper.SetDefault("CertPem", "cert.pem")
	viper.SetDefault("PageSize", 25)
	viper.SetDefault("NodeBackupFileName", "node-backup.json")
	viper.SetDefault("NodeBackupPath", "/home/masa/.masa") // This is the default path where the node backup file will be stored
	viper.SetDefault("Version", "v1.0.0")                  // This is the default version of the protocol
	viper.SetDefault("Environment", "test")
	viper.SetDefault("DefaultRPCURL", "http://localhost:8545")
}

func ProtocolWithVersion(protocolName string) protocol.ID {
	if Environment == "" {
		return protocol.ID(fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version))
	}
	return protocol.ID(fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, Version, Environment))
}

func TopicWithVersion(protocolName string) string {
	if Environment == "" {
		return fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version)
	}
	return fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, Version, Environment)
}
