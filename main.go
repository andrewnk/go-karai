package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/libp2p/go-libp2p"
	autonat "github.com/libp2p/go-libp2p-autonat-svc"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pquic "github.com/libp2p/go-libp2p-quic-transport"
	routing "github.com/libp2p/go-libp2p-routing"
	secio "github.com/libp2p/go-libp2p-secio"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/sirupsen/logrus"
)

const appName = "go-karai"
const appDev = "TurtleCoin Developers"
const appDescription = appName + " is helper software for the Karai network"
const appRepository = "https://github.com/rocksteadytc/go-karai"

func semverInfo() string {
	var majorSemver, minorSemver, patchSemver, wholeString string
	majorSemver = "0"
	minorSemver = "0"
	patchSemver = "1"

	wholeString = majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}

func ascii() {
	myFigure := figure.NewFigure("karai", "shadow", true)
	myFigure.Print()
}

func main() {

	ascii()
	fmt.Println("\nType \x1b[35m'menu'\x1b[0m to view a list of commands")
	inputHandler()

}
func handleStream(s network.Stream) {
	logrus.Debug("New stream")
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readData(rw)
	go writeData(rw)
}
func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[35m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}

}
func inputHandler() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("help", text) == 0 {
			menuHelp()
		} else if strings.Compare("?", text) == 0 {
			menuHelp()
		} else if strings.Compare("menu", text) == 0 {
			menuHelp()
		} else if strings.Compare("version", text) == 0 {
			logrus.Debug("Displaying version")
			menuVersion()
		} else if strings.Compare("peer-info", text) == 0 {
			menuCreatePeer()
		} else if strings.Compare("exit", text) == 0 {
			menuExit()
		} else if strings.Compare("quit", text) == 0 {
			menuExit()
		} else if strings.Compare("close", text) == 0 {
			menuExit()
		} else {
			logrus.Error("\nwtf is " + text + "???")
			logrus.Error("Please choose something I can actually do:")
			menuHelp()
		}

	}
}

func menuHelp() {
	fmt.Println("\n\x1b[35mversion \t\t \x1b[0mDisplays version")
	fmt.Println("\x1b[35mcreate-wallet \t\t \x1b[0mCreate a TRTL wallet")
	fmt.Println("\x1b[35mwallet-balance \t\t \x1b[0mDisplays wallet balance")
	fmt.Println("\x1b[35mlist-servers \t\t \x1b[0mLists pinning servers")
	fmt.Println("\x1b[35mcreate-peer \t\t \x1b[0mCreates IPFS peer")
	fmt.Println("\x1b[35mexit \t\t\t \x1b[0mQuit immediately")
}
func menuCreateWallet() {
	fmt.Println("create a wallet")
}
func menuBalance() {
	fmt.Println("display wallet balance")
}
func menuListPinServers() {
	fmt.Println("list known pinning servers")
}
func menuCreatePeer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, -1,
	)
	if err != nil {
		panic(err)
	}

	var idht *dht.IpfsDHT

	nodePeer, err := libp2p.New(ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/9000",
			"/ip4/0.0.0.0/udp/9000/quic",
		),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(secio.ID, secio.New),
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.DefaultTransports,
		libp2p.ConnectionManager(connmgr.NewConnManager(
			100,         // Lowwater
			400,         // HighWater,
			time.Minute, // GracePeriod
		)),
		libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		libp2p.EnableAutoRelay(),
	)
	if err != nil {
		panic(err)
	}
	_, err = autonat.NewAutoNATService(ctx, nodePeer,
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(secio.ID, secio.New),
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.DefaultTransports,
	)

	for _, addr := range dht.DefaultBootstrapPeers {
		pi, _ := peer.AddrInfoFromP2pAddr(addr)
		nodePeer.Connect(ctx, *pi)
	}

	fmt.Printf("Peer ID is %s\n", nodePeer.ID())
}

func menuVersion() {
	fmt.Println(appName + " - v" + semverInfo())
}

func menuExit() {
	// if err := node.Close(); err != nil {
	// 	panic(err)
	// }
	fmt.Println("\nExiting!")
	os.Exit(0)
}
