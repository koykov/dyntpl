package dyntpl

import (
	"sync"

	"github.com/koykov/fastconv"
)

// Template database.
// Contains two indexes describes two types of pairs between templates and keys/IDs.
type db struct {
	mux sync.RWMutex
	// ID index. Value is an offset in the tpl array.
	idxID map[int]int
	// Key index. Value is an offset in the tpl array as well.
	idxKey map[string]int
	// Hash index. Value is an offset in the tpl array.
	idxHash map[uint64]int
	// Templates storage.
	tpl []*Tpl
}

func initDB() *db {
	db := &db{
		idxID:   make(map[int]int),
		idxKey:  make(map[string]int),
		idxHash: make(map[uint64]int),
	}
	return db
}

// Save template tree in the storage and make two pairs (ID-tpl and key-tpl).
func (db *db) set(id int, key string, tree *Tree) {
	tpl := Tpl{
		ID:   id,
		Key:  key,
		tree: tree,
	}
	db.mux.Lock()
	var idx int
	if idx = db.getIdxLF(id, key); idx >= 0 && idx < len(db.tpl) {
		db.tpl[idx] = &tpl
	} else {
		db.tpl = append(db.tpl, &tpl)
		idx = len(db.tpl) - 1
		if id >= 0 {
			db.idxID[id] = idx
		}
		if key != "-1" {
			db.idxKey[key] = idx
		}
	}
	if _, ok := db.idxHash[tree.hsum]; !ok {
		db.idxHash[tree.hsum] = idx
	}
	db.mux.Unlock()
}

// Get first template found by key or ID.
func (db *db) get(id int, key string) (tpl *Tpl) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	if idx := db.getIdxLF(id, key); idx >= 0 && idx < len(db.tpl) {
		tpl = db.tpl[idx]
	}
	return
}

// Lock-free index getter.
//
// Returns first available index by key or ID.
func (db *db) getIdxLF(id int, key string) (idx int) {
	idx = -1
	if idx1, ok := db.idxKey[key]; ok && idx1 != -1 {
		idx = idx1
	} else if idx1, ok := db.idxID[id]; ok && idx1 != -1 {
		idx = idx1
	}
	return
}

// Get template by ID.
func (db *db) getID(id int) *Tpl {
	return db.get(id, "-1")
}

// Get template by key.
func (db *db) getKey(key string) *Tpl {
	return db.get(-1, key)
}

// Get template by key and fallback key.
func (db *db) getKey1(key, key1 string) (tpl *Tpl) {
	idx := -1
	db.mux.RLock()
	defer db.mux.RUnlock()
	idx1, ok := db.idxKey[key]
	if !ok {
		idx1, ok = db.idxKey[key1]
	}
	if ok {
		idx = idx1
	}
	if idx >= 0 && idx < len(db.tpl) {
		tpl = db.tpl[idx]
	}
	return
}

// Get parsed tree by hash sum.
func (db *db) getTreeByHash(hsum uint64) *Tree {
	db.mux.RLock()
	defer db.mux.RUnlock()
	if idx, ok := db.idxHash[hsum]; ok && idx >= 0 && idx < len(db.tpl) {
		return db.tpl[idx].tree
	}
	return nil
}

// Get template by list of keys describes as bytes arrays.
func (db *db) getBKeys(bkeys [][]byte) (tpl *Tpl) {
	l := len(bkeys)
	if l == 0 {
		return
	}
	db.mux.RLock()
	defer db.mux.RUnlock()
	_ = bkeys[l-1]
	for i := 0; i < l; i++ {
		if idx, ok := db.idxKey[fastconv.B2S(bkeys[i])]; ok && idx >= 0 && idx < len(db.tpl) {
			tpl = db.tpl[idx]
			return
		}
	}
	return
}
