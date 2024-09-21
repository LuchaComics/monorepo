package cmd

// DEVELOPERS NOTE:
// Special thanks to the following link:
// https://github.com/libp2p/go-libp2p/blob/master/examples/chat-with-rendezvous/flags.go

import (
	"strings"

	maddr "github.com/multiformats/go-multiaddr"
)

// A new type we need for writing a custom flag parser
type addrList []maddr.Multiaddr

//	func (al *addrList) String() string {
//		strs := make([]string, len(*al))
//		for i, addr := range *al {
//			strs[i] = addr.String()
//		}
//		return strings.Join(strs, ",")
//	}
//
//	func (al *addrList) Set(value string) error {
//		addr, err := maddr.NewMultiaddr(value)
//		if err != nil {
//			return err
//		}
//		*al = append(*al, addr)
//		return nil
//	}

func StringToAddres(addrString string) (maddrs []maddr.Multiaddr, err error) {
	// Defensive code: If no string specified then simply return empty values.
	if addrString == "" {
		return
	}
	addrStrings := strings.Split(addrString, ",")
	return StringsToAddrs(addrStrings)
}

func StringsToAddrs(addrStrings []string) (maddrs []maddr.Multiaddr, err error) {
	for _, addrString := range addrStrings {
		addr, err := maddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return
}

var (
	flagRendezvousString string
	flagBootstrapPeers   string
	flagListenAddresses  string
	flagProtocolID       string

	flagListenPort int
)
