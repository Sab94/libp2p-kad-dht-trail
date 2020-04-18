package main

import (
	"context"
	"fmt"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()

	sourceMultiAddr, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/4000")
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("This node is ", host.ID().String())
	fmt.Println("Available multiaddrs :")
	for _, v := range host.Addrs() {
		fmt.Printf("%s/p2p/%s\n", v, host.ID().String())
	}

	_, err = dht.New(ctx, host, []dht.Option{
		dht.ProtocolPrefix("/dht-trail"),
		dht.Mode(dht.ModeAutoServer),
		dht.RoutingTableRefreshPeriod(time.Minute * 1),
	}...)
	if err != nil {
		panic(err)
	}

	select {}
}