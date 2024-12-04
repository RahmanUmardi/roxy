package repository

import (
	"database/sql"
	"fmt"
	"roxy/entity"
)

type TransaksiRepository interface {
	CreateTransaksiWithDetail(header entity.TransaksiHeader, details []entity.TransaksiDetail) (string, error)
	GetTransaksiByID(idTrans string) (entity.TransaksiHeader, []entity.TransaksiDetail, error)
	DeleteTransaksi(idTrans string) error
	UpdateTransaksiWithDetail(transaksi entity.TransaksiHeader, details []entity.TransaksiDetail) error
}

type transaksiRepository struct {
	DB *sql.DB
}

func (t *transaksiRepository) CreateTransaksiWithDetail(header entity.TransaksiHeader, details []entity.TransaksiDetail) (string, error) {
	tx, err := t.DB.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var idTransaksi string
	queryHeader := `
        INSERT INTO transaksi_header (tgl_trans, total)
        VALUES ($1, $2) RETURNING id_trans
    `
	err = tx.QueryRow(queryHeader, header.TglTrans, header.Total).Scan(&idTransaksi)
	if err != nil {
		return "", err
	}

	for i := range details {
		details[i].IDTrans = idTransaksi
	}

	for _, detail := range details {
		queryDetail := `
            INSERT INTO transaksi_detail (id_trans, id_barang, qty, harga, subtotal)
            VALUES ($1, $2, $3, $4, $5)
        `
		_, err := tx.Exec(queryDetail, detail.IDTrans, detail.IDBarang, detail.Qty, detail.Harga, detail.Subtotal)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return idTransaksi, nil
}

func (t *transaksiRepository) GetTransaksiByID(idTrans string) (entity.TransaksiHeader, []entity.TransaksiDetail, error) {
	var transaksi entity.TransaksiHeader
	var details []entity.TransaksiDetail

	queryTransaksi := `SELECT id_trans, tgl_trans, total FROM transaksi_header WHERE id_trans = $1`
	row := t.DB.QueryRow(queryTransaksi, idTrans)
	err := row.Scan(&transaksi.IDTrans, &transaksi.TglTrans, &transaksi.Total)
	if err != nil {
		if err == sql.ErrNoRows {
			return transaksi, details, fmt.Errorf("transaksi not found")
		}
		return transaksi, details, err
	}

	queryDetail := `SELECT id_trans_detail, id_trans, id_barang, qty, harga, subtotal FROM transaksi_detail WHERE id_trans = $1`
	rows, err := t.DB.Query(queryDetail, idTrans)
	if err != nil {
		return transaksi, details, err
	}
	defer rows.Close()

	for rows.Next() {
		var detail entity.TransaksiDetail
		err := rows.Scan(&detail.IDTransDetail, &detail.IDTrans, &detail.IDBarang, &detail.Qty, &detail.Harga, &detail.Subtotal)
		if err != nil {
			return transaksi, details, err
		}
		details = append(details, detail)
	}

	return transaksi, details, nil
}

func (t *transaksiRepository) UpdateTransaksiWithDetail(transaksi entity.TransaksiHeader, details []entity.TransaksiDetail) error {
	tx, err := t.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateTransaksi := `UPDATE transaksi_header SET tgl_trans = $1, total = $2 WHERE id_trans = $3`
	_, err = tx.Exec(updateTransaksi, transaksi.TglTrans, transaksi.Total, transaksi.IDTrans)
	if err != nil {
		return err
	}

	deleteDetail := `DELETE FROM transaksi_detail WHERE id_trans = $1`
	_, err = tx.Exec(deleteDetail, transaksi.IDTrans)
	if err != nil {
		return err
	}

	for _, detail := range details {
		insertDetail := `INSERT INTO transaksi_detail (id_trans_detail, id_trans, id_barang, qty, harga, subtotal) 
						VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.Exec(insertDetail, detail.IDTransDetail, detail.IDTrans, detail.IDBarang, detail.Qty, detail.Harga, detail.Subtotal)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (t *transaksiRepository) DeleteTransaksi(idTrans string) error {
	tx, err := t.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deleteDetail := `DELETE FROM transaksi_detail WHERE id_trans = $1`
	_, err = tx.Exec(deleteDetail, idTrans)
	if err != nil {
		return err
	}

	deleteTransaksi := `DELETE FROM transaksi_header WHERE id_trans = $1`
	_, err = tx.Exec(deleteTransaksi, idTrans)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func NewTransaksiRepository(db *sql.DB) TransaksiRepository {
	return &transaksiRepository{DB: db}
}
