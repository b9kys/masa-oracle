package masa

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/chyeh/pubip"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	chain "github.com/masa-finance/masa-oracle/pkg/chain"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
)

type OracleNode struct {
	Host             host.Host
	PrivKey          *ecdsa.PrivateKey
	Protocol         protocol.ID
	priorityAddrs    multiaddr.Multiaddr
	multiAddrs       []multiaddr.Multiaddr
	DHT              *dht.IpfsDHT
	Context          context.Context
	PeerChan         chan myNetwork.PeerEvent
	NodeTracker      *pubsub2.NodeEventTracker
	PubSubManager    *pubsub2.Manager
	Signature        string
	IsStaked         bool
	IsValidator      bool
	IsTwitterScraper bool
	IsDiscordScraper bool
	IsWebScraper     bool
	IsLlmServer      bool
	StartTime        time.Time
	WorkerTracker    *pubsub2.WorkerEventTracker
	BlockTracker     *pubsub2.BlockEventTracker
	ActorEngine      *actor.RootContext
	ActorRemote      *remote.Remote
	Blockchain       *chain.Chain
}

// GetMultiAddrs returns the priority multiaddr for this node.
// It first checks if the priority address is already set, and returns it if so.
// If not, it determines the priority address from the available multiaddrs using
// the GetPriorityAddress utility function, sets it, and returns it.
func (node *OracleNode) GetMultiAddrs() multiaddr.Multiaddr {
	if node.priorityAddrs == nil {
		pAddr := myNetwork.GetPriorityAddress(node.multiAddrs)
		node.priorityAddrs = pAddr
	}
	return node.priorityAddrs
}

// getOutboundIP is a function that returns the outbound IP address of the current machine as a string.
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("Error getting outbound IP")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

// NewOracleNode creates a new OracleNode instance with the provided context and
// staking status. It initializes the libp2p host, DHT, pubsub manager, and other
// components needed for an Oracle node to join the network and participate.
func NewOracleNode(ctx context.Context, isStaked bool) (*OracleNode, error) {
	// Start with the default scaling limits.
	cfg := config.GetInstance()
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	resourceManager, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, err
	}

	var addrStr []string
	libp2pOptions := []libp2p.Option{
		libp2p.Identity(masacrypto.KeyManagerInstance().Libp2pPrivKey),
		libp2p.ResourceManager(resourceManager),
		libp2p.Ping(false), // disable built-in ping
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(), // Enable Circuit Relay v2 with hop
	}

	securityOptions := []libp2p.Option{
		libp2p.Security(noise.ID, noise.New),
	}
	// @note fix for increase buffer size warning on linux
	// sudo sysctl -w net.core.rmem_max=7500000
	// sudo sysctl -w net.core.wmem_max=7500000
	if cfg.UDP {
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", cfg.PortNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(quic.NewTransport))
	}
	if cfg.TCP {
		securityOptions = append(securityOptions, libp2p.Security(libp2ptls.ID, libp2ptls.New))
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.PortNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(tcp.NewTCPTransport))
		libp2pOptions = append(libp2pOptions, libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport))
	}
	libp2pOptions = append(libp2pOptions, libp2p.ChainOptions(securityOptions...))
	libp2pOptions = append(libp2pOptions, libp2p.ListenAddrStrings(addrStr...))

	hst, err := libp2p.New(libp2pOptions...)
	if err != nil {
		return nil, err
	}

	subscriptionManager, err := pubsub2.NewPubSubManager(ctx, hst)
	if err != nil {
		return nil, err
	}

	isValidator, _ := strconv.ParseBool(cfg.Validator)
	isTwitterScraper := cfg.TwitterScraper
	isDiscordScraper := cfg.DiscordScraper
	isWebScraper := cfg.WebScraper

	system := actor.NewActorSystemWithConfig(actor.Configure(
		actor.ConfigOption(func(config *actor.Config) {
			config.LoggerFactory = func(system *actor.ActorSystem) *slog.Logger {
				return slog.New(slog.NewTextHandler(io.Discard, nil))
			}
		}),
	))
	engine := system.Root

	var ip any
	if cfg.Environment == "local" {
		ip = getOutboundIP()
	} else {
		ip, _ = pubip.Get()
	}
	conf := remote.Configure("0.0.0.0", 4001,
		remote.WithAdvertisedHost(fmt.Sprintf("%s:4001", ip)))

	r := remote.NewRemote(system, conf)
	go r.Start()

	return &OracleNode{
		Host:             hst,
		PrivKey:          masacrypto.KeyManagerInstance().EcdsaPrivKey,
		Protocol:         config.ProtocolWithVersion(config.OracleProtocol),
		multiAddrs:       myNetwork.GetMultiAddressesForHostQuiet(hst),
		Context:          ctx,
		PeerChan:         make(chan myNetwork.PeerEvent),
		NodeTracker:      pubsub2.NewNodeEventTracker(config.Version, cfg.Environment),
		PubSubManager:    subscriptionManager,
		IsStaked:         isStaked,
		IsValidator:      isValidator,
		IsTwitterScraper: isTwitterScraper,
		IsDiscordScraper: isDiscordScraper,
		IsWebScraper:     isWebScraper,
		IsLlmServer:      cfg.LlmServer,
		ActorEngine:      engine,
		ActorRemote:      r,
		Blockchain:       &chain.Chain{},
	}, nil
}

