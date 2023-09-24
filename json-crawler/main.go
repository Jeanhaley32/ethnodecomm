package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	jsonpath, writefile string
)

func init() {
	flag.StringVar(&jsonpath, "jsonpath", "node-list.json", "Path to json file containing bootnodes. Defauklt is ./node-list.json")
	flag.StringVar(&writefile, "writefile", "finalized-node-list.json", "Path to write json file containing nodes with added neighbors. Default is ./finalized-node-list.json")
	flag.Parse()
}

// a list of nodeleaf structs.
type nodetree map[string]nodeleaf

// json struct that contains an enode struct and a list of that enodes neighbors.
type nodeleaf struct {
	Enode     enode.Node    `json:"enode"`
	Neighbors []*enode.Node `json:"neighbors"`
}

// JSON struct representing an ENR record
type enrJSON struct {
	Seq           uint64 `json:"seq"`
	Record        string `json:"record"`
	Score         int    `json:"score"`
	FirstResponse string `json:"firstResponse"`
	LastResponse  string `json:"lastResponse"`
	LastCheck     string `json:"lastCheck"`
}

func main() {
	// Check if the jsonpath and writefile are empty.
	// If they are empty, exit with a fatal error.
	if jsonpath == "" {
		log.Fatalf("jsonpath is empty")
	}
	if writefile == "" {
		log.Fatalf("writepath is empty")
	}
	// Create a target list of enode nodes to traverse.
	var target []*enode.Node

	// Create a nodetree to store the nodes and their neighbors.
	var nt nodetree

	// Open the JSON file containing the bootnodes.
	file, err := os.Open(jsonpath)
	if err != nil {
		log.Fatalf("Failed to open JSON file: %s", err.Error())
	}

	// Decode the JSON file into a list of ENR records.'
	var entries map[string]enrJSON               // map of string to enrJSON struct for each entry
	err = json.NewDecoder(file).Decode(&entries) // decode the JSON file into the entries map
	if err != nil {
		log.Fatalf("Failed to decode JSON file: %s", err.Error())
	}

	// populate target with the enode nodes from the JSON file.
	for _, entry := range entries {
		// Parse the ENR record from the JSON file.
		node, err := enode.Parse(enode.ValidSchemes) // I'm not sure how to parse the enrJSON struct into an enode.Node struct.
		if err != nil {
			log.Fatalf("Failed to parse ENR record: %s", err.Error())
		}
		// Append the node to the target list.
		target = append(target, node)
	}
	// This is going to be very slow. Ideally, we should use parallelism to speed this up. But for now, this will do.
	// i'm tired and i want to go to bed.
	for _, node := range target {
		disc, _, err := startV4("", node.ID().String(), "", "")
		if err != nil {
			log.Fatalf("Failed to start ephemeral discovery node: %s", err.Error())
		}
		defer disc.Close()
		// Check if the node is nil before accessing its ID
		if node != nil {
			// Lookup the neighbors for the current node.
			neighbors := disc.LookupPubkey(node.Pubkey())
			// Create a new nodeleaf struct with the current node and its neighbors.
			nt[node.ID().String()] = nodeleaf{Enode: *node, Neighbors: neighbors}
		} else {
			fmt.Println("Node is nil")
		}
	}

	// Marshal the nodeleaf list to JSON.
	json, err := json.MarshalIndent(nt, "", "  ")
	// write json to file.
	err = os.WriteFile(writefile, json, 0644)
	if err != nil {
		log.Fatalf("Failed to write json to file: %s", err.Error())
	}
}

// startV4 starts an ephemeral discovery V4 node.
func startV4(nodekey, bootnodes, nodedb, extaddr string) (*discover.UDPv4, discover.Config, error) {
	ln, config := makeDiscoveryConfig(nodekey, bootnodes, nodedb)
	socket := listen(ln, extaddr)
	disc, err := discover.ListenV4(socket, ln, config)
	if err != nil {
		return nil, config, err
	}
	return disc, config, nil
}

