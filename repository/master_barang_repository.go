package repository

import (
	"database/sql"
	"roxy/entity"
)

type MstBarangRepository interface {
	Create(barang entity.Barang) (entity.Barang, error)
	List() ([]entity.Barang, error)
	GetByID(id string) (entity.Barang, error)
	GetByName(name string) (entity.Barang, error)
	Update(barang entity.Barang) (entity.Barang, error)
	Delete(id string) error
}

type mstBarangRepository struct {
	db *sql.DB
}

func (b *mstBarangRepository) Create(barang entity.Barang) (entity.Barang, error) {

	err := b.db.QueryRow(`INSERT INTO master_barang (nm_barang, qty, harga) VALUES ($1, $2, $3) RETURNING id_barang`, barang.Nm_barang, barang.Qty, barang.Harga).Scan(&barang.Id_barang)

	if err != nil {
		return entity.Barang{}, err
	}
	return barang, nil
}

func (b *mstBarangRepository) List() ([]entity.Barang, error) {
	var barangs []entity.Barang

	rows, err := b.db.Query(`SELECT id_barang, nm_barang, qty, harga FROM master_barang`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var barang entity.Barang
		err := rows.Scan(&barang.Id_barang, &barang.Nm_barang, &barang.Qty, &barang.Harga)
		if err != nil {
			return nil, err
		}
		barangs = append(barangs, barang)
	}
	return barangs, nil
}

func (b *mstBarangRepository) GetByName(name string) (entity.Barang, error) {
	var barang entity.Barang

	err := b.db.QueryRow(`SELECT id_barang, nm_barang, qty, harga FROM master_barang WHERE nm_barang = $1`, name).Scan(&barang.Id_barang, &barang.Nm_barang, &barang.Qty, &barang.Harga)

	if err != nil {
		return entity.Barang{}, err
	}
	return barang, nil
}

func (b *mstBarangRepository) GetByID(id string) (entity.Barang, error) {
	var barang entity.Barang

	err := b.db.QueryRow(`SELECT id_barang, nm_barang, qty, harga FROM master_barang WHERE id_barang = $1`, id).Scan(&barang.Id_barang, &barang.Nm_barang, &barang.Qty, &barang.Harga)

	if err != nil {
		return entity.Barang{}, err
	}

	return barang, nil

}
func (b *mstBarangRepository) Update(barang entity.Barang) (entity.Barang, error) {

	_, err := b.db.Exec(`UPDATE master_barang SET nm_barang = $2, qty = $3, harga = $4 WHERE id_barang = $1`, barang.Id_barang, barang.Nm_barang, barang.Qty, barang.Harga)

	if err != nil {
		return entity.Barang{}, err
	}

	return barang, nil
}
func (b *mstBarangRepository) Delete(id string) error {

	_, err := b.db.Exec(`DELETE FROM master_barang WHERE id_barang = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

func NewBarangRepository(db *sql.DB) MstBarangRepository {
	return &mstBarangRepository{db: db}
}
