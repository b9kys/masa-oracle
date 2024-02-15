package masa

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/viper"
)

var (
	KeyFileKey           = viper.GetString("KeyFileKey")
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
	PageSize             = viper.GetInt("PageSize")
	NodeBackupFileName   = viper.GetString("NodeBackupFileName")
	NodeBackupPath       = viper.GetString("NodeBackupPath")
	Version              = viper.GetString("Version")
	DefaultRPCURL        = viper.GetString("DefaultRPCURL")
	Environment          = viper.GetString("Environment")
)

func init() {
	viper.SetDefault("KeyFileKey", "private.key")
	viper.SetDefault("CertPem", "cert.pem")
	// Set defaults for all other variables similarly
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
