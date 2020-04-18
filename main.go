package main

import (
	"context"
	"flag"
	"time"

	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	h "github.com/libp2p/go-libp2p-core/host"
	pnet "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
)
var dhtGlobal *dht.IpfsDHT
var hostGlobal h.Host
var log = logging.Logger("dht/sampleOne")
var port string
var bootstrapAddr string

func main() {
	logging.SetLogLevel("dht/sampleOne", "Debug")

	flag.StringVar(&port, "port", "8000", "--port")
	flag.StringVar(&bootstrapAddr, "bootstrap", "-", "--bootstrap")
	flag.Parse()

	if bootstrapAddr == "-" {
		log.Error("Please provide a valid bootstrap address")
		return
	}

	ctx := context.Background()

	opts := []libp2p.Option {
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/"+port),
	}
	host, err := libp2p.New(ctx, opts...)
	if err != nil {
		log.Fatal("Unable to create libp2p host :", err.Error())
		return
	}
	log.Debug("Host Created : ", host.Addrs(), host.ID().Pretty())
	host.Network().SetConnHandler(connHandler)

	d, err := dht.New(ctx, host, []dht.Option{
		dht.ProtocolPrefix("/dht-trail"),
		dht.Mode(dht.ModeAutoServer),
		dht.RoutingTableRefreshPeriod(time.Minute * 1),
	}...)
	if err != nil {
		log.Fatal("Unable to create dht :", err.Error())
		return
	}

	bootstrap, err := multiaddr.NewMultiaddr(bootstrapAddr)
	if err != nil {
		log.Fatal("Invalid bootstrap address :", err.Error())
		return
	}
	peerinfo, _ := peer.AddrInfoFromP2pAddr(bootstrap)
	ctx, _ = context.WithTimeout(context.Background(), time.Second*10)
	if err := host.Connect(ctx, *peerinfo); err != nil {
		log.Error("Could not connect to bootstrap")
	} else {
		log.Info("Connection established with bootstrap node:", *peerinfo)
	}

	dhtGlobal = d
	hostGlobal = host
	select{}
}

func connHandler(c pnet.Conn) {
	log.Debug("Handle connection request")
	v1b := cid.V1Builder{Codec: cid.Raw, MhType: mh.SHA2_256}
	rendezvousPoint, _ := v1b.Sum([]byte("10880 Malibu Point, 90265"))

	<-time.After(time.Second * 5)

	log.Debug("Announcing myself with : ", rendezvousPoint)
	tctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	if err := dhtGlobal.Provide(tctx, rendezvousPoint, true); err != nil {
		log.Fatal("Unable to provide :", err.Error())
		return
	}

	log.Debug("Looking for other peers...")
	tctx, _ = context.WithTimeout(context.Background(), time.Second*10)
	providers, err := dhtGlobal.FindProviders(tctx, rendezvousPoint)
	if err != nil {
		log.Fatal("Unable to find provider :", err.Error())
		return
	}

	if len(providers) != 0 {
		for _, v := range providers {
			if v.ID != hostGlobal.ID() {
				tctx, _ = context.WithTimeout(context.Background(), time.Second*10)
				err := hostGlobal.Connect(tctx, v)
				if err == nil {
					log.Debug("Connection established with : ", v.ID.String())
					log.Debugf("Remote node %s is providing %s\n", v.ID.Pretty(), rendezvousPoint)
				} else {
					log.Error("Unable to establish connection with : ", v.ID.String(), err.Error())
				}
			}
		}
	} else {
		log.Debug("no remote providers!\n")
	}
}