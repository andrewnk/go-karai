package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/common-nighthawk/go-figure"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
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
	ctx := context.Background()

	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	logrus.Info("karai node address:", addrs[0])

	if len(os.Args) > 1 {
		addr, err := multiaddr.NewMultiaddr(os.Args[1])
		if err != nil {
			panic(err)
		}
		peer, err := peerstore.AddrInfoFromP2pAddr(addr)
		if err != nil {
			panic(err)
		}
		if err := node.Connect(ctx, *peer); err != nil {
			panic(err)
		}
		logrus.Info("sending ping message to", addr)
		ch := pingService.Ping(ctx, peer.ID)
		for i := 0; i < 1; i++ {
			res := <-ch
			logrus.Info("pinged", addr, "in", res.RTT)
		}
	} else {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		logrus.Fatal("Received signal, shutting down... ")
	}

	if err := node.Close(); err != nil {
		panic(err)
	}
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
			logrus.Debug("Displaying peer-info")
			menuPeerInfo()
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
	fmt.Println("\x1b[35mpeer-info \t\t \x1b[0mDisplays IPFS peer address")
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
func menuPeerInfo() {

	logrus.Info("\nThis is your IPFS node address:")
	// fmt.Println(node.Addrs())
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
