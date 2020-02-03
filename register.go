package binary

import (
	"hash/adler32"
	"hash/fnv"
	"reflect"
	"sync"
)

type customTypes struct {
	*sync.Map
}

func (c customTypes) storeType(p reflect.Type) {
	typeHash := c.getHash(p)
	c.Store(typeHash, p)
}

func (c customTypes) loadType(typeHash uint32) (reflect.Type, bool) {
	p, ok := c.Load(typeHash)
	if !ok {
		return nil, false
	}
	return p.(reflect.Type), ok
}

func (c customTypes) isStoredType(p reflect.Type) bool {
	typeHash := c.getHash(p)

	_, ok := c.Load(typeHash)
	return ok
}

func (c customTypes) getHash(p reflect.Type) uint32 {
	s := p.Name()
	if p.Kind() == reflect.Ptr {
		s = p.Elem().Name()
	}
	fnvHash := fnv.New32()
	fnvHash.Write([]byte(s))

	adlerHash := adler32.New()
	adlerHash.Write([]byte(s))

	return fnvHash.Sum32() ^ adlerHash.Sum32()
}

var knownCustomTypes = customTypes{new(sync.Map)}

func RegisterType(rType reflect.Type) (err error) {
	knownCustomTypes.storeType(rType)
	_, err = scan(rType)
	return
}
