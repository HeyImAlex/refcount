package refcount

import (
    "errors"
    "math"
    "sync/atomic"
)

var (
    ErrDestroyed = errors.New("refcount: resource already destroyed")
    ErrReleased  = errors.New("refcount: reference already released")
)

type Reference struct {
    count *int32
    released int32
    destructor func()
}

func New(destructor func()) *Reference {
    count := new(int32)
    *count = 1
    return &Reference{
        count:      count, 
        destructor: destructor,
    }
}

func (r *Reference) Clone() (*Reference, error) {
    if atomic.LoadInt32(&r.released) == 1 {
        return nil, ErrReleased
    }
    if atomic.AddInt32(r.count, 1) < 1 {
        atomic.StoreInt32(r.count, math.MinInt32)
        return nil, ErrDestroyed
    }
    return &Reference{
        count:      r.count, 
        destructor: r.destructor,
    }, nil
}

// Like Clone, but panics on error. Safe to use so long as you know a reference
// hasn't been released yet.
func (r *Reference) MustClone() *Reference {
    ref, err := r.Clone()
    if err != nil {
        panic(err)
    }
    return ref
}

func (r *Reference) Release() error {
    if !atomic.CompareAndSwapInt32(&r.released, 0, 1) {
        return ErrReleased
    }
    if atomic.AddInt32(r.count, -1) == 0 {
        if atomic.CompareAndSwapInt32(r.count, 0, math.MinInt32) {
            if r.destructor != nil {
                r.destructor()
            }
        }
    }
    return nil
}
