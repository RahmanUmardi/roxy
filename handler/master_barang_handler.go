package handler

import (
	"net/http"
	"roxy/config"
	"roxy/entity"
	"roxy/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type MasterBarangHandler struct {
	barangUc usecase.MstBarangUseCase
	rg       *gin.RouterGroup
}

func (b *MasterBarangHandler) createHandler(ctx *gin.Context) {
	var payload entity.Barang

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response := struct {
			Message string
			Data    entity.Barang
		}{
			Message: "Invalid Payload for Barang",
			Data:    entity.Barang{},
		}

		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	barang, err := b.barangUc.Create(payload)
	if err != nil {
		response := struct {
			Message string
		}{
			Message: err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := struct {
		Message string
		Data    entity.Barang
	}{
		Message: "Barang Created",
		Data:    barang,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (b *MasterBarangHandler) listHandler(ctx *gin.Context) {

	barangs, err := b.barangUc.List()
	if err != nil {
		response := struct {
			Message string
		}{
			Message: err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if len(barangs) > 0 {
		response := struct {
			Message string
			Data    []entity.Barang
		}{
			Message: "Succes get all barang",
			Data:    barangs,
		}

		ctx.JSON(http.StatusOK, response)
		return
	}
	response := struct {
		Message string
	}{
		Message: "List of barang is empty",
	}

	ctx.JSON(http.StatusOK, response)
}

func (b *MasterBarangHandler) getHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	barang, err := b.barangUc.GetByID(id)
	if err != nil {
		response := struct {
			Message string
		}{
			Message: "Barang of Id " + id + " Not Found",
		}

		ctx.JSON(http.StatusNotFound, response)
		return
	}
	response := struct {
		Message string
		Data    entity.Barang
	}{
		Message: "Succes get barang by id",
		Data:    barang,
	}

	ctx.JSON(http.StatusOK, response)
}

func (b *MasterBarangHandler) updateHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var payload entity.Barang

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response := struct {
			Message string
		}{
			Message: "Invalid Payload for Barang",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	payload.Id_barang = id

	barang, err := b.barangUc.Update(payload)
	if err != nil {
		if strings.Contains(err.Error(), "name already exists") {
			// Specific error for name conflict
			response := struct {
				Message string
			}{
				Message: err.Error(),
			}
			ctx.JSON(http.StatusConflict, response)
			return
		}

		if strings.Contains(err.Error(), "not found") {
			response := struct {
				Message string
			}{
				Message: "Barang with ID " + id + " Not Found",
			}
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := struct {
			Message string
		}{
			Message: "Failed to update barang: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := struct {
		Message string
		Data    entity.Barang
	}{
		Message: "Barang of Id " + id + " Updated",
		Data:    barang,
	}
	ctx.JSON(http.StatusOK, response)
}

func (b *MasterBarangHandler) deleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := b.barangUc.Delete(id)
	if err != nil {
		response := struct {
			Message string
		}{
			Message: "Barang of Id " + id + " Not Found",
		}

		ctx.JSON(http.StatusNotFound, response)
		return
	}
	response := struct {
		Message string
	}{
		Message: "Barang of Id " + id + " Deleted",
	}

	ctx.JSON(http.StatusOK, response)
}

func (b *MasterBarangHandler) Route() {
	b.rg.POST(config.PostBarang, b.createHandler)
	b.rg.GET(config.GetBarangList, b.listHandler)
	b.rg.GET(config.GetBarang, b.getHandler)
	b.rg.PUT(config.PutBarang, b.updateHandler)
	b.rg.DELETE(config.DeleteBarang, b.deleteHandler)
}

func NewBarangHandler(barangUc usecase.MstBarangUseCase, rg *gin.RouterGroup) *MasterBarangHandler {
	return &MasterBarangHandler{barangUc: barangUc, rg: rg}
}