// Start initializes the OracleNode by setting up libp2p stream handlers,
// connecting to the DHT and bootnodes, and subscribing to topics. It launches
// goroutines to handle discovered peers, listen to the node tracker, and
// discover peers. If this is a bootnode, it adds itself to the node tracker.
func (node *OracleNode) Start() (err error) {
	logrus.Infof("Starting node with ID: %s", node.GetMultiAddrs().String())

	bootNodeAddrs, err := myNetwork.GetBootNodesMultiAddress(config.GetInstance().Bootnodes)
	if err != nil {
		return err
	}

	node.Host.SetStreamHandler(node.Protocol, node.handleStream)
	node.Host.SetStreamHandler(config.ProtocolWithVersion(config.NodeDataSyncProtocol), node.ReceiveNodeData)
	// IsStaked
	if node.IsStaked {
		node.Host.SetStreamHandler(config.ProtocolWithVersion(config.NodeGossipTopic), node.GossipNodeData)
	}
	node.Host.Network().Notify(node.NodeTracker)

	go node.ListenToNodeTracker()
	go node.handleDiscoveredPeers()

	node.DHT, err = myNetwork.WithDht(node.Context, node.Host, bootNodeAddrs, node.Protocol, config.MasaPrefix, node.PeerChan, node.IsStaked)
	if err != nil {
		return err
	}
	err = myNetwork.WithMDNS(node.Host, config.Rendezvous, node.PeerChan)
	if err != nil {
		return err
	}

	go myNetwork.Discover(node.Context, node.Host, node.DHT, node.Protocol)

	nodeData := node.NodeTracker.GetNodeData(node.Host.ID().String())
	if nodeData == nil {
		publicKeyHex := masacrypto.KeyManagerInstance().EthAddress
		nodeData = pubsub2.NewNodeData(node.GetMultiAddrs(), node.Host.ID(), publicKeyHex, pubsub2.ActivityJoined)
		nodeData.IsStaked = node.IsStaked
		nodeData.SelfIdentified = true
	}

	cfg := config.GetInstance()
	nodeData.IsDiscordScraper = cfg.DiscordScraper
	nodeData.IsTwitterScraper = cfg.TwitterScraper
	nodeData.IsWebScraper = cfg.WebScraper
	nodeData.IsValidator = cfg.Validator == "true"

	nodeData.Joined()
	node.NodeTracker.HandleNodeData(*nodeData)

	// call SubscribeToTopics on startup
	if err := SubscribeToTopics(node); err != nil {
		return err
	}
	node.StartTime = time.Now()

	return nil
}

// handleDiscoveredPeers listens on the PeerChan for discovered peers from the
// network discovery routines. It handles connecting to new peers and closing
// connections to peers that disconnect. This runs continuously to handle
// discovered peers.
func (node *OracleNode) handleDiscoveredPeers() {
	for {
		select {
		case peer := <-node.PeerChan: // will block until we discover a peer
			logrus.Debugf("Peer Event for: %s, Action: %s", peer.AddrInfo.ID.String(), peer.Action)
			// If the peer is a new peer, connect to it
			if peer.Action == myNetwork.PeerAdded {
				if err := node.Host.Connect(node.Context, peer.AddrInfo); err != nil {
					logrus.Errorf("Connection failed for peer: %s %v", peer.AddrInfo.ID.String(), err)
					// close the connection
					err := node.Host.Network().ClosePeer(peer.AddrInfo.ID)
					if err != nil {
						logrus.Error(err)
					}
					continue
				}
			}
		case <-node.Context.Done():
			return
		}
	}
}

