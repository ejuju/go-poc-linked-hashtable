# Doubly-linked-hashtable in Go

Basically an ordered map. Built as a POC for another project.

## Choosing the number of buckets

The number of buckets is static, in order to avoid needing to grow and shrink the underlying array
on insertion and deletion.
Consequently, the expected number of items in the bucket must be known by the user of the package in
order to set a sensible number of buckets.
