// Package enumcache provides simple methods to store and retrieve enumerations with a numerical ID and string name.
package enumcache

import (
	"fmt"
	"sort"
	"sync"

	"github.com/nexmoinc/gosrvlib/pkg/enumbitmap"
)

// EnumCache handles name and id value mapping.
type EnumCache struct {
	sync.RWMutex
	id   map[string]int
	name map[int]string
}

// New returns a new empty EnumCache.
func New() *EnumCache {
	return &EnumCache{
		id:   make(map[string]int),
		name: make(map[int]string),
	}
}

// Set a single id-name key-value.
func (ec *EnumCache) Set(id int, name string) {
	ec.Lock()
	defer ec.Unlock()

	ec.name[id] = name
	ec.id[name] = id
}

// ID returns the numerical ID associated to the given name.
func (ec *EnumCache) ID(name string) (int, error) {
	ec.RLock()
	defer ec.RUnlock()

	id, ok := ec.id[name]
	if !ok {
		return 0, fmt.Errorf("cache name not found: %s", name)
	}

	return id, nil
}

// Name returns the name associated with the given numerical ID.
func (ec *EnumCache) Name(id int) (string, error) {
	ec.RLock()
	defer ec.RUnlock()

	name, ok := ec.name[id]
	if !ok {
		return "", fmt.Errorf("cache ID not found: %d", id)
	}

	return name, nil
}

// SortNames returns a list of sorted names.
func (ec *EnumCache) SortNames() []string {
	sorted := make([]string, 0, len(ec.id))
	for name := range ec.id {
		sorted = append(sorted, name)
	}

	sort.Strings(sorted)

	return sorted
}

// SortIDs returns a list of sorted IDs.
func (ec *EnumCache) SortIDs() []int {
	sorted := make([]int, 0, len(ec.name))
	for id := range ec.name {
		sorted = append(sorted, id)
	}

	sort.Ints(sorted)

	return sorted
}

// DecodeBinaryMapUint64 decodes a uint64 binary map into a list of string names.
// The EnumCache must contain the mapping between the bit values and the names.
func (ec *EnumCache) DecodeBinaryMapUint64(v uint64) (s []string, err error) {
	return enumbitmap.MapUint64ToStrings(ec.name, v) // nolint:wrapcheck
}

// DecodeBinaryMapUint32 decodes a uint32 binary map into a list of string names.
// The EnumCache must contain the mapping between the bit values and the names.
func (ec *EnumCache) DecodeBinaryMapUint32(v uint32) (s []string, err error) {
	return enumbitmap.MapUint32ToStrings(ec.name, v) // nolint:wrapcheck
}

// DecodeBinaryMapUint16 decodes a uint16 binary map into a list of string names.
// The EnumCache must contain the mapping between the bit values and the names.
func (ec *EnumCache) DecodeBinaryMapUint16(v uint16) (s []string, err error) {
	return enumbitmap.MapUint16ToStrings(ec.name, v) // nolint:wrapcheck
}

// DecodeBinaryMapUint8 decodes a uint8 binary map into a list of string names.
// The EnumCache must contain the mapping between the bit values and the names.
func (ec *EnumCache) DecodeBinaryMapUint8(v uint8) (s []string, err error) {
	return enumbitmap.MapUint8ToStrings(ec.name, v) // nolint:wrapcheck
}

// EncodeBinaryMapUint64 encode a list of string names into a uint64 binary map.
// The EnumCache must contain the mapping between the bit values and the names.
func (ec *EnumCache) EncodeBinaryMapUint64(s []string) (v uint64, err error) {
	return enumbitmap.MapStringsToUint64(ec.id, s) // nolint:wrapcheck
}

// EncodeBinaryMapUint32 encode a list of string names into a uint32 binary map.
// The EnumCache must contain the mapping between the bit values and the names.
func (ec *EnumCache) EncodeBinaryMapUint32(s []string) (v uint32, err error) {
	return enumbitmap.MapStringsToUint32(ec.id, s) // nolint:wrapcheck
}

// EncodeBinaryMapUint16 encode a list of string names into a uint16 binary map.
// The EnumCache must contain the mapping between the bit values and the names.
func (ec *EnumCache) EncodeBinaryMapUint16(s []string) (v uint16, err error) {
	return enumbitmap.MapStringsToUint16(ec.id, s) // nolint:wrapcheck
}

// EncodeBinaryMapUint8 encode a list of string names into a uint8 binary map.
// The EnumCache must contain the mapping between the bit values and the names.
func (ec *EnumCache) EncodeBinaryMapUint8(s []string) (v uint8, err error) {
	return enumbitmap.MapStringsToUint8(ec.id, s) // nolint:wrapcheck
}
