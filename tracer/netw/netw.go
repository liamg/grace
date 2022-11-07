package netw

func ListConnections() ([]Connection, error) {
	tcpConnections, err := ListTCPConnections()
	if err != nil {
		return nil, err
	}
	udpConnections, err := ListUDPConnections()
	if err != nil {
		return nil, err
	}
	icmpConnections, err := ListICMPConnections()
	if err != nil {
		return nil, err
	}
	return append(append(tcpConnections, udpConnections...), icmpConnections...), nil
}
