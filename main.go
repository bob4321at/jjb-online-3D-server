package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerDataStruct struct {
	Time uint64
}

var ServerData = ServerDataStruct{}

func GetServerData(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, ServerData)
}

func main() {
	r := gin.Default()

	go func() {
		for true {
			time.Sleep(time.Second / 60)
			ServerData.Time += 1

			Projectiles.Range(func(key, value any) bool {
				old_projectile, _ := Projectiles.Load(key)
				new_projectile := old_projectile.(NetworkedProjectile)
				new_projectile.Pos_X += new_projectile.Vel_X
				new_projectile.Pos_Y += new_projectile.Vel_Y
				new_projectile.Pos_Z += new_projectile.Vel_Z

				Projectiles.Store(key, new_projectile)

				return true
			})

			PlayerDamageMap.Range(func(key, value any) bool {
				player_and_projectile := value.(PlayerAndProjectileNetworked)

				old_player, _ := Players.Load(player_and_projectile.Player.ID)

				new_player := old_player.(NetworkedPlayer)
				new_player.Health -= player_and_projectile.Projectile.Damage

				Players.Store(new_player.ID, new_player)

				PlayerDamageMap.Delete(key)
				return true
			})

			Players.Range(func(key, value any) bool {
				fmt.Println(value)
				return true
			})
		}
	}()

	r.GET("GetServerData", GetServerData)

	r.POST("AddPlayer", AddPlayer)
	r.POST("GetOtherPlayers", GetOtherPlayers)
	r.POST("UpdatePlayerPos", UpdatePlayerPos)
	r.POST("DamagePlayer", DamagePlayer)
	r.POST("GetPlayerHealth", GetPlayerHealth)

	r.POST("SpawnProjectile", SpawnProjectile)
	r.GET("GetProjectiles", GetProjectiles)

	r.Run()
}
