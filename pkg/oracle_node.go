package masa

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/ad"
	crypto2 "github.com/masa-finance/masa-oracle/pkg/crypto"
	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
)

type OracleNode struct {
	Host          host.Host
	PrivKey       *ecdsa.PrivateKey
	Protocol      protocol.ID
	priorityAddrs multiaddr.Multiaddr
	multiAddrs    []multiaddr.Multiaddr
	DHT           *dht.IpfsDHT
	Context       context.Context
	PeerChan      chan myNetwork.PeerEvent
	NodeTracker   *pubsub2.NodeEventTracker
	PubSubManager *pubsub2.Manager
	Signature     string
	IsStaked      bool
	StartTime     time.Time
	IDService     identify.IDService
}

func (node *OracleNode) GetMultiAddrs() multiaddr.Multiaddr {
	if node.priorityAddrs == nil {
		pAddr := myNetwork.GetPriorityAddress(node.multiAddrs)
		node.priorityAddrs = pAddr
	}
	return node.priorityAddrs
}

func NewOracleNode(ctx context.Context, privKey crypto.PrivKey, portNbr int, useUdp, useTcp bool, isStaked bool) (*OracleNode, error) {
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	resourceManager, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, err
	}

	var addrStr []string
	libp2pOptions := []libp2p.Option{
		libp2p.Identity(privKey),
		libp2p.ResourceManager(resourceManager),
		libp2p.Ping(false), // disable built-in ping
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(), // Enable Circuit Relay v2 with hop
	}

	securityOptions := []libp2p.Option{
		libp2p.Security(noise.ID, noise.New),
	}
	if useUdp {
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", portNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(quic.NewTransport))
	}
	if useTcp {
		securityOptions = append(securityOptions, libp2p.Security(libp2ptls.ID, libp2ptls.New))
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", portNbr))
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

	// Create a new Identify service
	ids, err := identify.NewIDService(hst)
	if err != nil {
		return nil, err
	}

	ecdsaPrivKey, err := crypto2.Libp2pPrivateKeyToEcdsa(privKey)
	if err != nil {
		return nil, err
	}
	return &OracleNode{
		Host:          hst,
		PrivKey:       ecdsaPrivKey,
		Protocol:      oracleProtocol,
		multiAddrs:    myNetwork.GetMultiAddressesForHostQuiet(hst),
		Context:       ctx,
		PeerChan:      make(chan myNetwork.PeerEvent),
		NodeTracker:   pubsub2.NewNodeEventTracker(),
		PubSubManager: subscriptionManager,
		IsStaked:      isStaked,
		IDService:     ids,
	}, nil
}

func (node *OracleNode) Start() (err error) {
	logrus.Infof("Starting node with ID: %s", node.GetMultiAddrs().String())
	node.Host.SetStreamHandler(node.Protocol, node.handleStream)
	node.Host.SetStreamHandler(NodeDataSyncProtocol, node.ReceiveNodeData)
	node.Host.SetStreamHandler(NodeGossipTopic, node.GossipNodeData)

	node.Host.Network().Notify(node.NodeTracker)

	go node.ListenToNodeTracker()
	go node.handleDiscoveredPeers()

	err = myNetwork.WithMDNS(node.Host, rendezvous, node.PeerChan)
	if err != nil {
		return err
	}

	bootNodeAddrs, err := myNetwork.GetBootNodesMultiAddress(os.Getenv(Peers))
	if err != nil {
		return err
	}

	node.DHT, err = myNetwork.WithDht(node.Context, node.Host, bootNodeAddrs, oracleProtocol, masaPrefix, node.PeerChan)
	if err != nil {
		return err
	}

	go myNetwork.Discover(node.Context, node.Host, node.DHT, node.Protocol, node.GetMultiAddrs())

	// Subscribe to a topics
	err = node.PubSubManager.AddSubscription(NodeGossipTopic, node.NodeTracker)
	if err != nil {
		return err
	}
	err = node.PubSubManager.AddSubscription(AdTopic, &ad.SubscriptionHandler{})
	node.StartTime = time.Now()
	return nil
}

func (node *OracleNode) handleDiscoveredPeers() {
	for {
		select {
		case peer := <-node.PeerChan: // will block until we discover a peer
			logrus.Info("Peer Event for:", peer, ", Action:", peer.Action)

			if err := node.Host.Connect(node.Context, peer.AddrInfo); err != nil {
				logrus.Error("Connection failed:", err)
				continue
			}

			// open a stream, this stream will be handled by handleStream other end
			stream, err := node.Host.NewStream(node.Context, peer.AddrInfo.ID, node.Protocol)

			if err != nil {
				logrus.Error("Stream open failed", err)
			} else {
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

				go node.writeData(rw, peer, stream)
				go node.readData(rw, peer, stream)
			}
		case <-node.Context.Done():
			return
		}
	}
}

func (node *OracleNode) handleStream(stream network.Stream) {
	data := node.handleStreamData(stream)
	logrus.Info("handleStream -> Received data:", string(data))
	remotePeer := stream.Conn().RemotePeer()

	// Wait for the Identify protocol to complete
	<-node.IDService.IdentifyWait(stream.Conn())

	// Now the public key should be available in the Peerstore
	pubKey := node.Host.Peerstore().PubKey(remotePeer)
	if pubKey == nil {
		logrus.Warnf("No public key found for peer %s", remotePeer)
	} else {
		logrus.Infof("Public key found for peer %s", remotePeer.String())
	}
	//data := node.handleStreamData(stream)
	//logrus.Info("handleStream -> Received data:", string(data))
}

func (node *OracleNode) readData(rw *bufio.ReadWriter, event myNetwork.PeerEvent, stream network.Stream) {
	defer func() {
		err := stream.Close()
		if err != nil {
			logrus.Error("Error closing stream:", err)
		}
	}()

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			logrus.Error("Error reading from buffer:", err)
			return
		}
		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

func (node *OracleNode) writeData(rw *bufio.ReadWriter, event myNetwork.PeerEvent, stream network.Stream) {
	defer func() {
		err := stream.Close()
		if err != nil {
			logrus.Error("Error closing stream:", err)
		}
	}()

	for {
		// Generate a message including the multiaddress of the sender
		sendData := fmt.Sprintf("%s: Hello from %s\n", event.Source, node.GetMultiAddrs().String())

		_, err := rw.WriteString(sendData)
		if err != nil {
			logrus.Error("Error writing to buffer:", err)
			return
		}
		err = rw.Flush()
		if err != nil {
			logrus.Error("Error flushing buffer:", err)
			return
		}
		// Sleep for a while before sending the next message
		time.Sleep(time.Second * 30)
	}
}

func (node *OracleNode) IsPublisher() bool {
	// Node is a publisher if it has a non-empty signature
	return node.Signature != ""
}
