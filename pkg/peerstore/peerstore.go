package peerstore

import (
	"sync"
)

// Peer represents a connected peer
type Peer struct {
	ID   string
	Addr string
}

// PeerStore manages connected peers
type PeerStore struct {
	mu    sync.RWMutex
	peers map[string]*Peer // key: peer ID
	addrs map[string]string // key: addr -> peer ID (for reverse lookup)
}

// NewPeerStore creates a new peer store
func NewPeerStore() *PeerStore {
	return &PeerStore{
		peers: make(map[string]*Peer),
		addrs: make(map[string]string),
	}
}

// Add adds or updates a peer in the store
func (ps *PeerStore) Add(id, addr string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Remove old address mapping if peer exists
	if existing, ok := ps.peers[id]; ok {
		delete(ps.addrs, existing.Addr)
	}

	peer := &Peer{
		ID:   id,
		Addr: addr,
	}
	ps.peers[id] = peer
	ps.addrs[addr] = id
}

// Get retrieves a peer by ID
func (ps *PeerStore) Get(id string) (*Peer, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	peer, ok := ps.peers[id]
	return peer, ok
}

// GetByAddr retrieves a peer by address
func (ps *PeerStore) GetByAddr(addr string) (*Peer, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	id, ok := ps.addrs[addr]
	if !ok {
		return nil, false
	}
	return ps.peers[id], true
}

// Remove removes a peer from the store
func (ps *PeerStore) Remove(id string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if peer, ok := ps.peers[id]; ok {
		delete(ps.addrs, peer.Addr)
		delete(ps.peers, id)
	}
}

// GetAll returns all connected peers
func (ps *PeerStore) GetAll() []*Peer {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	peers := make([]*Peer, 0, len(ps.peers))
	for _, peer := range ps.peers {
		peers = append(peers, peer)
	}
	return peers
}

// Count returns the number of connected peers
func (ps *PeerStore) Count() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return len(ps.peers)
}

// Exists checks if a peer exists by ID
func (ps *PeerStore) Exists(id string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, ok := ps.peers[id]
	return ok
}
