package dyntpl

import (
	"sync"

	"github.com/koykov/fastconv"
)

// Template database.
// Contains two indexes describes two types of pairs between templates and keys/IDs.
type db struct {
	mux sync.RWMutex
	// ID index. Value is a offset in the tpl array.
	idxID map[int]int
	// Key index. Value is a offset in the tpl array as well.
	idxKey map[string]int
	// Templates storage.
	tpl []*Tpl
}

func initDB() *db {
	db := &db{
		idxID:  make(map[int]int),
		idxKey: make(map[string]int),
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
	if idx := db.getIdxLF(id, key); idx >= 0 && idx < len(db.tpl) {
		db.tpl[idx] = &tpl
	} else {
		db.tpl = append(db.tpl, &tpl)
		if id >= 0 {
			db.idxID[id] = len(db.tpl) - 1
		}
		if key != "-1" {
			db.idxKey[key] = len(db.tpl) - 1
		}
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
