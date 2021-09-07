package dyntpl

import (
	"sync"

	"github.com/koykov/fastconv"
)

type db struct {
	mux    sync.RWMutex
	idxID  map[int]int
	idxKey map[string]int
	tpl    []*Tpl
}

func initDB() *db {
	db := &db{
		idxID:  make(map[int]int),
		idxKey: make(map[string]int),
	}
	return db
}

func (db *db) set(id int, key string, tree *Tree) {
	tpl := Tpl{
		Id:   id,
		Key:  key,
		tree: tree,
	}
	db.mux.Lock()
	db.tpl = append(db.tpl, &tpl)
	db.idxID[id] = len(db.tpl) - 1
	db.idxKey[key] = len(db.tpl) - 1
	db.mux.Unlock()
}

func (db *db) get(id int, key string) (tpl *Tpl) {
	idx := -1
	db.mux.RLock()
	defer db.mux.RUnlock()
	if idx1, ok := db.idxKey[key]; ok && idx1 != -1 {
		idx = idx1
	} else if idx1, ok := db.idxID[id]; ok && idx1 != -1 {
		idx = idx1
	}
	if idx >= 0 && idx < len(db.tpl) {
		tpl = db.tpl[idx]
	}
	return
}

func (db *db) getID(id int) *Tpl {
	return db.get(id, "-2")
}

func (db *db) getKey(key string) *Tpl {
	return db.get(-2, key)
}

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
