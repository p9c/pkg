package kx

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/p9c/pkg/app/slog"
)

type Key struct {
	x, y  *big.Int
	group *Group
}

func GetGroup1() (out *big.Int) {
	out, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A63A3620FFFFFFFFFFFFFFFF", 16)
	return
}
func GetGroup2() (out *big.Int) {
	out, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFF", 16)
	return
}
func GetGroup14() (out *big.Int) {
	out, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF", 16)
	return
}

// copyLeftPad copies the source to the end of the destination
func copyLeftPad(dest, src []byte) {
	padLen := len(dest) - len(src)
	for i := 0; i < padLen; i++ {
		dest[i] = 0
	}
	copy(dest[padLen:], src)
}

func (k *Key) Bytes() (out []byte) {
	if k.y == nil {
		return
	}
	if k.group != nil {
		bitLen := (k.group.p.BitLen() + 7) / 8
		out = make([]byte, bitLen)
		copyLeftPad(out, k.y.Bytes())
	}
	return
}

func (k *Key) String() (out string) {
	if k.y == nil {
		out = ""
	} else {
		out = k.y.String()
	}
	return
}

func (k *Key) IsPrivKey() bool {
	return k.x != nil
}

func NewPubKey(b []byte) (out *Key) {
	out = &Key{y: new(big.Int).SetBytes(b)}
	return
}

type Group struct {
	p, g *big.Int
}

func (grp *Group) P() (p *big.Int) {
	p = &big.Int{}
	p.Set(grp.p)
	return
}

func (grp *Group) G() (g *big.Int) {
	g = &big.Int{}
	g.Set(grp.g)
	return
}

func (grp *Group) GenPrivKey() (key *Key, err error) {
	s := rand.Reader
	var x *big.Int
	if x, err = rand.Int(s, grp.p); slog.Check(err) {
		return
	}
	zero := big.NewInt(0)
	for x.Cmp(zero) == 0 {
		if x, err = rand.Int(s, grp.p); slog.Check(err) {
			return
		}
	}
	key = &Key{
		x:     x,
		y:     new(big.Int).Exp(grp.g, x, grp.p),
		group: grp,
	}
	return
}

// GetGroup returns a Diffie Hellman group by its ID as defined in RFC2409 and 3526
// an id of 0 will select the recommended group 14
func GetGroup(gID int) (group *Group, err error) {
	if gID <= 0 {
		gID = 14
	}
	switch gID {
	case 1:
		group = &Group{
			g: new(big.Int).SetInt64(2),
			p: GetGroup1(),
		}
	case 2:
		group = &Group{
			g: new(big.Int).SetInt64(2),
			p: GetGroup2(),
		}
	case 14:
		group = &Group{
			g: new(big.Int).SetInt64(2),
			p: GetGroup14(),
		}
	default:
		group = nil
		err = errors.New("Unknown group")
	}
	return
}

func (grp *Group) ComputeKey(pubkey *Key, privkey *Key) (key *Key, err error) {
	if grp.p == nil {
		err = errors.New("invalid group")
		return
	}
	if pubkey.y == nil {
		err = errors.New("invalid public key")
		return
	}
	if pubkey.y.Sign() <= 0 || pubkey.y.Cmp(grp.p) >= 0 {
		err = errors.New("Diffie Hellman parameter out of bounds")
		return
	}
	if privkey.x == nil {
		err = errors.New("invalid private key")
		return
	}
	k := new(big.Int).Exp(pubkey.y, privkey.x, grp.p)
	key = &Key{y: k, group: grp}
	return
}
