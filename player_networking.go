package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type NetworkedPlayer struct {
	Pos_X  float32
	Pos_Y  float32
	Pos_Z  float32
	Health uint8
	ID     uint8
}
type CollectionOfNetworkedPlayers struct {
	Players []NetworkedPlayer
}

type PlayerAndProjectileNetworked struct {
	Player     NetworkedPlayer
	Projectile NetworkedProjectile
}

var Players sync.Map
var PlayerDamageMap sync.Map

func GetSyncMapSize(m *sync.Map) int {
	count := 0
	m.Range(func(key, value any) bool {
		count++
		return true // Continue iteration
	})
	return count
}

func AddPlayer(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		panic(err)
	}

	new_player_data := NetworkedPlayer{}
	if err := json.Unmarshal(json_data, &new_player_data); err != nil {
		ctx.Status(http.StatusBadRequest)
		panic(err)
	}

	new_player_data.ID = uint8(GetSyncMapSize(&Players) + 1)
	Players.Store(new_player_data.ID, new_player_data)

	ctx.JSON(http.StatusAccepted, new_player_data)
}

func GetOtherPlayers(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		panic(err)
	}

	new_player_data := NetworkedPlayer{}
	if err := json.Unmarshal(json_data, &new_player_data); err != nil {
		ctx.Status(http.StatusBadRequest)
		panic(err)
	}

	var players_to_send CollectionOfNetworkedPlayers

	Players.Range(func(key, value any) bool {
		player := value.(NetworkedPlayer)
		if player.ID != new_player_data.ID {
			players_to_send.Players = append(players_to_send.Players, player)
		}
		return true
	})

	ctx.JSON(http.StatusAccepted, players_to_send)
}

func UpdatePlayerPos(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	new_player_data := NetworkedPlayer{}
	if err := json.Unmarshal(json_data, &new_player_data); err != nil {
		panic(err)
	}

	Players.Store(new_player_data.ID, new_player_data)
}

func GetPlayerHealth(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	get_player_id := NetworkedPlayer{}
	if err := json.Unmarshal(json_data, &get_player_id); err != nil {
		panic(err)
	}

	player, _ := Players.Load(get_player_id.ID)
	player_health := player.(NetworkedPlayer)

	get_player_id.Health = player_health.Health

	ctx.JSON(http.StatusAccepted, get_player_id)
}

func DamagePlayer(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	new_data := PlayerAndProjectileNetworked{}
	if err := json.Unmarshal(json_data, &new_data); err != nil {
		panic(err)
	}

	PlayerDamageMap.Store(new_data, new_data)
}
