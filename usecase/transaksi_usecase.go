package usecase

import (
	"errors"
	"fmt"
	"roxy/entity"
	"roxy/repository"
)

type TransaksiUsecase interface {
	CreateTransaksiWithDetail(transaksi entity.TransaksiHeader, details []entity.TransaksiDetail) (string, error)
	GetTransaksiByID(idTrans string) (entity.TransaksiHeader, []entity.TransaksiDetail, error)
	UpdateTransaksiWithDetail(transaksi entity.TransaksiHeader, details []entity.TransaksiDetail) error
	DeleteTransaksi(idTrans string) error
}

type transaksiUsecase struct {
	TransaksiRepo repository.TransaksiRepository
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
		if details[i].Harga <= 0 {
			return "", errors.New("harga harus lebih dari 0")
		}
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

func (t *transaksiUsecase) GetTransaksiByID(idTrans string) (entity.TransaksiHeader, []entity.TransaksiDetail, error) {
	transaksi, details, err := t.TransaksiRepo.GetTransaksiByID(idTrans)
	if err != nil {
		return transaksi, details, err
	}

	return transaksi, details, nil
}

func (t *transaksiUsecase) UpdateTransaksiWithDetail(header entity.TransaksiHeader, details []entity.TransaksiDetail) error {
	oldTransaksi, _, err := t.TransaksiRepo.GetTransaksiByID(header.IDTrans)
	if err != nil {
		return fmt.Errorf("Message: %s, ID transaksi: %s", err.Error(), header.IDTrans)
	}

	if header.IDTrans != oldTransaksi.IDTrans {
		return errors.New("ID transaksi tidak boleh diubah")
	}

	if header.TglTrans.IsZero() {
		header.TglTrans = oldTransaksi.TglTrans
	}
	if header.Total == 0 {
		header.Total = oldTransaksi.Total
	}

	var total float64
	for i := range details {
		if details[i].Qty <= 0 {
			return errors.New("qty harus lebih dari 0")
		}
		if details[i].Harga <= 0 {
			return errors.New("harga harus lebih dari 0")
		}

		details[i].Subtotal = details[i].Harga * float64(details[i].Qty)
		total += details[i].Subtotal
	}

	header.Total = total

	err = t.TransaksiRepo.UpdateTransaksiWithDetail(header, details)
	if err != nil {
		return err
	}

	return nil
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

func NewTransaksiUsecase(transaksiRepo repository.TransaksiRepository) TransaksiUsecase {
	return &transaksiUsecase{
		TransaksiRepo: transaksiRepo,
	}
}
