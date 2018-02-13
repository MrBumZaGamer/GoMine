package chunks

type EmptySubChunk struct {
	Blocks []byte
	BlockLight []byte
	SkyLight []byte
	Metadata []byte
}

func NewEmptySubChunk() *EmptySubChunk {
	return &EmptySubChunk{[]byte{}, []byte{}, []byte{}, []byte{}}
}

func (subChunk *EmptySubChunk) IsAllAir() bool {
	return true
}

func (subChunk *EmptySubChunk) GetIdIndex(x, y, z int) int {
	return (x << 8) | (z << 4) | y
}

func (subChunk *EmptySubChunk) GetDataIndex(x, y, z int) int {
	return (x << 7) + (z << 3) + (y >> 1)
}

func (subChunk *EmptySubChunk) GetBlockId(x, y, z int) byte {
	return 0
}

func (subChunk *EmptySubChunk) SetBlockId(x, y, z int, id byte) {

}

func (subChunk *EmptySubChunk) GetBlockLight(x, y, z int) byte {
	return byte(0)
}

func (subChunk *EmptySubChunk) SetBlockLight(x, y, z int, data byte) {

}

func (subChunk *EmptySubChunk) GetSkyLight(x, y, z int) byte {
	return byte(0)
}

func (subChunk *EmptySubChunk) SetSkyLight(x, y, z int, data byte) {

}

func (subChunk *EmptySubChunk) GetBlockData(x, y, z int) byte {
	return byte(0)
}

func (subChunk *EmptySubChunk) SetBlockData(x, y, z int, data byte) {

}

func (subChunk *EmptySubChunk) GetHighestBlockId(x, z int) byte {
	return 0
}

func (subChunk *EmptySubChunk) GetHighestBlockData(x, z int) byte {
	return 0
}

func (subChunk *EmptySubChunk) GetHighestBlock(x, z int) int {
	return 0
}

func (subChunk *EmptySubChunk) ToBinary() []byte {
	return make([]byte, 4096 + 2048 + 1)
}