package main

import (
	"log"

	"github.com/ejuju/go-poc-linked-hashtable/lht"
)

func main() {
	lhtable := lht.NewLHT(1)
	lhtable.Put([]byte{2}, "2")
	lhtable.Put([]byte{1}, "1")
	lhtable.Put([]byte{2}, "22")

	for item := lhtable.Oldest(); item != nil; item = item.Next() {
		log.Println(item)
	}

	lhtable.Delete([]byte{1})
	lhtable.Delete([]byte{2})
}
