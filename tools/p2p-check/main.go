package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	peer "github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
	"gitlab.com/thorchain/thornode/common/cosmos"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

func extractIPAddress(s string) (string, error) {
	parts := strings.Split(s, "/")
	if len(parts) < 3 || parts[1] != "ip4" {
		return "", fmt.Errorf("invalid format")
	}
	return parts[2], nil
}

func getPeers(peers map[string]string, source string) (map[string]string, error) {
	ctx := context.Background()

	// Create a new libp2p Host that listens on a random TCP port
	host, err := libp2p.New(ctx)
	if err != nil {
		return peers, err
	}

	// Create a DHT client (or a full DHT)
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		return peers, err
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	err = kademliaDHT.Bootstrap(ctx)
	if err != nil {
		return peers, err
	}

	// Connect to a known peer
	peerAddr, _ := ma.NewMultiaddr(source)
	peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
	if err = host.Connect(ctx, *peerInfo); err != nil {
		return peers, err
	}

	fmt.Println("Connected to:", peerInfo.ID)

	// Let's wait a bit to allow the DHT to do its magic and discover more peers
	time.Sleep(5 * time.Second)

	// Fetch the address book (peerstore)
	ps := host.Peerstore()
	for _, peerID := range ps.Peers() {
		addrs := ps.Addrs(peerID)
		for _, addr := range addrs {
			if strings.Contains(addr.String(), "5040") {
				peers[peerID.String()] = addr.String()
			}
		}
	}
	return peers, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func main() {
	var err error

	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount("thor", "thorpub")
	config.SetBech32PrefixForValidator("thorv", "thorvpub")
	config.SetBech32PrefixForConsensusNode("thorc", "thorcpub")
	config.Seal()

	nodesByPeerID := make(map[string]string)
	nodes := make([]openapi.Node, 0)
	url := fmt.Sprintf("%s/thorchain/nodes", getEnvOrDefault("THORNODE", "https://thornode.ninerealms.com"))
	// nolint
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("fail to get thornode status", err.Error())
	} else {
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
			fmt.Println("fail to decode thornode status", err.Error())
		} else {
			for _, node := range nodes {
				if node.PreflightStatus.Status == types.NodeStatus_Ready.String() {
					nodesByPeerID[node.PeerId] = ""
				}
			}
		}
	}

	fmt.Println("Discovering IP addresses for nodes...")
	for _, node := range nodes {
		if node.PreflightStatus.Status != types.NodeStatus_Ready.String() {
			continue
		}
		if nodesByPeerID[node.PeerId] != "" {
			continue
		}
		peers := make(map[string]string)
		peers, err = getPeers(peers, fmt.Sprintf("/ip4/%s/tcp/5040/p2p/%s", node.IpAddress, node.PeerId))
		if err != nil {
			fmt.Println("Error", err.Error())
		}
		for k, v := range peers {
			var ip string
			ip, err = extractIPAddress(v)
			if err != nil {
				fmt.Println("Error", err.Error())
				continue
			}
			nodesByPeerID[k] = ip
		}

		// see if we got all of our IPs
		var remaining int
		for _, v := range nodesByPeerID {
			if v == "" {
				remaining += 1
			}
		}
		fmt.Println("Remaining:", remaining)
		if remaining == 0 {
			break
		}
	}

	fmt.Println("\nTesting p2p health...")

	// ping and http get /p2pid on all peers
	height, err := fetchBlockHeight(getEnvOrDefault("RPC", "https://rpc.ninerealms.com"))
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	type rslt struct {
		node        openapi.Node
		description string
	}
	results := make(map[string]rslt)
	for peerID, ipAddr := range nodesByPeerID {
		if ipAddr == "" {
			continue
		}
		wg.Add(1)
		go func(peerID, ipAddr string) {
			defer func() {
				wg.Done()
			}()

			node := fetchNode(peerID, nodes)
			if node.PreflightStatus.Status != types.NodeStatus_Ready.String() {
				return
			}
			r := rslt{
				node: node,
			}

			// get the peer id
			_, err = fetchPeerId(ipAddr, peerID)
			if err != nil {
				r.description = err.Error()
				results[node.NodeAddress] = r
				return
			}

			// check the p2p port
			ok := checkPortOpen(ipAddr, 5040)
			if !ok {
				r.description = "port 5040 not open"
				results[node.NodeAddress] = r
				return
			}

			// check the thornode port
			ok = checkPortOpen(ipAddr, 27147)
			if !ok {
				r.description = "port 27147 not open"
				results[node.NodeAddress] = r
				return
			}

			// check thornode sync at tip
			var nodeHeight int64
			nodeHeight, err = fetchBlockHeight(fmt.Sprintf("http://%s:27147", ipAddr))
			if err != nil {
				r.description = err.Error()
				results[node.NodeAddress] = r
				return
			}
			if height > nodeHeight+5 {
				r.description = fmt.Sprintf("thornode behind the tip %d", height-nodeHeight)
				results[node.NodeAddress] = r
				return
			}
		}(peerID, ipAddr)
	}
	wg.Wait()

	if len(results) > 0 {
		fmt.Println("Node\t", "Operator\t", "Status\t", "Error\t")
		for k, v := range results {
			fmt.Println(k[len(k)-4:], "\t", v.node.NodeOperatorAddress[len(k)-4:], "\t\t", v.node.Status, "\t", v.description)
		}
		fmt.Println("Count failing:", len(results))
	} else {
		fmt.Println("OK")
	}
}

func fetchNode(peerID string, nodes []openapi.Node) openapi.Node {
	for _, node := range nodes {
		if node.PeerId == peerID {
			return node
		}
	}
	return openapi.Node{}
}

func fetchPeerId(addr, peerID string) (string, error) {
	// get the peer id
	resp, err := http.Get(fmt.Sprintf("http://%s:6040/p2pid", addr))
	if err != nil {
		msg := err.Error()
		if resp != nil {
			msg = fmt.Sprintf("%s: %s", resp.Status, err.Error())
		}
		return "", fmt.Errorf(msg)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 300 {
		if resp.StatusCode == 503 {
			// retry
			return fetchPeerId(addr, peerID)
		}
		return "", fmt.Errorf("bifrost p2pid bad request: %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if string(b) != peerID {
		return "", fmt.Errorf("peer id mismatch: %s != %s", string(b), peerID)
	}
	return peerID, nil
}

func fetchBlockHeight(addr string) (int64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/status", addr))
	if err != nil {
		return 0, err
	}

	var result struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}

	if resp.StatusCode > 300 {
		if resp.StatusCode == 503 {
			// retry
			return fetchBlockHeight(addr)
		}
		return 0, fmt.Errorf("thornode status bad request: %d", resp.StatusCode)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if err = json.Unmarshal(buf, &result); err != nil {
		return 0, fmt.Errorf("failed to unmarshal tendermint status: %w", err)
	}

	return strconv.ParseInt(result.Result.SyncInfo.LatestBlockHeight, 10, 64)
}

func checkPortOpen(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
