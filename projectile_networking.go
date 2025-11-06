package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type NetworkedProjectile struct {
	Pos_X      float32
	Pos_Y      float32
	Pos_Z      float32
	Vel_X      float32
	Vel_Y      float32
	Vel_Z      float32
	Speed      float32
	Damage     int
	Name       string
	ServerTime uint64
}

type CollectionOfNetworkedProjectiles struct {
	Projectiles []NetworkedProjectile
}

var Projectiles sync.Map

func SpawnProjectile(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		panic(err)
	}

	new_projectile := NetworkedProjectile{}
	if err := json.Unmarshal(json_data, &new_projectile); err != nil {
		ctx.Status(http.StatusBadRequest)
		panic(err)
	}

	new_projectile.Pos_X += new_projectile.Vel_X * (float32(ServerData.Time) - float32(new_projectile.ServerTime)) * 60
	new_projectile.Pos_Y += new_projectile.Vel_Y * (float32(ServerData.Time) - float32(new_projectile.ServerTime)) * 60
	new_projectile.Pos_Z += new_projectile.Vel_Z * (float32(ServerData.Time) - float32(new_projectile.ServerTime)) * 60

	Projectiles.Store(GetSyncMapSize(&Projectiles)+1, new_projectile)
}

func GetProjectiles(ctx *gin.Context) {
	var projectiles_to_send CollectionOfNetworkedProjectiles

	Projectiles.Range(func(key, value any) bool {
		projectile := value.(NetworkedProjectile)
		projectiles_to_send.Projectiles = append(projectiles_to_send.Projectiles, projectile)
		return true
	})

	ctx.JSON(http.StatusAccepted, projectiles_to_send)
}
