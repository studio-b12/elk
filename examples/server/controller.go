package main

import (
	"sync"

	"github.com/studio-b12/whoops"
)

type Controller struct {
	rwx sync.RWMutex
	db  *Database
}

func NewController(db *Database) *Controller {
	return &Controller{
		db: db,
	}
}

func (t *Controller) GetCount(id string) (Count, error) {
	t.rwx.RLock()
	defer t.rwx.RUnlock()

	count, ok, err := t.db.GetCount(id)
	if err != nil {
		return Count{}, whoops.Wrap(ErrorInternal, err, "failed getting count from database")
	}

	if !ok {
		return Count{}, whoops.Detailed(ErrorCountNotFound)
	}

	c := Count{
		Id:    id,
		Count: count,
	}
	return c, nil
}

func (t *Controller) IncrementCount(id string) (Count, error) {
	t.rwx.Lock()
	t.rwx.Unlock()

	count, _, err := t.db.GetCount(id)
	if err != nil {
		return Count{}, whoops.Wrap(ErrorInternal, err, "failed getting count from database")
	}

	count++

	err = t.db.SetCount(id, count)
	if err != nil {
		return Count{}, whoops.Wrap(ErrorInternal, err, "failed setting count to database")
	}

	c := Count{
		Id:    id,
		Count: count,
	}
	return c, nil
}
