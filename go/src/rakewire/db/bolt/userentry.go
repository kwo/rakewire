package bolt

import (
	"github.com/boltdb/bolt"
	m "rakewire/model"
)

// UserEntrySave saves userentries to the database.
func (z *Service) UserEntrySave(userentries []*m.UserEntry) error {

	z.Lock()
	defer z.Unlock()
	err := z.db.Update(func(tx *bolt.Tx) error {
		for _, userentry := range userentries {
			if err := kvSave(m.UserEntryEntity, userentry, tx); err != nil {
				return err
			}
		}
		return nil
	})

	return err

}
