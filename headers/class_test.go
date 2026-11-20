package headers

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestRetirarExito(t *testing.T) {
	c := NewCuenta("Alice", 234.56)
	ok, err := c.Retirar(200.00)
	if err != nil || !ok {
		t.Fatalf("esperado retiro exitoso, obtenido ok=%v err=%v", ok, err)
	}
	if math.Abs(c.ObtenerSaldo()-34.56) > 0.001 {
		t.Fatalf("saldo esperado 34.56, obtenido %.2f", c.ObtenerSaldo())
	}
}

func TestRetirarSaldoInsuficiente(t *testing.T) {
	c := NewCuenta("Hamza", 198.76)
	ok, err := c.Retirar(200.00)
	if err != nil || ok {
		t.Fatalf("esperado rechazo sin error, obtenido ok=%v err=%v", ok, err)
	}
	if math.Abs(c.ObtenerSaldo()-198.76) > 0.001 {
		t.Fatalf("el saldo no debe cambiar, obtenido %.2f", c.ObtenerSaldo())
	}
}

func TestRetirarNegativo(t *testing.T) {
	c := NewCuenta("Pedro", 300.00)
	ok, err := c.Retirar(-20.00)
	if err == nil || ok {
		t.Fatalf("esperado error por monto negativo, obtenido ok=%v err=%v", ok, err)
	}
	if math.Abs(c.ObtenerSaldo()-300.00) > 0.001 {
		t.Fatalf("el saldo no debe cambiar tras error, obtenido %.2f", c.ObtenerSaldo())
	}
}

func TestAbonarExitoYFallo(t *testing.T) {
	c := NewCuenta("Bob", 100.00)
	ok, err := c.Abonar(50.00)
	if err != nil || !ok {
		t.Fatalf("esperado abono exitoso, obtenido ok=%v err=%v", ok, err)
	}
	if math.Abs(c.ObtenerSaldo()-150.00) > 0.001 {
		t.Fatalf("saldo esperado 150.00, obtenido %.2f", c.ObtenerSaldo())
	}

	ok, err = c.Abonar(-20.00)
	if err == nil || ok {
		t.Fatalf("esperado error por abono negativo, obtenido ok=%v err=%v", ok, err)
	}
	if math.Abs(c.ObtenerSaldo()-150.00) > 0.001 {
		t.Fatalf("el saldo debe permanecer en 150.00, obtenido %.2f", c.ObtenerSaldo())
	}
}

func TestIntegridadTransaccionYLogJSON(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test_tx.json")

	c := NewCuenta("Alice", 234.56)
	txExito, ok, err := c.EjecutarTransaccion(TipoRetiro, 200.00)
	if !ok || err != nil || txExito.Estado != EstadoExito {
		t.Fatalf("esperado EXITO, obtenido estado=%s ok=%v err=%v", txExito.Estado, ok, err)
	}
	if err := GuardarTransaccionJSON(logPath, txExito); err != nil {
		t.Fatalf("error guardando log JSON: %v", err)
	}

	c2 := NewCuenta("Hamza", 198.76)
	txFallo, ok, err := c2.EjecutarTransaccion(TipoRetiro, 200.00)
	if ok || err != nil || txFallo.Estado != EstadoFallo || c2.ObtenerSaldo() != 198.76 {
		t.Fatalf("esperado FALLO atómico, obtenido estado=%s saldo=%.2f", txFallo.Estado, c2.ObtenerSaldo())
	}
	if err := GuardarTransaccionJSON(logPath, txFallo); err != nil {
		t.Fatalf("error guardando log JSON: %v", err)
	}

	c3 := NewCuenta("Pedro", 300.00)
	txError, ok, err := c3.EjecutarTransaccion(TipoRetiro, -20.00)
	if ok || err == nil || txError.Estado != EstadoError || c3.ObtenerSaldo() != 300.00 {
		t.Fatalf("esperado ERROR atómico, obtenido estado=%s saldo=%.2f", txError.Estado, c3.ObtenerSaldo())
	}
	if err := GuardarTransaccionJSON(logPath, txError); err != nil {
		t.Fatalf("error guardando log JSON: %v", err)
	}

	raw, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("error leyendo log JSON: %v", err)
	}
	var registro RegistroLog
	if err := json.Unmarshal(raw, &registro); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if registro.TotalEventos != 3 || len(registro.Transacciones) != 3 {
		t.Fatalf("esperados 3 eventos en log JSON, obtenidos %d", registro.TotalEventos)
	}
}
