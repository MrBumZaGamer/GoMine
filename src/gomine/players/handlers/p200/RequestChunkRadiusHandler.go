package p200

import (
	"gomine/interfaces"
	"goraklib/server"
	"gomine/utils"
	"gomine/net/packets/p200"
	"gomine/net/packets/data"
	"gomine/players/handlers"
	"math"
)

type RequestChunkRadiusHandler struct {
	*handlers.PacketHandler
}

func NewRequestChunkRadiusHandler() RequestChunkRadiusHandler {
	return RequestChunkRadiusHandler{handlers.NewPacketHandler()}
}

/**
 * Handles the chunk radius requests and initial spawns.
 */
func (handler RequestChunkRadiusHandler) Handle(packet interfaces.IPacket, player interfaces.IPlayer, session *server.Session, server interfaces.IServer) bool {
	if chunkRadiusPacket, ok := packet.(*p200.RequestChunkRadiusPacket); ok {

		player.SetViewDistance(chunkRadiusPacket.Radius)

		player.SendChunkRadiusUpdated(player.GetViewDistance())

		var hasChunksInUse = player.HasAnyChunkInUse()

		server.GetLevelManager().GetDefaultLevel().GetDefaultDimension().RequestChunks(player, 10)

		if !hasChunksInUse {
			player.SetSpawned(true)

			var players = server.GetPlayerFactory().GetPlayers()
			for name, pl := range players {
				if !pl.HasSpawned() {
					delete(players, name)
				}
			}
			player.SendPlayerList(data.ListTypeAdd, players)

			for _, receiver := range server.GetPlayerFactory().GetPlayers() {
				if player != receiver {
					receiver.SendPlayerList(data.ListTypeAdd, map[string]interfaces.IPlayer{player.GetName(): player})

					receiver.SpawnTo(player)
					receiver.SpawnPlayerTo(player)
				}
			}

			var x = int32(math.Floor(float64(player.GetPosition().X))) >> 4
			var z = int32(math.Floor(float64(player.GetPosition().Z))) >> 4
			if !player.GetDimension().IsChunkLoaded(x, z) {
				player.GetDimension().LoadChunk(x, z, func(chunk interfaces.IChunk) {
					player.SpawnToAll()
					player.SpawnPlayerToAll()
				})
			}

			player.UpdateAttributes()
			player.SendSetEntityData(player, player.GetEntityData())

			server.BroadcastMessage(utils.Yellow + player.GetDisplayName() + " has joined the server")
		}

		player.SendPlayStatus(data.StatusSpawn)

		return true
	}

	return false
}