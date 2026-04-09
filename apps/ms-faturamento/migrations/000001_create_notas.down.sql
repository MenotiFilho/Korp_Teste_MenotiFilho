DROP TRIGGER IF EXISTS trg_notas_set_updated_at ON notas;
DROP TABLE IF EXISTS notas;
DROP FUNCTION IF EXISTS set_updated_at();
DROP SEQUENCE IF EXISTS nota_numero_seq;
