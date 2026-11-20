package headers

import (
	"errors"
	"fmt"
)

// Cuenta representa la entidad de negocio descrita en los Features Gherkin.
type Cuenta struct {
	saldo  float64
	nombre string
}

// NewCuenta es el constructor de Cuenta.
func NewCuenta(nombre string, saldoInicial float64) *Cuenta {
	return &Cuenta{
		nombre: nombre,
		saldo:  saldoInicial,
	}
}

// ObtenerSaldo devuelve el saldo actual.
func (c *Cuenta) ObtenerSaldo() float64 {
	return c.saldo
}

// Nombre devuelve el nombre del titular.
func (c *Cuenta) Nombre() string {
	return c.nombre
}

/*
Retirar intenta retirar una cantidad.
Regla: No se puede retirar más de lo que hay en la cuenta ni valores negativos.
*/
func (c *Cuenta) Retirar(cantidad float64) (bool, error) {
	if cantidad < 0 {
		return false, errors.New("no se pueden retirar cantidades negativas")
	}
	if cantidad > c.saldo {
		return false, nil // Escenario: Withdrawal declined
	}
	c.saldo -= cantidad
	return true, nil // Escenario: Successful withdrawal
}

/*
Abonar intenta ingresar dinero a la cuenta.
Regla: El monto debe ser estrictamente mayor a 0.
*/
func (c *Cuenta) Abonar(cantidad float64) (bool, error) {
	if cantidad <= 0 {
		return false, errors.New("el monto a abonar debe ser mayor a cero")
	}
	c.saldo += cantidad
	return true, nil // Escenario: Successful deposit
}

// String implementa fmt.Stringer para depuración.
func (c Cuenta) String() string {
	return fmt.Sprintf("Cuenta {%s, saldo=%.2f}", c.nombre, c.saldo)
}
