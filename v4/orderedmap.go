package orderedmap

import (
	"iter"
	"slices"
)

type Of[K comparable, V any] struct {
	kv map[K]V
	ll *[]K
}

func Make[K comparable, V any]() Of[K, V] {
	return Of[K, V]{
		kv: map[K]V{},
		ll: new([]K),
	}
}

// NewOrderedMapWithCapacity creates a map with enough pre-allocated space to
// hold the specified number of elements.
func MakeWithCapacity[K comparable, V any](capacity int) Of[K, V] {
	ll := make([]K, 0, capacity)
	return Of[K, V]{
		kv: map[K]V{},
		ll: &ll,
	}
}

// Get returns the value for a key. If the key does not exist, the second return
// parameter will be false and the value will be nil.
func (m Of[K, V]) Get(key K) (value V, ok bool) {
	value, ok = m.kv[key]
	return
}

// Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).
func (m Of[K, V]) Set(key K, value V) bool {
	_, alreadyExist := m.kv[key]
	m.kv[key] = value
	if alreadyExist {
		return false
	}

	*m.ll = append(*m.ll, key)
	return true
}

// ReplaceKey replaces an existing key with a new key while preserving order of
// the value. This function will return true if the operation was successful, or
// false if 'originalKey' is not found OR 'newKey' already exists (which would be an overwrite).
func (m Of[K, V]) ReplaceKey(originalKey, newKey K) bool {
	element, originalExists := m.kv[originalKey]
	_, newKeyExists := m.kv[newKey]
	if originalExists && !newKeyExists {
		delete(m.kv, originalKey)
		m.kv[newKey] = element
		lli := slices.IndexFunc(*m.ll, func(key K) bool { return key == originalKey })
		(*m.ll)[lli] = newKey
		return true
	}
	return false
}

// GetOrDefault returns the value for a key. If the key does not exist, returns
// the default value instead.
func (m Of[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := m.kv[key]; ok {
		return value
	}

	return defaultValue
}

// Len returns the number of elements in the map.
func (m Of[K, V]) Len() int {
	return len(*m.ll)
}

// All returns an iterator that yields all elements in the map starting
// at the front (oldest Set element).
func (m Of[K, V]) All() iter.Seq2[K, V] {
	return m.AllFromFront()
}

// AllFromFront returns an iterator that yields all elements in the map starting
// at the front (oldest Set element).
func (m Of[K, V]) AllFromFront() iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		for _, key := range *m.ll {
			if !yield(key, m.kv[key]) {
				return
			}
		}
	}
}

// AllFromBack returns an iterator that yields all elements in the map starting
// at the back (most recent Set element).
func (m Of[K, V]) AllFromBack() iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		for _, key := range slices.Backward(*m.ll) {
			if !yield(key, m.kv[key]) {
				return
			}
		}
	}
}

// Keys returns an iterator that yields all the keys in the map starting at the
// front (oldest Set element). To create a slice containing all the map keys,
// use the slices.Collect function on the returned iterator.
func (m Of[K, V]) Keys() iter.Seq[K] {
	return func(yield func(key K) bool) {
		for _, key := range *m.ll {
			if !yield(key) {
				return
			}
		}
	}
}

// Values returns an iterator that yields all the values in the map starting at
// the front (oldest Set element). To create a slice containing all the map
// values, use the slices.Collect function on the returned iterator.
func (m Of[K, V]) Values() iter.Seq[V] {
	return func(yield func(value V) bool) {
		for _, key := range *m.ll {
			if !yield(m.kv[key]) {
				return
			}
		}
	}
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (m Of[K, V]) Delete(key K) (didDelete bool) {
	_, ok := m.kv[key]
	if ok {
		lli := slices.Index(*m.ll, key)
		*m.ll = slices.Delete(*m.ll, lli, lli+1)
		delete(m.kv, key)
	}

	return ok
}

// Copy returns a new OrderedMap with the same elements.
// Using Copy while there are concurrent writes may mangle the result.
func (m Of[K, V]) Copy() Of[K, V] {
	m2 := MakeWithCapacity[K, V](m.Len())
	for _, key := range *m.ll {
		m2.Set(key, m.kv[key])
	}
	return m2
}

// Filter returns a new OrderedMap with the same elements for each keep(key, value) returns true.
// Using Filter while there are concurrent writes may mangle the result.
func (m Of[K, V]) Filter(keep func(key K, value V) bool) Of[K, V] {
	m2 := MakeWithCapacity[K, V](m.Len())
	for _, key := range *m.ll {
		value := m.kv[key]
		if keep(key, value) {
			m2.Set(key, value)
		}
	}
	return m2
}

// Has checks if a key exists in the map.
func (m Of[K, V]) Has(key K) bool {
	_, exists := m.kv[key]
	return exists
}

// Sort keys in the map.
func (m Of[K, V]) SortFunc(cmp func(a, b K) int) {
	slices.SortFunc(*m.ll, cmp)
}

// Append value to multimap
func AppendMultiMap[K comparable, V any](m Of[K, []V], key K, value V) Of[K, []V] {
	if a, ok := m.Get(key); ok {
		m.Set(key, append(a, value))
	} else {
		m.Set(key, []V{value})
	}

	return m
}
