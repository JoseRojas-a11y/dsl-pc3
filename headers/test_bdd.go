package headers

import (
	"fmt"
	"math"
)

func imprimirSaldo(cuenta *Cuenta) {
	fmt.Printf("Saldo de %.2f €\n", cuenta.ObtenerSaldo())
}

// --- Escenarios de Retiro ---

func TestEscenario1Exito(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 1: Successful withdrawal within balance (%s) ---\n", nombre)
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	if err != nil {
		fmt.Printf("Resultado: FALLO. Error inesperado: %v\n\n", err)
		return
	}

	if exito {
		fmt.Println("Resultado: ÉXITO. La retirada fue procesada.")
		imprimirSaldo(cuenta)
		saldoEsperado := saldoInicial - cantidadRetirar
		if math.Abs(cuenta.ObtenerSaldo()-saldoEsperado) < 0.001 {
			fmt.Printf("Validación: El saldo restante es correcto (%.2f).\n\n", saldoEsperado)
		} else {
			fmt.Printf("ERROR: El saldo restante no coincide.\n\n")
		}
	} else {
		fmt.Printf("Resultado: FALLO. La retirada fue rechazada (incorrecto).\n\n")
	}
}

func TestEscenario2Fallo(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 2: Declined withdrawal in excess of balance (%s) ---\n", nombre)
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	if err != nil {
		fmt.Printf("Resultado: FALLO. Error inesperado: %v\n\n", err)
		return
	}

	if !exito {
		fmt.Println("Resultado: ÉXITO. La retirada fue rechazada correctamente.")
		imprimirSaldo(cuenta)
		if math.Abs(cuenta.ObtenerSaldo()-saldoInicial) < 0.001 {
			fmt.Printf("Validación: El saldo no ha cambiado (%.2f).\n\n", saldoInicial)
		} else {
			fmt.Printf("ERROR: El saldo se modificó incorrectamente.\n\n")
		}
	} else {
		fmt.Printf("Resultado: FALLO. La retirada fue procesada.\n\n")
	}
}

func TestEscenario3Error(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 3: Validation Error (%s) ---\n", nombre)
	cuenta := NewCuenta(nombre, saldoInicial)
	exito, err := cuenta.Retirar(cantidadRetirar)

	if err != nil {
		fmt.Printf("Resultado: ÉXITO. El sistema atrapó el error correctamente: %v\n", err)
		imprimirSaldo(cuenta)
		fmt.Println()
	} else if !exito {
		fmt.Printf("Resultado: RECHAZADO. Fondos insuficientes (no es error de validación).\n\n")
	} else {
		fmt.Printf("Resultado: FALLO. Se procesó un retiro que debía dar error.\n\n")
	}
}

// --- Escenarios de Abono (Nuevo Feature) ---

func TestEscenarioAbonoExito(nombre string, saldoInicial float64, cantidadAbonar float64) {
	fmt.Printf("--- Escenario Abono: Successful deposit (%s) ---\n", nombre)
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)
	fmt.Printf("When: %s intenta abonar %.2f\n", nombre, cantidadAbonar)
	exito, err := cuenta.Abonar(cantidadAbonar)

	if err != nil {
		fmt.Printf("Resultado: FALLO. Error inesperado: %v\n\n", err)
		return
	}

	if exito {
		fmt.Println("Resultado: ÉXITO. El abono fue procesado.")
		imprimirSaldo(cuenta)
		saldoEsperado := saldoInicial + cantidadAbonar
		if math.Abs(cuenta.ObtenerSaldo()-saldoEsperado) < 0.001 {
			fmt.Printf("Validación: El saldo actualizado es correcto (%.2f).\n\n", saldoEsperado)
		} else {
			fmt.Printf("ERROR: El saldo actualizado no coincide.\n\n")
		}
	}
}

func TestEscenarioAbonoFallo(nombre string, saldoInicial float64, cantidadAbonar float64) {
	fmt.Printf("--- Escenario Abono: Declined invalid deposit (%s) ---\n", nombre)
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)
	fmt.Printf("When: %s intenta abonar %.2f\n", nombre, cantidadAbonar)
	exito, err := cuenta.Abonar(cantidadAbonar)

	if err != nil {
		fmt.Printf("Resultado: ÉXITO. El abono fue rechazado correctamente: %v\n", err)
		imprimirSaldo(cuenta)
		if math.Abs(cuenta.ObtenerSaldo()-saldoInicial) < 0.001 {
			fmt.Printf("Validación: El saldo se mantiene intacto (%.2f).\n\n", saldoInicial)
		} else {
			fmt.Printf("ERROR: El saldo se modificó indebidamente.\n\n")
		}
	} else if exito {
		fmt.Printf("Resultado: FALLO. Se procesó un abono inválido.\n\n")
	}
}

// --- Escenarios de Integridad y Log JSON (US-04) ---

// DefaultLogFile es el archivo por defecto para registrar transacciones en formato JSON.
var DefaultLogFile = "transacciones.json"

func registrarYMostrarLog(tx Transaccion, rutaLog string, mostrarJSON bool) {
	if rutaLog != "" {
		if err := GuardarTransaccionJSON(rutaLog, tx); err != nil {
			fmt.Printf("Advertencia: no se pudo guardar el log JSON: %v\n", err)
		}
	}
	if mostrarJSON {
		if jsonStr, err := FormatearTransaccionJSON(tx); err == nil {
			fmt.Println("Log JSON generado:")
			fmt.Println(jsonStr)
		}
	}
}

// TestEscenario4Integridad valida US-04: atomicidad y registro en formato JSON (EXITO, FALLO, ERROR).
func TestEscenario4Integridad(nombre string, saldoInicial float64, tipo string, monto float64, rutaLog string, mostrarJSON bool) {
	fmt.Printf("--- Escenario 4: Integridad de Transacción [%s] (%s) ---\n", tipo, nombre)
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)
	fmt.Printf("When: %s ejecuta %s atómico de %.2f\n", nombre, tipo, monto)

	tx, exito, _ := cuenta.EjecutarTransaccion(tipo, monto)
	registrarYMostrarLog(tx, rutaLog, mostrarJSON)

	switch tx.Estado {
	case EstadoExito:
		if exito && math.Abs(cuenta.ObtenerSaldo()-tx.SaldoNuevo) < 0.001 {
			fmt.Printf("Resultado: EXITO. Transacción atómica completada (%s). Nuevo saldo: %.2f\n\n", tx.ID, cuenta.ObtenerSaldo())
		} else {
			fmt.Printf("ERROR: Inconsistencia de saldo en transacción exitosa.\n\n")
		}
	case EstadoFallo:
		if math.Abs(cuenta.ObtenerSaldo()-saldoInicial) < 0.001 {
			fmt.Printf("Resultado: FALLO controlado (%s). Saldo intacto: %.2f\n\n", tx.Mensaje, cuenta.ObtenerSaldo())
		} else {
			fmt.Printf("ERROR: Violación de atomicidad; el saldo cambió tras un fallo.\n\n")
		}
	case EstadoError:
		if math.Abs(cuenta.ObtenerSaldo()-saldoInicial) < 0.001 {
			fmt.Printf("Resultado: ERROR controlado (%s). Saldo intacto: %.2f\n\n", tx.Mensaje, cuenta.ObtenerSaldo())
		} else {
			fmt.Printf("ERROR: Violación de atomicidad; el saldo cambió tras un error.\n\n")
		}
	}
}
