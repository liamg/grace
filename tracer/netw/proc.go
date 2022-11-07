package netw

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Connection struct {
	Protocol      string
	LocalAddress  net.IP
	LocalPort     int
	RemoteAddress net.IP
	RemotePort    int
	INode         int
	State         int // TCP_ESTABLISHED etc. see https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/net/tcp_states.h
}

const netPath = "/proc/net/"

func ListTCPConnections() ([]Connection, error) {
	return parseBaseProtocol("tcp")
}

func ListUDPConnections() ([]Connection, error) {
	return parseBaseProtocol("udp")
}

func ListICMPConnections() ([]Connection, error) {
	return parseBaseProtocol("icmp")
}

func parseBaseProtocol(protocol string) ([]Connection, error) {
	v4, err := parseFile(protocol)
	if err != nil {
		return nil, err
	}

	v6, err := parseFile(protocol + "6")
	if err != nil {
		return nil, err
	}

	return append(v4, v6...), nil
}

func parseFile(protocol string) ([]Connection, error) {

	file := netPath + protocol
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var connections []Connection
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "sl") { // skip headings
			continue
		}
		connection, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		connection.Protocol = protocol
		connections = append(connections, *connection)
	}
	return connections, nil
}

func parseLine(line string) (*Connection, error) {

	line = strings.TrimSpace(line)
	fields := strings.Fields(line)

	if len(fields) < 10 {
		return nil, fmt.Errorf("invalid tcp connection: %s", line)
	}

	inode, err := strconv.Atoi(fields[9])
	if err != nil {
		return nil, fmt.Errorf("invalid inode '%s': %w", fields[9], err)
	}

	localip, localport, err := parseIPAndPortFromHex(fields[1])
	if err != nil {
		return nil, fmt.Errorf("invalid local ip '%s': %w", fields[1], err)
	}

	remoteip, remoteport, err := parseIPAndPortFromHex(fields[2])
	if err != nil {
		return nil, fmt.Errorf("invalid remote ip '%s': %w", fields[2], err)
	}

	return &Connection{
		INode:         inode,
		LocalAddress:  localip,
		LocalPort:     localport,
		RemoteAddress: remoteip,
		RemotePort:    remoteport,
		State:         int(hexToByte(fields[3])),
	}, nil
}

func parseIPAndPortFromHex(hex string) (net.IP, int, error) {

	rawip, rawport, found := strings.Cut(hex, ":")
	if !found {
		return nil, 0, fmt.Errorf("invalid hex '%s'", hex)
	}

	var ip net.IP

	switch len(rawip) {
	case 8:
		ip = net.IPv4(
			hexToByte(rawip[6:8]),
			hexToByte(rawip[4:6]),
			hexToByte(rawip[2:4]),
			hexToByte(rawip[0:2]),
		)
	case 32:
		ip = []byte{
			hexToByte(rawip[30:32]),
			hexToByte(rawip[28:30]),
			hexToByte(rawip[26:28]),
			hexToByte(rawip[24:26]),
			hexToByte(rawip[22:24]),
			hexToByte(rawip[20:22]),
			hexToByte(rawip[18:20]),
			hexToByte(rawip[16:18]),
			hexToByte(rawip[14:16]),
			hexToByte(rawip[12:14]),
			hexToByte(rawip[10:12]),
			hexToByte(rawip[8:10]),
			hexToByte(rawip[6:8]),
			hexToByte(rawip[4:6]),
			hexToByte(rawip[2:4]),
			hexToByte(rawip[0:2]),
		}
	default:
		return nil, 0, fmt.Errorf("invalid ipv4 hex '%s'", hex)
	}

	port, err := strconv.ParseInt(rawport, 16, 32)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid port '%s': %w", rawport, err)
	}

	return ip, int(port), nil
}

func hexToByte(hex string) byte {
	b, err := strconv.ParseUint(hex, 16, 8)
	if err != nil {
		return 0
	}
	return byte(b)
}
