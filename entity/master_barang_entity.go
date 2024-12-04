package entity

type Barang struct {
	Id_barang string  `json:"id_barang"`
	Nm_barang string  `json:"nm_barang"`
	Qty       int     `json:"qty"`
	Harga     float32 `json:"harga"`
}
