package filter

import "reflect"

type fieldCache struct {
	cache map[string]map[string]reflectPath
}

func (c *fieldCache) Get(t reflect.Type, fieldPath string) (reflectPath, bool) {
	fields := c.get(t)
	path, ok := fields[fieldPath]

	return path, ok
}

func (c *fieldCache) Set(t reflect.Type, fieldPath string, path reflectPath) {
	fields := c.get(t)
	fields[fieldPath] = path
}

func (c *fieldCache) get(t reflect.Type) map[string]reflectPath {
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
