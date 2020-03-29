package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	patchSemver = "2"

	wholeString = majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}

func ascii() {
	fmt.Println("\033[1;32m")
	myFigure := figure.NewFigure("karai", "block", true)
	myFigure.Print()
	fmt.Println("\x1b[0m")
}

func main() {

	ascii()
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
			fmt.Printf("\033[1;32m%s\x1b[0m> ", str)
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

		fmt.Println("\033[0;37mType \033[1;32m'menu'\033[0;37m to view a list of commands\033[1;37m")
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
		} else if strings.Compare("license", text) == 0 {
			logrus.Debug("Displaying license")
			printLicense()
		} else if strings.Compare("create-wallet", text) == 0 {
			logrus.Debug("Creating Wallet")
			menuCreateWallet()
		} else if strings.Compare("open-wallet", text) == 0 {
			logrus.Debug("Opening Wallet")
			menuOpenWallet()
		} else if strings.Compare("transaction-history", text) == 0 {
			logrus.Debug("Opening Transaction History")
			menuGetContainerTransactions()
		} else if strings.Compare("open-wallet-info", text) == 0 {
			logrus.Debug("Opening Wallet Info")
			menuOpenWalletInfo()
		} else if strings.Compare("create-peer", text) == 0 {
			menuCreatePeer()
		} else if strings.Compare("exit", text) == 0 {
			logrus.Warning("Exiting")
			menuExit()
		} else if strings.Compare("quit", text) == 0 {
			logrus.Warning("Exiting")
			menuExit()
		} else if strings.Compare("close", text) == 0 {
			logrus.Warning("Exiting")
			menuExit()
		} else if strings.Compare("\n", text) == 0 {
			fmt.Println("")
		} else {
			fmt.Println("\nChoose an option from the menu")
			menuHelp()
		}

	}
}

func menuHelp() {
	fmt.Println("\n\033[1;32mWALLET_OPTIONS\033[1;37m\x1b[0m")
	fmt.Println("\033[1;37mopen-wallet \t\t \033[0;37mOpen a TRTL wallet\x1b[0m")
	fmt.Println("\033[1;37mopen-wallet-info \t \033[0;37mShow wallet and connection info\x1b[0m")
	fmt.Println("\033[1;37mcreate-wallet \t\t \033[0;37mCreate a TRTL wallet\x1b[0m")
	fmt.Println("\033[1;30mwallet-balance \t\t Displays wallet balance\x1b[0m")

	fmt.Println("\n\033[1;32mIPFS_OPTIONS\033[1;37m\x1b[0m")
	fmt.Println("\033[1;37mcreate-peer \t\t \033[0;37mCreates IPFS peer\x1b[0m")
	fmt.Println("\033[1;30mlist-servers \t\t Lists pinning servers\x1b[0m")

	fmt.Println("\n\033[1;32mGENERAL_OPTIONS\033[1;37m\x1b[0m")
	fmt.Println("\033[1;37mversion \t\t \033[0;37mDisplays version\033[0m")
	fmt.Println("\033[1;37mlicense \t\t \033[0;37mDisplays license\033[0m")
	fmt.Println("\033[1;37mexit \t\t\t \033[0;37mQuit immediately\x1b[0m")

	fmt.Println("")
}

func menuOpenWalletInfo() {
	walletInfoPrimaryAddressBalance()
	getNodeInfo()
	getWalletAPIStatus()
}

func menuGetContainerTransactions() {
	// logrus.Info("[Container transactions]")
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/transactions", nil)
	if err != nil {
		log.Fatal("Error reading request for transactions. ", err)
	}

	req.Header.Set("X-API-KEY", "pineapples")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response from transactions query. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body of transactions response. ", err)
	}

	fmt.Printf("%s\n", body)
}

func getWalletAPIStatus() {
	logrus.Info("[Wallet-API Status]")
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/status", nil)
	if err != nil {
		log.Fatal("Error reading request for wallet-api status. ", err)
	}

	req.Header.Set("X-API-KEY", "pineapples")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response from wallet-api status. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body of wallet-api status response. ", err)
	}

	fmt.Printf("%s\n", body)
}

func getNodeInfo() {
	logrus.Info("[Node Info]")
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/node", nil)
	if err != nil {
		log.Fatal("Error reading request for node info. ", err)
	}

	req.Header.Set("X-API-KEY", "pineapples")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response from node info. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body of node info response. ", err)
	}

	fmt.Printf("%s\n", body)
}

func walletInfoPrimaryAddressBalance() {
	logrus.Info("[Primary Address]")
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/balances", nil)
	if err != nil {
		log.Fatal("Error reading request for balances. ", err)
	}

	req.Header.Set("X-API-KEY", "pineapples")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response from balances. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body of balances response. ", err)
	}

	fmt.Printf("%s\n", body)
}

func printLicense() {
	fmt.Println("\n\033[1;32m" + appName + " \033[0;32mv" + semverInfo() + "\033[0;37m by \033[1;37m" + appDev)
	fmt.Println("\033[0;32m" + appRepository + "\n")
	fmt.Println("\033[1;37mMIT License\n\nCopyright (c) 2020-2021 RockSteady, TurtleCoin Developers\n\033[1;30mPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in allcopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.\n")
}

func menuCreateWallet() {
	logrus.Debug("Creating Wallet")
	url := "http://127.0.0.1:8070/wallet/create"

	data := []byte(`{"daemonHost": "127.0.0.1",	"daemonPort": 11898, "filename": "karai-wallet.wallet", "password": "supersecretpassword"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "pineapples")

	client := &http.Client{Timeout: time.Second * 10}
	logrus.Info(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	logrus.Info("response Status:", resp.Status)
	logrus.Info("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	fmt.Printf("%s\n", body)
}

func menuOpenWallet() {
	logrus.Debug("Opening Wallet")
	url := "http://127.0.0.1:8070/wallet/open"

	data := []byte(`{"daemonHost": "127.0.0.1",	"daemonPort": 11898, "filename": "karai-wallet.wallet", "password": "supersecretpassword"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "pineapples")

	client := &http.Client{Timeout: time.Second * 10}
	logrus.Info(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	logrus.Info("response Status:", resp.Status)
	logrus.Info("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	fmt.Printf("%s\n", body)

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
