package main

import (
	goflaker "github.com/MCausc78/goflaker"
)

var (
	Campfire3000Builder    goflaker.SnowflakeBuilder = goflaker.NewBuilder(1)
	Campfire3000Generators map[int]*goflaker.DefaultSnowflakeGenerator
)

func InitializeSnowflakes() {
	goflaker.Initialize()
	Campfire3000Generators = map[int]*goflaker.DefaultSnowflakeGenerator{}
}

func GenerateSnowflake(shardId int) goflaker.Snowflake {
	if g, ok := Campfire3000Generators[shardId]; ok {
		return g.Make(0)
	}
	g := Campfire3000Builder.DefaultGenerator(uint8(shardId & 0xFF))
	Campfire3000Generators[shardId] = g.(*goflaker.DefaultSnowflakeGenerator)
	return g.Make(0)
}

func GenerateSnowflakeValue(shardId int) uint64 {
	return GenerateSnowflake(shardId).Value()
}
