package refcount

import (
    "sync"
    "sync/atomic"
    "testing"
)

func TestReleaseOnce(t *testing.T) {
    r := New(nil)
    if err := r.Release(); err != nil {
        t.Fatalf("unexpected error on release: %s", err)
    }
    if err := r.Release(); err != ErrReleased {
        t.Fatalf("expected ErrReleased, actual: %s", nil)
    }
    if err := r.Release(); err != ErrReleased {
        t.Fatalf("expected ErrReleased, actual: %s", nil)
    }
}

func TestDestroyOnce(t *testing.T) {

    var callCount int32
    destructor := func() {
        atomic.AddInt32(&callCount, 1)
    }

    root := New(destructor)
    references := make([]*Reference, 25)
    for i := range references {
        references[i] = root.MustClone()
    }
    if err := root.Release(); err != nil {
        t.Fatalf("unexpected error releasing root reference: %s", err)
    }

    if callCount != 0 {
        t.Fatal("destructor called before all references released")
    }

    var wg sync.WaitGroup
    wg.Add(25)
    for i := range references {
        go func(idx int) {
            defer wg.Done()
            if err := references[idx].Release(); err != nil {
                t.Fatalf("unexpected error releasing reference #%d: %s", idx, err)
            }
        }(i)
    }
    wg.Wait()

    if callCount != 1 {
        t.Fatalf("expected destructor to be called 1 time, actual: %d", callCount)
    }
}
