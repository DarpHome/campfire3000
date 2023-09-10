package main

import (
	"math/rand"
	"time"
)

var (
	RNG *rand.Rand = nil
)

func InitializeRandom() {
	RNG = rand.New(rand.NewSource(time.Now().UnixNano()))
}
