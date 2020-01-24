package memory

import (
	"git.maxset.io/maxsetdev/srverror"
	"git.maxset.io/web/knaxim/internal/database"
)

type Ownerbase struct {
	Database
}

func (ob *Ownerbase) Reserve(id database.OwnerID, name string) (database.OwnerID, error) {
	ob.lock.Lock()
	defer ob.lock.Unlock()
	if id.Type == 'u' {
		if _, ok := ob.Owners.UserName[name]; ok {
			return id, database.ErrNameTaken
		}
	} else if id.Type == 'g' {
		if _, ok := ob.Owners.GroupName[name]; ok {
			return id, database.ErrNameTaken
		}
	} else {
		return id, srverror.Basic(500, "Server Error", "unrecognized id type")
	}

	for idstr := id.String(); true; idstr = id.String() {
		if _, ok := ob.Owners.ID[idstr]; !ok {
			break
		}
		id = id.Mutate()
	}
	ob.Owners.ID[id.String()] = nil
	switch id.Type {
	case 'u':
		ob.Owners.UserName[name] = nil
	case 'g':
		ob.Owners.GroupName[name] = nil
	}
	return id, nil
}

func (ob *Ownerbase) Insert(u database.Owner) error {
	ob.lock.Lock()
	defer ob.lock.Unlock()
	idstr := u.GetID().String()
	if expectnil, ok := ob.Owners.ID[idstr]; !ok {
		return database.ErrIDNotReserved
	} else if expectnil != nil {
		return database.ErrNameTaken
	}
	switch v := u.(type) {
	case database.UserI:
		expectnil, ok := ob.Owners.UserName[v.GetName()]
		if !ok {
			return database.ErrIDNotReserved
		}
		if expectnil != nil {
			return database.ErrNameTaken
		}
		ob.Owners.UserName[v.GetName()] = v
	case database.GroupI:
		expectnil, ok := ob.Owners.GroupName[v.GetName()]
		if !ok {
			return database.ErrIDNotReserved
		}
		if expectnil != nil {
			return database.ErrNameTaken
		}
		ob.Owners.GroupName[v.GetName()] = v
	default:
		return srverror.Basic(500, "Server Error", "Unrecognized Owner Type")
	}
	ob.Owners.ID[idstr] = u
	return nil
}

func (ob *Ownerbase) Get(id database.OwnerID) (database.Owner, error) {
	ob.lock.RLock()
	defer ob.lock.RUnlock()
	if ob.Owners.ID[id.String()] == nil {
		return nil, database.ErrNotFound
	}
	return ob.Owners.ID[id.String()].Copy(), nil
}

func (ob *Ownerbase) FindUserName(name string) (database.UserI, error) {
	ob.lock.RLock()
	defer ob.lock.RUnlock()
	if ob.Owners.UserName[name] == nil {
		return nil, database.ErrNotFound
	}
	return ob.Owners.UserName[name].Copy().(database.UserI), nil
}

func (ob *Ownerbase) FindGroupName(name string) (database.GroupI, error) {
	ob.lock.RLock()
	defer ob.lock.RUnlock()
	if ob.Owners.GroupName[name] == nil {
		return nil, database.ErrNotFound
	}
	return ob.Owners.GroupName[name].Copy().(database.GroupI), nil
}

func (ob *Ownerbase) GetGroups(id database.OwnerID) (owned []database.GroupI, member []database.GroupI, err error) {
	ob.lock.RLock()
	defer ob.lock.RUnlock()
LOOP:
	for _, grp := range ob.Owners.GroupName {
		if grp.GetOwner().GetID().Equal(id) {
			owned = append(owned, grp.Copy().(database.GroupI))
			continue
		}
		for _, mem := range grp.GetMembers() {
			if mem.GetID().Equal(id) {
				member = append(member, grp.Copy().(database.GroupI))
				continue LOOP
			}
		}
	}
	return
}

func (ob *Ownerbase) Update(o database.Owner) error {
	ob.lock.Lock()
	defer ob.lock.Unlock()
	if ob.Owners.ID[o.GetID().String()] == nil {
		return database.ErrNotFound
	}
	switch v := o.(type) {
	case database.UserI:
		ob.Owners.UserName[v.GetName()] = v
	case database.GroupI:
		ob.Owners.GroupName[v.GetName()] = v
	default:
		return srverror.Basic(500, "Server Error", "Unrecognized owner type")
	}
	ob.Owners.ID[o.GetID().String()] = o
	return nil
}

func (ob *Ownerbase) GetSpace(o database.OwnerID) (int64, error) {
	ob.lock.RLock()
	defer ob.lock.RUnlock()
	if ob.Owners.ID[o.String()] == nil {
		return 0, database.ErrNotFound
	}
	var total int64
	for _, file := range ob.Files {
		if file.GetOwner().GetID().Equal(o) {
			total += ob.Stores[file.GetID().StoreID.String()].FileSize
		}
	}
	return total, nil
}

func (ob *Ownerbase) GetTotalSpace(o database.OwnerID) (int64, error) {
	ob.lock.RLock()
	defer ob.lock.RUnlock()
	if ob.Owners.ID[o.String()] == nil {
		return 0, database.ErrNotFound
	}
	if o.Type == 'u' {
		return 50 << 20, nil
	}
	return 0, nil
}
