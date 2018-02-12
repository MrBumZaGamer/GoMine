package providers

import (
	"sync"
	"gomine/worlds/chunks"
	"libraries/anvil"
	"libraries/nbt"
	"os"
	"strconv"
)

type AnvilChunkProvider struct {
	path string
	regions sync.Map

	*ChunkProvider
}

func NewAnvilChunkProvider(path string) *AnvilChunkProvider {
	var provider = &AnvilChunkProvider{path, sync.Map{}, newChunkProvider()}
	go provider.Process()

	return provider
}

/**
 * Continuously processes chunk requests for chunks that were not yet loaded when requested.
 */
func (provider *AnvilChunkProvider) Process() {
	for {
		var request = <-provider.requests
		if provider.IsChunkLoaded(request.x, request.z) {
			provider.completeRequest(request)
			continue
		}

		var regionX, regionZ = request.x >> 5, request.z >> 5
		if provider.IsRegionLoaded(regionX, regionZ) {

			var region, _ = provider.GetRegion(regionX, regionZ)
			var compression, data = region.GetChunkData(request.x, request.z)
			if len(data) < 1 {
				provider.GenerateChunk(request.x, request.z)
				provider.completeRequest(request)
				return
			}

			var reader = GoNBT.NewNBTReader(data, false, GoNBT.BigEndian)
			var c = reader.ReadIntoCompound(int(compression))

			provider.SetChunk(request.x, request.z, chunks.GetAnvilChunkFromNBT(c))
			provider.completeRequest(request)

		} else {
			var path = provider.path + "r." + strconv.Itoa(int(regionX)) + "." + strconv.Itoa(int(regionZ)) + ".mca"
			var _, err = os.Stat(path)
			if err != nil {
				os.Create(path)
			}
			provider.OpenRegion(regionX, regionZ, path)
		}
	}
}

/**
 * Checks if a region with the given region X and Z is loaded.
 */
func (provider *AnvilChunkProvider) IsRegionLoaded(regionX, regionZ int32) bool {
	var _, ok = provider.regions.Load(provider.GetChunkIndex(regionX, regionZ))
	return ok
}

/**
 * Returns a region with the given region X and Z, or nil if it is not loaded, and a bool indicating success.
 */
func (provider *AnvilChunkProvider) GetRegion(regionX, regionZ int32) (*goanvil.Region, bool) {
	var region, ok = provider.regions.Load(provider.GetChunkIndex(regionX, regionZ))
	return region.(*goanvil.Region), ok
}

/**
 * Sets a region with the given X and Z and path.
 */
func (provider *AnvilChunkProvider) OpenRegion(regionX, regionZ int32, path string) {
	var region, _ = goanvil.OpenRegion(path)
	provider.regions.Store(provider.GetChunkIndex(regionX, regionZ), region)
}

/**
 * Closes the provider and saves all chunks.
 */
func (provider *AnvilChunkProvider) Close(async bool) {
	if async {
		go func() {
			provider.regions.Range(func(index, region interface{}) bool {
				region.(*goanvil.Region).Close(true)
				return true
			})
		}()
	} else {
		provider.regions.Range(func(index, region interface{}) bool {
			region.(*goanvil.Region).Close(true)
			return true
		})
	}
}

/**
 * Saves all regions in the provider.
 */
func (provider *AnvilChunkProvider) Save() {
	go func() {
		provider.regions.Range(func(index, region interface{}) bool {
			region.(*goanvil.Region).Save()
			return true
		})
	}()
}