package dht

import (
	"github.com/anacrolix/dht/v2/krpc"
	"math/rand"
)

// todo version number
const peerIDPrefix = "-BM0001-"

func RandomPeerID() krpc.ID {
	/* > The peer_id is exactly 20 bytes (characters) long.
	 * >
	 * > There are mainly two conventions how to encode client and client version information into the peer_id,
	 * > Azureus-style and Shadow's-style.
	 * >
	 * > Azureus-style uses the following encoding: '-', two characters for client id, four ascii digits for version
	 * > number, '-', followed by random numbers.
	 * >
	 * > For example: '-AZ2060-'...
	 *
	 * https://wiki.theory.org/BitTorrentSpecification#peer_id
	 *
	 * We encode the version number as:
	 *  - First two digits for the major version number
	 *  - Last two digits for the minor version number
	 *  - Patch version number is not encoded.
	 */
	var peerID [20]byte

	i := 0
	for _, c := range peerIDPrefix {
		peerID[i] = byte(c)
		i++
	}
	for {
		if i >= 20 {
			break
		}
		peerID[i] = randomDigit()
		i++
	}

	return peerID
}

func randomDigit() byte {
	var max, min int
	max, min = '9', '0'
	return byte(rand.Intn(max-min) + min)
}
