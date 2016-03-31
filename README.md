# refcount

Reference counting package for go.

## Install

```shell
go get github.com/HeyImAlex/refcount
```

## Usage

```go
import "github.com/HeyImAlex/refcount"

// Create a new reference. The function passed in is a destructor that will run
// once the refcount hits zero.
ref := refcount.New(func() { fmt.Println("destructor called!") })

// Clone the reference
anotherRef, err := ref.Clone()

// Cloned references can be cloned themselves
oneMoreRef, err := anotherRef.Clone()

// Release each reference when you're done with it
ref.Release()
anotherRef.Release()

// Release is idempotent, so you may call it multiple times without worrying
anotherRef.Release()
anotherRef.Release()
anotherRef.Release()

// References cannot be cloned once they are released
_, err = ref.Clone() // err == refcount.ErrReleased

// Destructor will be called once all references have been released
oneMoreRef.Release() // destructor called!
```

## Motivation

I wrote it for managing deletion of files stored in an lru cache; when a file is evicted from the cache, we don't want to remove it from disk until everyone using the file is done with it. This could be achieved with `runtime.SetFinalizer`, but I like how explicit references are and it was a lot of fun to write.

## Guarantees

- All operations are safe for concurrent access and wait free.
- The destructor will only be called once, and only when there are no more live references.
- The destructor is run synchronously when the final reference is released.
- References can only be released once.
- No new references can be created after the destructor has run.
- If `Clone` is called after `Release` has returned, the clone will fail. If `Release` and `Clone` execute concurrently, the clone may succeed or it may fail.

## TODO

- Document code
- More comprehensive tests
- Add optional auto-release with `runtime.SetFinalizer`