// makeDiscoveryConfig creates a discovery configuration.
// A discovery configuration is used to create a discovery node.
func makeDiscoveryConfig(nodekey, bootnodes, nodedb string) (*enode.LocalNode, discover.Config) {
	var cfg discover.Config

	if nodekey != "" {
		key, err := crypto.HexToECDSA(nodekey)
		if err != nil {
			exit(fmt.Errorf("-%s: %v", nodekey, err))
		}
		cfg.PrivateKey = key
	} else {
		cfg.PrivateKey, _ = crypto.GenerateKey()
	}

	if bootnodes != "" {
		bn, err := parseBootnodes(bootnodes)
		if err != nil {
			exit(err)
		}
		cfg.Bootnodes = bn
	}

	dbpath := nodedb
	db, err := enode.OpenDB(dbpath)
	if err != nil {
		exit(err)
	}
	ln := enode.NewLocalNode(db, cfg.PrivateKey)
	return ln, cfg
}

func listen(ln *enode.LocalNode, extAddr string) *net.UDPConn {
	addr := "0.0.0.0:0"
	socket, err := net.ListenPacket("udp4", addr)
	if err != nil {
		exit(err)
	}

	// Configure UDP endpoint in ENR from listener address.
	usocket := socket.(*net.UDPConn)
	uaddr := socket.LocalAddr().(*net.UDPAddr)
	if uaddr.IP.IsUnspecified() {
		ln.SetFallbackIP(net.IP{127, 0, 0, 1})
	} else {
		ln.SetFallbackIP(uaddr.IP)
	}
	ln.SetFallbackUDP(uaddr.Port)

	if extAddr != "" {
		ip, port, ok := parseExtAddr(extAddr)
		if !ok {
			exit(fmt.Errorf("invalid external address %q", extAddr))
		}
		ln.SetStaticIP(ip)
		if port != 0 {
			ln.SetFallbackUDP(port)
		}
	}

	return usocket
}

// exit prints the error to stderr and exits with status 1.
func exit(err interface{}) {
	if err == nil {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// parseExtAddr parses an external address specification.
func parseExtAddr(spec string) (ip net.IP, port int, ok bool) {
	ip = net.ParseIP(spec)
	if ip != nil {
		return ip, 0, true
	}
	host, portstr, err := net.SplitHostPort(spec)
	if err != nil {
		return nil, 0, false
	}
	ip = net.ParseIP(host)
	if ip == nil {
		return nil, 0, false
	}
	port, err = strconv.Atoi(portstr)
	if err != nil {
		return nil, 0, false
	}
	return ip, port, true
}

// parseBootnodes parses a comma-separated list of bootnodes.
func parseBootnodes(bootNodes string) ([]*enode.Node, error) {
	s := params.MainnetBootnodes
	if bootNodes != "" {
		input := bootNodes
		if input == "" {
			return nil, nil
		}
		s = strings.Split(input, ",")
	}
	nodes := make([]*enode.Node, len(s))
	var err error
	for i, record := range s {
		nodes[i], err = parseNode(record)
		if err != nil {
			return nil, fmt.Errorf("invalid bootstrap node: %v", err)
		}
	}
	return nodes, nil
}

// parseNode parses a node record and verifies its signature.
func parseNode(source string) (*enode.Node, error) {
	if strings.HasPrefix(source, "enode://") {
		return enode.ParseV4(source)
	}
	r, err := parseRecord(source)
	if err != nil {
		return nil, err
	}
	return enode.New(enode.ValidSchemes, r)
}

// pulled from enrcmd.go in dsp2p cli library.
// parseRecord parses a node record from hex, base64, or raw binary input.
func parseRecord(source string) (*enr.Record, error) {
	bin := []byte(source)
	if d, ok := decodeRecordHex(bytes.TrimSpace(bin)); ok {
		bin = d
	} else if d, ok := decodeRecordBase64(bytes.TrimSpace(bin)); ok {
		bin = d
	}
	var r enr.Record
	err := rlp.DecodeBytes(bin, &r)
	return &r, err
}

// decodeRecordHex decodes a hex-encoded node record.
func decodeRecordHex(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("0x")) {
		b = b[2:]
	}
	dec := make([]byte, hex.DecodedLen(len(b)))
	_, err := hex.Decode(dec, b)
	return dec, err == nil
}

// decodeRecordBase64 decodes a base64-encoded node record.
func decodeRecordBase64(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("enr:")) {
		b = b[4:]
	}
	dec := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(dec, b)
	return dec[:n], err == nil
}

func (n *nodeleaf) AppendNeighbor(n2 []*enode.Node) {
	n.Neighbors = append(n.Neighbors, n2...)
}
