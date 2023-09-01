package lht

import "bytes"

type LHT struct {
	count          int
	oldest, latest *Item
	buckets        []*Item
}

type Item struct {
	key            []byte
	value          any
	previous, next *Item
	nextInBucket   *Item
}

func NewLHT(numBuckets int) *LHT { return &LHT{buckets: make([]*Item, numBuckets)} }

func (lht *LHT) Count() int    { return lht.count }
func (lht *LHT) Oldest() *Item { return lht.oldest }
func (lht *LHT) Latest() *Item { return lht.latest }

func (it *Item) Key() []byte     { return it.key }
func (it *Item) Value() any      { return it.value }
func (it *Item) Next() *Item     { return it.next }
func (it *Item) Previous() *Item { return it.previous }

func (lht *LHT) Put(key []byte, value any) {
	bucketIndex := lht.hashFNV1a(key)
	root := lht.buckets[bucketIndex]
	var prevInBucket *Item
	for item := root; item != nil; prevInBucket, item = item, item.nextInBucket {
		if bytes.Equal(item.key, key) {
			// Replace existing value and move to end of linked-list
			item.value = value
			if item != lht.latest {
				if item == lht.oldest {
					lht.oldest = item.next
				}
				item.next.previous, item.next = item.previous, nil
				item.previous, lht.latest.next = lht.latest, item // link to previous item
				lht.latest = item                                 // set latest to current item
			}
			return
		}
	}

	// Add new item to bucket and increment count
	lht.count++
	newItem := &Item{key: key, value: value}
	if prevInBucket == nil {
		lht.buckets[bucketIndex] = newItem
	} else {
		prevInBucket.nextInBucket = newItem
	}

	// Add at the end of linked-list
	if lht.latest == nil {
		lht.oldest, lht.latest = newItem, newItem
	} else {
		newItem.previous, lht.latest.next = lht.latest, newItem // link to previous item
		lht.latest = newItem                                    // set latest to new item
	}
}

func (lht *LHT) Delete(key []byte) {
	bucketIndex := lht.hashFNV1a(key)
	root := lht.buckets[bucketIndex]
	var prevInBucket *Item
	for item := root; item != nil; prevInBucket, item = item, item.nextInBucket {
		if bytes.Equal(item.key, key) {
			// Remove from bucket and decrement count
			lht.count--
			if prevInBucket != nil {
				prevInBucket.nextInBucket = item.nextInBucket
			}

			// Unlink from linked-list
			if item.previous == nil {
				lht.oldest = item.next
			} else {
				item.previous.next = item.next
			}
			if item.next == nil {
				lht.latest = item.previous
			} else {
				item.next.previous = item.previous
			}
			return
		}
	}
}

func (lht *LHT) Get(key []byte) *Item {
	root := lht.buckets[lht.hashFNV1a(key)]
	for item := root; item != nil; item = item.nextInBucket {
		if bytes.Equal(item.key, key) {
			return item
		}
	}
	return nil
}

func (lht *LHT) hashFNV1a(key []byte) int {
	const offset, prime = uint64(14695981039346656037), uint64(1099511628211) // fnv-1a constants
	hash := offset
	for _, char := range key {
		hash *= prime
		hash ^= uint64(char)
	}
	index := int(hash % uint64(len(lht.buckets)))
	return index
}
