package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NetworkedBlockStruct struct {
	Pos_X  float32
	Pos_Y  float32
	Pos_Z  float32
	Size_X float32
	Size_Y float32
	Size_Z float32
	Color  uint8
}

type NetworkedLevel struct {
	Blocks []NetworkedBlockStruct
}

func GetLevel(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, Level)
}

func SendLevel(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		panic(err)
	}

	new_level := NetworkedLevel{}
	if err := json.Unmarshal(json_data, &new_level); err != nil {
		ctx.Status(http.StatusBadRequest)
		panic(err)
	}

	Level = new_level

	PlayerUpdateMap.Range(func(key, value any) bool {
		PlayerUpdateMap.Store(key, false)
		return true
	})
}
