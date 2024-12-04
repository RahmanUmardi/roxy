CREATE DATABASE roxy;

-- CREATE MASTER_BARANG
CREATE TABLE master_barang (
    id_barang VARCHAR(15) PRIMARY KEY,
    nm_barang VARCHAR(30) NOT NULL,
    qty INT NOT NULL,
    harga DOUBLE PRECISION NOT NULL
);

CREATE SEQUENCE barang_seq START 1 INCREMENT 1;

CREATE OR REPLACE FUNCTION generate_barang_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id_barang := 'BR-' || LPAD(nextval('barang_seq')::TEXT, 4, '0');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_generate_barang_id
BEFORE INSERT ON master_barang
FOR EACH ROW
WHEN (NEW.id_barang IS NULL)
EXECUTE FUNCTION generate_barang_id();

-- CREATE TRANSAKSI HEADER
CREATE TABLE transaksi_header (
    id_trans VARCHAR(15) PRIMARY KEY,
    tgl_trans TIMESTAMP,
    total DOUBLE PRECISION NOT NULL
);

CREATE SEQUENCE transaksi_seq START 1 INCREMENT 1;

-- Fungsi untuk generate id_trans otomatis
CREATE OR REPLACE FUNCTION generate_transaksi_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id_trans := 'TR-' || LPAD(nextval('transaksi_seq')::TEXT, 4, '0');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Trigger untuk generate id_trans sebelum insert
CREATE TRIGGER trg_generate_transaksi_id
BEFORE INSERT ON transaksi_header
FOR EACH ROW
WHEN (NEW.id_trans IS NULL)
EXECUTE FUNCTION generate_transaksi_id();

-- CREATE TRANSAKSI DETAIL
CREATE TABLE transaksi_detail (
    id_trans_detail VARCHAR(15) PRIMARY KEY,
    id_trans VARCHAR(15) REFERENCES transaksi_header(id_trans) ON DELETE CASCADE,
    id_barang VARCHAR(15) REFERENCES master_barang(id_barang) ON DELETE CASCADE,
    qty INT NOT NULL,
    harga DOUBLE PRECISION NOT NULL,
    subtotal DOUBLE PRECISION NOT NULL
);

-- Sequence untuk transaksi_detail
CREATE SEQUENCE transaksi_detail_seq START 1 INCREMENT 1;

-- Function untuk generate id_trans_detail otomatis
CREATE OR REPLACE FUNCTION generate_transaksi_detail_id()
RETURNS TRIGGER AS $$ 
BEGIN
    NEW.id_trans_detail := 'TD-' || LPAD(nextval('transaksi_detail_seq')::TEXT, 4, '0');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk generate id_trans_detail otomatis
CREATE TRIGGER trg_generate_transaksi_detail_id
BEFORE INSERT ON transaksi_detail
FOR EACH ROW
WHEN (NEW.id_trans_detail IS NULL)
EXECUTE FUNCTION generate_transaksi_detail_id();

CREATE OR REPLACE FUNCTION update_total_transaksi()
RETURNS TRIGGER AS $$
BEGIN
    -- Menghitung total transaksi berdasarkan subtotal dari transaksi_detail
    UPDATE transaksi_header
    SET total = (SELECT SUM(subtotal) FROM transaksi_detail WHERE id_trans = NEW.id_trans)
    WHERE id_trans = NEW.id_trans;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_total_transaksi_after_update()
RETURNS TRIGGER AS $$
BEGIN
    -- Menghitung total transaksi setelah update pada transaksi_detail
    UPDATE transaksi_header
    SET total = (SELECT SUM(subtotal) FROM transaksi_detail WHERE id_trans = NEW.id_trans)
    WHERE id_trans = NEW.id_trans;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_total_transaksi_after_delete()
RETURNS TRIGGER AS $$
BEGIN
    -- Menghitung total transaksi setelah detail transaksi dihapus
    UPDATE transaksi_header
    SET total = (SELECT SUM(subtotal) FROM transaksi_detail WHERE id_trans = OLD.id_trans)
    WHERE id_trans = OLD.id_trans;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- TRIGER FOR INSERT
CREATE TRIGGER trg_update_total_transaksi_insert
AFTER INSERT ON transaksi_detail
FOR EACH ROW
EXECUTE FUNCTION update_total_transaksi();

-- TRIGER FOR UPDATE
CREATE TRIGGER trg_update_total_transaksi_update
AFTER UPDATE ON transaksi_detail
FOR EACH ROW
EXECUTE FUNCTION update_total_transaksi_after_update();

-- TRIGER FOR DELETE
CREATE TRIGGER trg_update_total_transaksi_delete
AFTER DELETE ON transaksi_detail
FOR EACH ROW
EXECUTE FUNCTION update_total_transaksi_after_delete();


CREATE OR REPLACE FUNCTION update_stok_barang()
RETURNS TRIGGER AS $$
BEGIN
    -- Mengurangi stok barang sesuai dengan jumlah transaksi
    UPDATE master_barang
    SET qty = qty - NEW.qty
    WHERE id_barang = NEW.id_barang;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_stok_barang
AFTER INSERT ON transaksi_detail
FOR EACH ROW
EXECUTE FUNCTION update_stok_barang();
