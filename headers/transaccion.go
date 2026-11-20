/*
Source file: transaccion.go
Entidad Transaccion y Registro JSON para US-04 (Integridad de Transacción).
Date: nov-2026
*/
package headers

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync/atomic"
	"time"
)

// Estados de transacción según US-04
const (
	EstadoExito = "EXITO"
	EstadoFallo = "FALLO"
	EstadoError = "ERROR"
)

// Tipos de operación soportados
const (
	TipoRetiro = "RETIRO"
	TipoAbono  = "ABONO"
)

var txCounter uint64

// Transaccion representa el registro auditable e inmutable de una operación financiera.
type Transaccion struct {
	ID            string  `json:"id"`
	Timestamp     string  `json:"timestamp"`
	Cliente       string  `json:"cliente"`
	Tipo          string  `json:"tipo"`
	Monto         float64 `json:"monto"`
	SaldoAnterior float64 `json:"saldo_anterior"`
	SaldoNuevo    float64 `json:"saldo_nuevo"`
	Estado        string  `json:"estado"`
	Mensaje       string  `json:"mensaje"`
	Atomica       bool    `json:"atomica"`
}

// RegistroLog agrupa el historial de transacciones para serialización JSON.
type RegistroLog struct {
	Proyecto      string        `json:"proyecto"`
	TotalEventos  int           `json:"total_eventos"`
	Transacciones []Transaccion `json:"transacciones"`
}

func generarIDTransaccion() string {
	seq := atomic.AddUint64(&txCounter, 1)
	return fmt.Sprintf("TXN-%d-%03d", time.Now().UnixMilli(), seq)
}

func redondear2(val float64) float64 {
	return math.Round(val*100) / 100
}

// EjecutarTransaccion ejecuta un RETIRO o ABONO de forma atómica (RN-03) y devuelve su registro auditable.
func (c *Cuenta) EjecutarTransaccion(tipo string, monto float64) (Transaccion, bool, error) {
	saldoAnterior := redondear2(c.saldo)
	tx := Transaccion{
		ID:            generarIDTransaccion(),
		Timestamp:     time.Now().Format(time.RFC3339),
		Cliente:       c.nombre,
		Tipo:          tipo,
		Monto:         redondear2(monto),
		SaldoAnterior: saldoAnterior,
		SaldoNuevo:    saldoAnterior,
		Atomica:       true,
	}

	var exito bool
	var err error

	switch tipo {
	case TipoRetiro:
		exito, err = c.Retirar(monto)
		if err != nil {
			// Rollback explícito para garantizar atomicidad (RN-03)
			c.saldo = saldoAnterior
			tx.Estado = EstadoError
			tx.Mensaje = err.Error()
		} else if !exito {
			c.saldo = saldoAnterior
			tx.Estado = EstadoFallo
			tx.Mensaje = "Saldo insuficiente. Intente con un monto menor"
		} else {
			c.saldo = redondear2(c.saldo)
			tx.SaldoNuevo = c.saldo
			tx.Estado = EstadoExito
			tx.Mensaje = "Retiro procesado correctamente"
		}

	case TipoAbono:
		exito, err = c.Abonar(monto)
		if err != nil {
			c.saldo = saldoAnterior
			tx.Estado = EstadoError
			tx.Mensaje = err.Error()
		} else {
			c.saldo = redondear2(c.saldo)
			tx.SaldoNuevo = c.saldo
			tx.Estado = EstadoExito
			tx.Mensaje = "Abono procesado correctamente"
		}

	default:
		err = fmt.Errorf("tipo de operación no soportado: %s", tipo)
		tx.Estado = EstadoError
		tx.Mensaje = err.Error()
	}

	return tx, exito, err
}

// GuardarTransaccionJSON agrega una transacción al archivo JSON especificado, preservando registros previos.
func GuardarTransaccionJSON(rutaArchivo string, tx Transaccion) error {
	var registro RegistroLog
	registro.Proyecto = "prj_withdrawal"

	if data, err := os.ReadFile(rutaArchivo); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &registro)
	}

	registro.Transacciones = append(registro.Transacciones, tx)
	registro.TotalEventos = len(registro.Transacciones)

	salida, err := json.MarshalIndent(registro, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando log JSON: %w", err)
	}

	if err := os.WriteFile(rutaArchivo, salida, 0644); err != nil {
		return fmt.Errorf("error escribiendo log JSON: %w", err)
	}
	return nil
}

// FormatearTransaccionJSON devuelve la representación JSON indentada de una transacción.
func FormatearTransaccionJSON(tx Transaccion) (string, error) {
	b, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
