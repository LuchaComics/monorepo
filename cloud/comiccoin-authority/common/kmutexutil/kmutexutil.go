// Package kmutexutil provides utilities for working with kmutex.
package kmutexutil

import (
	"fmt"

	"github.com/im7mortal/kmutex"
)

// KMutexProvider provides interface for abstracting KMutex generation.
type KMutexProvider interface {
	Acquire(key string)
	Acquiref(format string, a ...any)
	Release(key string)
	Releasef(format string, a ...any)
}

type kMutex struct {
	kMutex *kmutex.Kmutex
}

// NewKMutexProvider constructor that returns the default KMutex generator.
func NewKMutexProvider() KMutexProvider {
	kmux := kmutex.New()
	return &kMutex{kMutex: kmux}
}

// Acquire function blocks the current thread if the lock key is currently locked.
func (u *kMutex) Acquire(k string) {
	u.kMutex.Lock(k)
}

// Acquiref function blocks the current thread if the lock key is currently locked.
func (u *kMutex) Acquiref(format string, a ...any) {
	k := fmt.Sprintf(format, a...)
	u.kMutex.Lock(k)
}

// Release function blocks the current thread if the lock key as currently locked.
func (u *kMutex) Release(k string) {
	u.kMutex.Unlock(k)
}

// Releasef function blocks the current thread if the lock key as currently locked.
func (u *kMutex) Releasef(format string, a ...any) {
	k := fmt.Sprintf(format, a...)
	u.kMutex.Unlock(k)
}
