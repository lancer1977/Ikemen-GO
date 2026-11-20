package main

import "sync"

// newTestSystem returns a *System with the pointer fields that production
// code dereferences unconditionally (bgm, sel, loader, selMutex, loadMutex,
// statePool, savePool, loadPool) initialized, mirroring the guarantees the
// package-level `sys` variable gets from its own initializer. Tests that
// build ad-hoc System values (rather than saving/restoring the real `sys`)
// should use this instead of a bare `&System{}`/`System{}` literal so those
// pointer fields are never left nil.
func newTestSystem() *System {
	return &System{
		bgm:       newBgm(),
		sel:       newSelect(),
		loader:    newLoader(),
		selMutex:  newTestRWMutex(),
		loadMutex: newTestMutex(),
		statePool: NewGameStatePool(),
		savePool:  NewGameStatePool(),
		loadPool:  NewGameStatePool(),
	}
}

// newTestRWMutex and newTestMutex let other test files fill in the
// System.selMutex/loadMutex pointer fields directly in a struct literal
// without each of them needing their own "sync" import.
func newTestRWMutex() *sync.RWMutex { return &sync.RWMutex{} }
func newTestMutex() *sync.Mutex     { return &sync.Mutex{} }
