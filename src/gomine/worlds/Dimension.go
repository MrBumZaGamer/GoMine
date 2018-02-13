package worlds

import (
	"gomine/interfaces"
	"gomine/worlds/generation"
	"sync"
	"gomine/worlds/providers"
	"os"
)

const (
	OverworldId = 0
	NetherId    = 1
	EndId	    = 2
)

type Dimension struct {
	name 		string
	dimensionId int
	level       interfaces.ILevel

	chunkProvider interfaces.IChunkProvider
	updatedBlocks map[int][]interfaces.IBlock

	mux sync.Mutex
}

/**
 * Returns a new dimension with the given dimension ID.
 */
func NewDimension(name string, dimensionId int, level *Level, generator string) *Dimension {
	var path = level.GetServer().GetServerPath() + "worlds/" + level.GetName() + "/" + name + "/region/"
	os.MkdirAll(path, os.ModeDir)

	var dimension = &Dimension{
		name:  name,
		dimensionId: dimensionId,
		level: level,
		updatedBlocks: make(map[int][]interfaces.IBlock),
		chunkProvider: providers.NewAnvilChunkProvider(path),
	}

	if len(generator) == 0 {
		dimension.chunkProvider.SetGenerator(generation.GetGeneratorByName(level.server.GetConfiguration().DefaultGenerator))
	} else {
		dimension.chunkProvider.SetGenerator(generation.GetGeneratorByName(generator))
	}

	return dimension
}

/**
 * Returns the dimension ID of this dimension.
 */
func (dimension *Dimension) GetDimensionId() int {
	return dimension.dimensionId
}

/**
 * Returns the name of this dimension.
 */
func (dimension *Dimension) GetName() string {
	return dimension.name
}

/**
 * Returns the level this dimension is in.
 */
func (dimension *Dimension) GetLevel() interfaces.ILevel {
	return dimension.level
}

/**
 * Closes the dimension and saves it.
 */
func (dimension *Dimension) Close(async bool) {
	dimension.chunkProvider.Close(async)
}

/**
 * Saves the dimension.
 */
func (dimension *Dimension) Save() {
	dimension.chunkProvider.Save()
}

/**
 * Returns if a chunk with the given X and Z is loaded.
 */
func (dimension *Dimension) IsChunkLoaded(x, z int32) bool {
	return dimension.chunkProvider.IsChunkLoaded(x, z)
}

/**
 * Sets the chunk at the given X and Z unloaded.
 */
func (dimension *Dimension) UnloadChunk(x, z int32) {
	dimension.chunkProvider.UnloadChunk(x, z)
}

/**
 * Submits a chunk X and Z to get loaded.
 * The function provided gets run as soon as the chunk gets loaded.
 */
func (dimension *Dimension) LoadChunk(x, z int32, function func(interfaces.IChunk)) {
	dimension.chunkProvider.LoadChunk(x, z, function)
}

/**
 * Sets a new chunk in the dimension at the x/z coordinates.
 */
func (dimension *Dimension) SetChunk(x, z int32, chunk interfaces.IChunk) {
	dimension.chunkProvider.SetChunk(x, z, chunk)
}

/**
 * Gets the chunk in the dimension at the x/z coordinates.
 */
func (dimension *Dimension) GetChunk(x, z int32) (interfaces.IChunk, bool) {
	return dimension.chunkProvider.GetChunk(x, z)
}

/**
 * Sets the generator of this dimension.
 */
func (dimension *Dimension) SetGenerator(generator interfaces.IGenerator) {
	dimension.chunkProvider.SetGenerator(generator)
}

/**
 * Returns the generator of this level.
 */
func (dimension *Dimension) GetGenerator() interfaces.IGenerator {
	return dimension.chunkProvider.GetGenerator()
}

func (dimension *Dimension) TickDimension() {

}

/**
 * Sends all chunks required around the player.
 */
func (dimension *Dimension) RequestChunks(player interfaces.IPlayer, distance int32) {
	xD, zD := int32(player.GetPosition().X) >> 4, int32(player.GetPosition().Z) >> 4

	for x := -distance + xD; x <= distance + xD; x++ {
		for z := -distance + zD; z <= distance + zD; z++ {

			var xRel = x - xD
			var zRel = z - zD
			if xRel * xRel + zRel * zRel <= distance * distance {
				index := GetChunkIndex(x, z)

				if player.HasChunkInUse(index) {
					continue
				}

				f := func(c interfaces.IChunk) {
					c.AddViewer(player)
					player.SendChunk(c, index)

					for _, entity := range c.GetEntities() {
						entity.SpawnTo(player)
					}
				}
				if !dimension.chunkProvider.IsChunkLoaded(x, z) {
					dimension.chunkProvider.LoadChunk(x, z, f)
				} else {
					chunk, _ := dimension.GetChunk(x, z)
					f(chunk)
				}
			}
		}
	}
}
