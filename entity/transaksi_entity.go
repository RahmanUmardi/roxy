package entity

import (
	"time"
)

type TransaksiHeader struct {
	IDTrans  string    `json:"id_trans"`
	TglTrans time.Time `json:"tgl_trans"`
	Total    float64   `json:"total"`
}

type TransaksiDetail struct {
	IDTransDetail string  `json:"id_trans_detail"`
	IDTrans       string  `json:"id_trans"`
	IDBarang      string  `json:"id_barang"`
	Qty           int     `json:"qty"`
	Harga         float64 `json:"harga"`
	Subtotal      float64 `json:"subtotal"`
}
