package chunks

import (
	"gomine/interfaces"
	"libraries/nbt"
)

/**
 * Returns a new Anvil chunk from the given NBT compound.
 */
func GetAnvilChunkFromNBT(compound *GoNBT.Compound) interfaces.IChunk {
	var chunk = NewChunk(compound.GetInt("xPos", 0), compound.GetInt("zPos", 0))
	chunk.LightPopulated = getBool(compound.GetByte("LightPopulated", 0))
	chunk.TerrainPopulated = getBool(compound.GetByte("TerrainPopulated", 0))
	chunk.biomes = compound.GetByteArray("Biomes", make([]byte, 256))
	chunk.InhabitedTime = compound.GetLong("InhabitedTime", 0)
	chunk.LastUpdate = compound.GetLong("LastUpdate", 0)
	var heightMap = [256]int16{}
	for i, b := range compound.GetByteArray("HeightMap", make([]byte, 256)) {
		heightMap[i] = int16(b)
	}
	chunk.heightMap = heightMap

	var sections = compound.GetList("Sections", GoNBT.TAG_Compound)
	for _, comp := range sections.GetTags() {
		section := comp.(*GoNBT.Compound)
		subChunk := NewSubChunk()
		subChunk.BlockLight = section.GetByteArray("BlockLight", make([]byte, 2048))
		subChunk.BlockData = section.GetByteArray("Data", make([]byte, 2048))
		subChunk.SkyLight = section.GetByteArray("SkyLight", make([]byte, 2048))
		subChunk.BlockIds = section.GetByteArray("Blocks", make([]byte, 4096))

		chunk.subChunks[section.GetByte("Y", 0)] = subChunk
	}

	return chunk
}

func getBool(value byte) bool {
	if value > 0 {
		return true
	}
	return false
}