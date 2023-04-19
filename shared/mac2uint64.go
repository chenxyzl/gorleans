package shared

import (
	"net"
)

// Mac2Uint64 mac addr to uint64 "00:11:22:33:44:55" Mac2Uint64 big:75384242876037 lit:596532804581120
func Mac2Uint64(mac1 string, bigOrLit bool) (uint64, error) {
	mac, err := net.ParseMAC(mac1)
	if err != nil {
		return 0, err
	}

	if bigOrLit {
		return uint64(mac[0])<<40 |
			uint64(mac[1])<<32 |
			uint64(mac[2])<<24 |
			uint64(mac[3])<<16 |
			uint64(mac[4])<<8 |
			uint64(mac[5]), nil
	} else {
		return uint64(mac[0]) |
			uint64(mac[1])<<8 |
			uint64(mac[2])<<16 |
			uint64(mac[3])<<24 |
			uint64(mac[4])<<32 |
			uint64(mac[5])<<40, nil
	}
}
