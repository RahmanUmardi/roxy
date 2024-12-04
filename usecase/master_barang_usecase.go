package usecase

import (
	"fmt"
	"roxy/entity"
	"roxy/repository"
	"strings"
)

type MstBarangUseCase interface {
	Create(barang entity.Barang) (entity.Barang, error)
	List() ([]entity.Barang, error)
	GetByID(id string) (entity.Barang, error)
	GetByName(name string) (entity.Barang, error)
	Update(barang entity.Barang) (entity.Barang, error)
	Delete(id string) error
}

type mstBarangUseCase struct {
	barangRepository repository.MstBarangRepository
}

func (b *mstBarangUseCase) Create(barang entity.Barang) (entity.Barang, error) {
	existBarang, _ := b.barangRepository.GetByName(barang.Nm_barang)
	if strings.TrimSpace(barang.Nm_barang) == "" {
		return entity.Barang{}, fmt.Errorf("name cannot be empty")
	}
	if existBarang.Nm_barang == barang.Nm_barang {
		return entity.Barang{}, fmt.Errorf("name already exist")
	}

	return b.barangRepository.Create(barang)
}

func (b *mstBarangUseCase) List() ([]entity.Barang, error) {
	return b.barangRepository.List()
}

func (b *mstBarangUseCase) GetByID(id string) (entity.Barang, error) {
	return b.barangRepository.GetByID(id)
}

func (b *mstBarangUseCase) GetByName(name string) (entity.Barang, error) {
	return b.barangRepository.GetByName(name)
}

func (b *mstBarangUseCase) Update(barang entity.Barang) (entity.Barang, error) {
	payload, err := b.barangRepository.GetByID(barang.Id_barang)
	if err != nil {
		return entity.Barang{}, fmt.Errorf("barang with ID %s not found", barang.Id_barang)
	}

	if strings.TrimSpace(barang.Nm_barang) != "" && barang.Nm_barang != payload.Nm_barang {
		existBarang, _ := b.barangRepository.GetByName(barang.Nm_barang)
		if existBarang.Id_barang != "" && existBarang.Id_barang != barang.Id_barang {
			return entity.Barang{}, fmt.Errorf("name %s already exists", barang.Nm_barang)
		}
	}

	if strings.TrimSpace(barang.Nm_barang) == "" {
		barang.Nm_barang = payload.Nm_barang
	}

	if barang.Qty == 0 {
		barang.Qty = payload.Qty
	}
	if barang.Harga == 0 {
		barang.Harga = payload.Harga
	}

	updatedBarang, err := b.barangRepository.Update(barang)
	if err != nil {
		return entity.Barang{}, fmt.Errorf("failed to update barang: %v", err)
	}

	return updatedBarang, nil
}

func (b *mstBarangUseCase) Delete(id string) error {
	_, err := b.barangRepository.GetByID(id)
	if err != nil {
		return fmt.Errorf("barang with ID %s not found", id)
	}

	err = b.barangRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete barang: %v", err)
	}

	return nil
}

func NewBarangUseCase(barangRepository repository.MstBarangRepository) MstBarangUseCase {
	return &mstBarangUseCase{barangRepository: barangRepository}
}
