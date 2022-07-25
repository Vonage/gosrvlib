package filter

import "reflect"

// fieldCache caches reflectPath by type and field.
type fieldCache struct {
	cache map[string]map[string]reflectPath
}

// Get return the reflectPath from the cache for a field given its type and path, and true if it's found.
// Returns (nil, false) if not found.
func (c *fieldCache) Get(t reflect.Type, fieldPath string) (reflectPath, bool) {
	fields := c.getFieldsMap(t)
	path, ok := fields[fieldPath]

	return path, ok
}

// Set stores a reflectPath in the cache by its type and path.
func (c *fieldCache) Set(t reflect.Type, fieldPath string, path reflectPath) {
	fields := c.getFieldsMap(t)
	fields[fieldPath] = path
}

func (c *fieldCache) getFieldsMap(t reflect.Type) map[string]reflectPath {
	if c.cache == nil {
		c.cache = map[string]map[string]reflectPath{}
	}

	tKey := t.String()

	fields, ok := c.cache[tKey]
	if !ok {
		fields = map[string]reflectPath{}
		c.cache[tKey] = fields
	}

	return fields
}
