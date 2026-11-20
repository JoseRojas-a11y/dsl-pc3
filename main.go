package main

import (
	"flag"
	"fmt"
	"os"
	"prj_withdrawal/headers"
)

func main() {
	nombre := flag.String("nombre", "", "Nombre del titular de la cuenta")
	saldo := flag.Float64("saldo", 0.0, "Saldo inicial de la cuenta")
	retiro := flag.Float64("retiro", 0.0, "Cantidad a retirar")
	abono := flag.Float64("abono", 0.0, "Cantidad a abonar")
	operacion := flag.String("op", "retiro", "Tipo de operación: 'retiro', 'abono' o 'integridad'")
	escenario := flag.String("test", "", "Tipo de test forzado: 'exito', 'fallo', 'error'")
	logFile := flag.String("log", "transacciones.json", "Ruta del archivo JSON para registrar transacciones (US-04)")
	mostrarJSON := flag.Bool("json", false, "Imprimir el registro JSON de la transacción en consola")
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("=== EJECUTANDO PRUEBAS BASE DE GHERKIN ===")
		fmt.Println("[Feature: Withdrawing cash]")
		headers.TestEscenario1Exito("Alice", 234.56, 200.00)
		headers.TestEscenario2Fallo("Hamza", 198.76, 200.00)
		fmt.Println("[Feature: Depositing cash]")
		headers.TestEscenarioAbonoExito("Bob", 100.00, 50.00)
		headers.TestEscenarioAbonoFallo("Bob", 100.00, -20.00)
		fmt.Println("[Feature: Transaction integrity & JSON logging (US-04)]")
		headers.TestEscenario4Integridad("Alice", 234.56, headers.TipoRetiro, 200.00, *logFile, true)
		headers.TestEscenario4Integridad("Hamza", 198.76, headers.TipoRetiro, 200.00, *logFile, true)
		headers.TestEscenario4Integridad("Pedro", 300.00, headers.TipoRetiro, -20.00, *logFile, true)
		headers.TestEscenario4Integridad("Bob", 100.00, headers.TipoAbono, 50.00, *logFile, true)
		os.Exit(0)
	}

	if *nombre == "" {
		fmt.Println("Error: Debes proporcionar al menos un nombre. Usa -h para ayuda.")
		os.Exit(1)
	}

	if *operacion == "integridad" {
		fmt.Println("=== PROCESANDO TRANSACCIÓN ATÓMICA CON LOG JSON (US-04) ===")
		if *abono != 0 {
			headers.TestEscenario4Integridad(*nombre, *saldo, headers.TipoAbono, *abono, *logFile, *mostrarJSON)
		} else {
			headers.TestEscenario4Integridad(*nombre, *saldo, headers.TipoRetiro, *retiro, *logFile, *mostrarJSON)
		}
		os.Exit(0)
	}

	if *operacion == "abono" {
		fmt.Println("=== PROCESANDO OPERACIÓN DE ABONO ===")
		if *mostrarJSON || *logFile != "" {
			headers.TestEscenario4Integridad(*nombre, *saldo, headers.TipoAbono, *abono, *logFile, *mostrarJSON)
		} else if *escenario == "fallo" || *abono <= 0 {
			headers.TestEscenarioAbonoFallo(*nombre, *saldo, *abono)
		} else {
			headers.TestEscenarioAbonoExito(*nombre, *saldo, *abono)
		}
		os.Exit(0)
	}

	// Operación de retiro por defecto
	fmt.Println("=== PROCESANDO OPERACIÓN DE RETIRO ===")
	if *mostrarJSON {
		headers.TestEscenario4Integridad(*nombre, *saldo, headers.TipoRetiro, *retiro, *logFile, true)
		os.Exit(0)
	}

	switch *escenario {
	case "exito":
		headers.TestEscenario1Exito(*nombre, *saldo, *retiro)
	case "fallo":
		headers.TestEscenario2Fallo(*nombre, *saldo, *retiro)
	case "error":
		headers.TestEscenario3Error(*nombre, *saldo, *retiro)
	default:
		if *retiro < 0 {
			headers.TestEscenario3Error(*nombre, *saldo, *retiro)
		} else if *retiro > *saldo {
			headers.TestEscenario2Fallo(*nombre, *saldo, *retiro)
		} else {
			headers.TestEscenario1Exito(*nombre, *saldo, *retiro)
		}
	}
	os.Exit(0)
}
