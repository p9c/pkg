package kx

import (
	"fmt"
	"math/big"
	"testing"
)

type peer struct {
	priv  *Key
	group *Group
	pub   *Key
}

func newPeer(g *Group) *peer {
	ret := new(peer)
	ret.priv, _ = g.GenPrivKey()
	ret.group = g
	return ret
}

func (self *peer) getPubKey() []byte {
	return self.priv.Bytes()
}

func (self *peer) recvPeerPubKey(pub []byte) {
	pubKey := NewPubKey(pub)
	self.pub = pubKey
}

func (self *peer) getKey() []byte {
	k, err := self.group.ComputeKey(self.pub, self.priv)
	if err != nil {
		return nil
	}
	return k.Bytes()
}

func exchangeKey(p1, p2 *peer) error {
	pub1 := p1.getPubKey()
	pub2 := p2.getPubKey()

	p1.recvPeerPubKey(pub2)
	p2.recvPeerPubKey(pub1)

	key1 := p1.getKey()
	key2 := p2.getKey()

	if key1 == nil {
		return fmt.Errorf("p1 has nil key")
	}
	if key2 == nil {
		return fmt.Errorf("p2 has nil key")
	}

	for i, k := range key1 {
		if key2[i] != k {
			return fmt.Errorf("%vth byte does not same", i)
		}
	}
	return nil
}

func TestKeyExchange(t *testing.T) {
	group, _ := GetGroup(14)
	p1 := newPeer(group)
	p2 := newPeer(group)

	err := exchangeKey(p1, p2)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestPIsNotMutable(t *testing.T) {
	d, _ := GetGroup(0)
	p := d.p.String()
	d.P().Set(big.NewInt(1))
	if p != d.p.String() {
		t.Errorf("group's prime mutated externally, should be %s, was changed to %s", p, d.p.String())
	}
}

func TestGIsNotMutable(t *testing.T) {
	d, _ := GetGroup(0)
	g := d.g.String()
	d.G().Set(big.NewInt(0))
	if g != d.g.String() {
		t.Errorf("group's generator mutated externally, should be %s, was changed to %s", g, d.g.String())
	}
}
