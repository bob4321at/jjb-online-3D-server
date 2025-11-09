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
	Health int
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
var PlayerUpdateMap sync.Map

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
	PlayerUpdateMap.Store(new_player_data.ID, false)

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

	playerr, there := Players.Load(new_player_data.ID)
	if there {
		player_health := playerr.(NetworkedPlayer)

		new_player_data.Health = player_health.Health

		Players.Store(new_player_data.ID, new_player_data)
	}
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

	playerr, there := Players.Load(get_player_id.ID)
	if there {
		player_health := playerr.(NetworkedPlayer)

		get_player_id.Health = player_health.Health

		ctx.JSON(http.StatusAccepted, get_player_id)
	}
}

func GetPlayerMapState(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	get_player_id := NetworkedPlayer{}
	if err := json.Unmarshal(json_data, &get_player_id); err != nil {
		panic(err)
	}
	get_player_map_state, exists := PlayerUpdateMap.Load(get_player_id.ID)
	player_map_state := get_player_map_state.(bool)
	if exists {
		if player_map_state {
			get_player_id.Pos_X = 0
		} else {
			get_player_id.Pos_X = -1
			PlayerUpdateMap.Delete(get_player_id.ID)
			PlayerUpdateMap.Store(get_player_id.ID, true)
		}
	}

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

func CheckPlayers(ctx *gin.Context) {
	var players_to_send CollectionOfNetworkedPlayers

	Players.Range(func(key, value any) bool {
		player := value.(NetworkedPlayer)
		players_to_send.Players = append(players_to_send.Players, player)
		return true
	})

	ctx.JSON(http.StatusOK, players_to_send)
}