// handleStream handles an incoming libp2p stream from a remote peer.
// It reads the stream data, validates the remote peer ID, updates the node tracker
// with the remote peer's information, and logs the event.
func (node *OracleNode) handleStream(stream network.Stream) {
	remotePeer, nodeData, err := node.handleStreamData(stream)
	if err != nil {
		if strings.HasPrefix(err.Error(), "un-staked") {
			// just ignore the error
			return
		}
		logrus.Errorf("Failed to read stream: %v", err)
		return
	}
	if remotePeer.String() != nodeData.PeerId.String() {
		logrus.Warnf("Received data from unexpected peer %s", remotePeer)
		return
	}
	multiAddr := stream.Conn().RemoteMultiaddr()
	newNodeData := pubsub2.NewNodeData(multiAddr, remotePeer, nodeData.EthAddress, pubsub2.ActivityJoined)
	newNodeData.IsStaked = nodeData.IsStaked
	err = node.NodeTracker.AddOrUpdateNodeData(newNodeData, false)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("nodeStream -> Received data from: %s", remotePeer.String())
}

// IsWorker determines if the OracleNode is configured to act as an actor.
// An actor node is one that has at least one of the following scrapers enabled:
// TwitterScraper, DiscordScraper, or WebScraper.
// It returns true if any of these scrapers are enabled, otherwise false.
func (node *OracleNode) IsWorker() bool {
	// need to get this by node data
	cfg := config.GetInstance()
	if cfg.TwitterScraper || cfg.DiscordScraper || cfg.WebScraper {
		return true
	}
	return false
}

// IsPublisher returns true if this node is a publisher node.
// A publisher node is one that has a non-empty signature.
func (node *OracleNode) IsPublisher() bool {
	// Node is a publisher if it has a non-empty signature
	return node.Signature != ""
}

// FromUnixTime converts a Unix timestamp into a formatted string.
// The Unix timestamp is expected to be in seconds.
// The returned string is in the format "2006-01-02T15:04:05.000Z".
func (node *OracleNode) FromUnixTime(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02T15:04:05.000Z")
}

// ToUnixTime converts a formatted string time into a Unix timestamp.
// The input string is expected to be in the format "2006-01-02T15:04:05.000Z".
// The returned Unix timestamp is in seconds.
func (node *OracleNode) ToUnixTime(stringTime string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", stringTime)
	return t.Unix()
}

// Version returns the current version string of the oracle node software.
func (node *OracleNode) Version() string {
	return config.Version
}

// LogActiveTopics logs the currently active topic names to the
// default logger. It gets the list of active topics from the
// PubSubManager and logs them if there are any, otherwise it logs
// that there are no active topics.
func (node *OracleNode) LogActiveTopics() {
	topicNames := node.PubSubManager.GetTopicNames()
	if len(topicNames) > 0 {
		logrus.Infof("Active topics: %v", topicNames)
	} else {
		logrus.Info("No active topics.")
	}
}

// Blockchain Implementation
var (
	blocksCh = make(chan *pubsub.Message)
)

// SubscribeToBlocks is a function that takes in a context and an OracleNode as parameters.
// It is used to subscribe the given OracleNode to the blockchain blocks.
func SubscribeToBlocks(ctx context.Context, node *OracleNode) {
	node.BlockTracker = &pubsub2.BlockEventTracker{BlocksCh: blocksCh}
	err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.BlockTopic), node.BlockTracker, true)
	if err != nil {
		logrus.Errorf("Subscribe error %v", err)
	}

	if !node.IsValidator {
		return
	}

	go node.Blockchain.Init()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logrus.Debug("tick")
		case block := <-node.BlockTracker.BlocksCh:
			_ = node.Blockchain.AddBlock(block.Data)
			if node.Blockchain.LastHash != nil {
				b, e := node.Blockchain.GetBlock(node.Blockchain.LastHash)
				if e != nil {
					logrus.Errorf("Blockchain.GetBlock err: %v", e)
				}
				b.Print()
			}

		case <-ctx.Done():
			return
		}
	}
}
