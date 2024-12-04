package handler

import (
	"net/http"
	"roxy/config"
	"roxy/entity"
	"roxy/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

type TransaksiHandler struct {
	TransaksiUsecase usecase.TransaksiUsecase
	rg               *gin.RouterGroup
}

func (t *TransaksiHandler) CreateTransaksiHandler(c *gin.Context) {
	var req struct {
		Header struct {
			TanggalTransaksi string `json:"tanggal_transaksi"`
		} `json:"header"`
		Detail []entity.TransaksiDetail `json:"detail"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	tglTrans, err := time.Parse("2006-01-02", req.Header.TanggalTransaksi)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	header := entity.TransaksiHeader{
		TglTrans: tglTrans,
	}

	idTransaksi, err := t.TransaksiUsecase.CreateTransaksiWithDetail(header, req.Detail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaksi berhasil dibuat",
		"data": gin.H{
			"id_trans":          idTransaksi,
			"tanggal_transaksi": tglTrans.Format("2006-01-02"),
			"detail":            req.Detail,
		},
	})
}

func (t *TransaksiHandler) GetTransaksiHandler(c *gin.Context) {
	idTrans := c.Param("id")

	header, detail, err := t.TransaksiUsecase.GetTransaksiByID(idTrans)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"header": header,
		"detail": detail,
	}

	c.JSON(http.StatusOK, response)
}

func (t *TransaksiHandler) UpdateTransaksiHandler(c *gin.Context) {
	var req struct {
		Header entity.TransaksiHeader   `json:"header"`
		Detail []entity.TransaksiDetail `json:"detail"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.Header.IDTrans == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID transaksi tidak boleh kosong"})
		return
	}

	err := t.TransaksiUsecase.UpdateTransaksiWithDetail(req.Header, req.Detail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil diperbarui"})
}

func (t *TransaksiHandler) DeleteTransaksiHandler(c *gin.Context) {
	idTrans := c.Param("id")

	err := t.TransaksiUsecase.DeleteTransaksi(idTrans)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil dihapus"})
}

func (t *TransaksiHandler) Route() {
	t.rg.POST(config.PostTransaksi, t.CreateTransaksiHandler)
	t.rg.GET(config.GetTransaksiByID, t.GetTransaksiHandler)
	t.rg.PUT(config.PutTransaksi, t.UpdateTransaksiHandler)
	t.rg.DELETE(config.DeleteTransaksi, t.DeleteTransaksiHandler)
}

func NewTransaksiHandler(transaksiUc usecase.TransaksiUsecase, rg *gin.RouterGroup) *TransaksiHandler {
	return &TransaksiHandler{
		TransaksiUsecase: transaksiUc, rg: rg}
}
