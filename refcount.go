package refcount

import (
    "errors"
    "math"
    "sync/atomic"
)

var ErrReleased = errors.New("refcount: reference already released")

type Reference struct {
    *resource
    released uint32
}

type resource struct {
    count      int32
    destructor func()
}

func New(destructor func()) *Reference {
    return &Reference{
        resource: &resource{
            count:      1,
            destructor: destructor,
        },
    }
}

func (r *Reference) Clone() (*Reference, error) {
    if atomic.LoadUint32(&r.released) == 1 {
        return nil, ErrReleased
    }
    if atomic.AddInt32(&r.count, 1) < 1 {
        atomic.StoreInt32(&r.count, math.MinInt32)
        return nil, ErrReleased
    }
    return &Reference{ resource: r.resource }, nil
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
    if !atomic.CompareAndSwapUint32(&r.released, 0, 1) {
        return ErrReleased
    }
    if atomic.AddInt32(&r.count, -1) == 0 {
        if atomic.CompareAndSwapInt32(&r.count, 0, math.MinInt32) {
            if r.destructor != nil {
                r.destructor()
            }
        }
    }
    return nil
}
