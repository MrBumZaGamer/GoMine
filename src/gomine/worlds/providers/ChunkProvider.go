package providers

import (
	"gomine/interfaces"
	"sync"
	"gomine/worlds/chunks"
)

type ChunkProvider struct {
	generator interfaces.IGenerator
	chunks sync.Map
	requests chan ChunkRequest
}

type ChunkRequest struct {
	function func(interfaces.IChunk)
	x int32
	z int32
}

func newChunkProvider() *ChunkProvider {
	return &ChunkProvider{chunks: sync.Map{}, requests: make(chan ChunkRequest, 4096)}
}

/**
 * Submits a chunk X and Z to get loaded.
 * The function provided gets run as soon as the chunk gets loaded.
 */
func (provider *ChunkProvider) LoadChunk(x, z int32, function func(interfaces.IChunk)) {
	provider.requests <-ChunkRequest{function, x, z}
}

/**
 * Checks if a chunk with the given chunk X and Z is loaded.
*/
func (provider *ChunkProvider) IsChunkLoaded(x, z int32) bool {
	var _, ok = provider.chunks.Load(provider.GetChunkIndex(x, z))
	return ok
}

/**
 * Unloads a chunk with the given chunk X and Z.
 */
func (provider *ChunkProvider) UnloadChunk(x, z int32) {
	if provider.IsChunkLoaded(x, z) {
		provider.chunks.Delete(provider.GetChunkIndex(x, z))
	}
}

/**
 * Sets a new chunk in the provider at the x/z coordinates.
 */
func (provider *ChunkProvider) SetChunk(x, z int32, chunk interfaces.IChunk) {
	provider.chunks.Store(provider.GetChunkIndex(x, z), chunk)
}

/**
 * Returns the chunk in the provider at the x/z coordinates and a bool indicating success.
 */
func (provider *ChunkProvider) GetChunk(x, z int32) (interfaces.IChunk, bool) {
	var chunk, ok = provider.chunks.Load(provider.GetChunkIndex(x, z))
	return chunk.(interfaces.IChunk), ok
}

/**
 * Sets the generator of this provider.
 */
func (provider *ChunkProvider) SetGenerator(generator interfaces.IGenerator) {
	provider.generator = generator
}

/**
 * Returns the generator of this provider.
 */
func (provider *ChunkProvider) GetGenerator() interfaces.IGenerator {
	return provider.generator
}

/**
 * Completes a request made to load a chunk.
 */
func (provider *ChunkProvider) completeRequest(request ChunkRequest) {
	var chunk, ok = provider.GetChunk(request.x, request.z)
	if ok {
		request.function(chunk)
	}
}

/**
 * Generates a new chunk with the given X and Z values.
 */
func (provider *ChunkProvider) GenerateChunk(chunkX, chunkZ int32) {
	var chunk = provider.generator.GetNewChunk(chunks.NewChunk(chunkX, chunkZ))
	provider.SetChunk(chunkX, chunkZ, chunk)
}

/**
 * Gets the chunk index for a certain position in a chunk
 */
func (provider *ChunkProvider) GetChunkIndex(x, z int32) int {
	return int(((int64(x) & 0xffffffff) << 32) | (int64(z) & 0xffffffff))
}