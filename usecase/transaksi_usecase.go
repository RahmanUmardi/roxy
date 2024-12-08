package usecase

import (
	"errors"
	"fmt"
	"roxy/entity"
	"roxy/repository"
)

type TransaksiUsecase interface {
	CreateTransaksiWithDetail(transaksi entity.TransaksiHeader, details []entity.TransaksiDetail) (string, error)
	GetAllTransaksi() ([]entity.TransaksiHeader, error)
	GetTransaksiByID(idTrans string) (entity.TransaksiHeader, []entity.TransaksiDetail, error)
	UpdateTransaksiWithDetail(idTrans string, transaksi entity.TransaksiHeader, details []entity.TransaksiDetail) (entity.TransaksiHeader, []entity.TransaksiDetail, error)
	DeleteTransaksi(idTrans string) error
}

type transaksiUsecase struct {
	TransaksiRepo repository.TransaksiRepository
	barangRepo    repository.MstBarangRepository
}

func (t *transaksiUsecase) CreateTransaksiWithDetail(transaksi entity.TransaksiHeader, details []entity.TransaksiDetail) (string, error) {
	if len(details) == 0 {
		return "", errors.New("transaksi detail tidak boleh kosong")
	}

	var total float64
	for i := range details {
		if details[i].Qty <= 0 {
			return "", errors.New("qty harus lebih dari 0")
		}

		barang, err := t.barangRepo.GetByID(details[i].IDBarang)
		if err != nil {
			return "", fmt.Errorf("gagal mendapatkan data barang dengan ID %s: %v", details[i].IDBarang, err)
		}

		details[i].Harga = float64(barang.Harga)

		details[i].Subtotal = details[i].Harga * float64(details[i].Qty)

		total += details[i].Subtotal
	}

	transaksi.Total = total

	transaksi.IDTrans = ""

	idTransaksi, err := t.TransaksiRepo.CreateTransaksiWithDetail(transaksi, details)
	if err != nil {
		return "", err
	}

	return idTransaksi, nil
}

func (t *transaksiUsecase) GetAllTransaksi() ([]entity.TransaksiHeader, error) {
	return t.TransaksiRepo.GetAllTransaksi()
}

func (t *transaksiUsecase) GetTransaksiByID(idTrans string) (entity.TransaksiHeader, []entity.TransaksiDetail, error) {
	transaksi, details, err := t.TransaksiRepo.GetTransaksiByID(idTrans)
	if err != nil {
		return transaksi, details, err
	}

	return transaksi, details, nil
}

func (t *transaksiUsecase) UpdateTransaksiWithDetail(idTrans string, header entity.TransaksiHeader, details []entity.TransaksiDetail) (entity.TransaksiHeader, []entity.TransaksiDetail, error) {
	oldTransaksi, _, err := t.TransaksiRepo.GetTransaksiByID(idTrans)
	if err != nil {
		return header, details, fmt.Errorf("Message: %s, ID transaksi: %s", err.Error(), header.IDTrans)
	}

	if header.TglTrans.IsZero() {
		header.TglTrans = oldTransaksi.TglTrans
	}

	var total float64
	for i := range details {
		if details[i].Qty <= 0 {
			return header, details, fmt.Errorf("qty harus lebih dari 0")
		}

		barang, err := t.barangRepo.GetByID(details[i].IDBarang)
		if err != nil {
			return header, details, fmt.Errorf("gagal mendapatkan data barang dengan ID %s: %v", details[i].IDBarang, err)
		}

		details[i].Subtotal = float64(barang.Harga) * float64(details[i].Qty)
		total += details[i].Subtotal
	}

	header.Total = total

	header, details, err = t.TransaksiRepo.UpdateTransaksiWithDetail(header, details)
	if err != nil {
		return header, details, err
	}

	return header, details, nil
}

func (t *transaksiUsecase) DeleteTransaksi(idTrans string) error {
	_, _, err := t.TransaksiRepo.GetTransaksiByID(idTrans)
	if err != nil {
		return fmt.Errorf("transaksi dengan id %s tidak ditemukan", idTrans)
	}

	err = t.TransaksiRepo.DeleteTransaksi(idTrans)
	if err != nil {
		return err
	}

	return nil
}

func NewTransaksiUsecase(transaksiRepo repository.TransaksiRepository, barangRepo repository.MstBarangRepository) TransaksiUsecase {
	return &transaksiUsecase{
		TransaksiRepo: transaksiRepo,
		barangRepo:    barangRepo,
	}
}
