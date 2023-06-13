package main

import (
	"encoding/json"
	"os"
	"sync"
)

type Database struct {
	mtx  sync.Mutex
	file string
}

func NewDatabase(file string) *Database {
	return &Database{
		file: file,
	}
}

func (t *Database) GetCount(id string) (int, bool, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	m, err := t.getCounts()
	if err != nil {
		return 0, false, err
	}

	v, ok := m[id]
	return v, ok, nil
}

func (t *Database) SetCount(id string, v int) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	m, err := t.getCounts()
	if err != nil {
		return err
	}

	m[id] = v

	return t.writeCounts(m)
}

func (t *Database) getCounts() (map[string]int, error) {
	f, err := t.openFile(false)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[string]int)
	err = json.NewDecoder(f).Decode(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (t *Database) writeCounts(m map[string]int) error {
	f, err := t.openFile(true)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(m)
}

func (t *Database) openFile(write bool) (f *os.File, err error) {

	if write {
		f, err = os.Create(t.file)
	} else {
		f, err = os.Open(t.file)
		if os.IsNotExist(err) {
			f, err = os.Create(t.file)
			if err != nil {
				return nil, err
			}

			err = json.NewEncoder(f).
				Encode(make(map[string]int))

			// This is an intended bug ti demonstrate the case where the
			// dabase driver fails.
			// The initial write of the empty map should actually be
			// followed by f.Seek(0, 0) so that the next read can
			// receive the written content.
		}
	}

	return f, err
}
