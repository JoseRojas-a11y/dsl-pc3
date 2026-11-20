# UNIVERSIDAD NACIONAL DE INGENIERÍA
**Facultad de Ingeniería Industrial y de Sistemas**  
**Escuela Profesional de Ingeniería de Software**

---

# Especificación: DSL / Gherkin-BDD
**Análisis, Diseño e Implementación**  
*Lenguaje - Metodología - Herramientas*

## DSL: Especificación Gherkin
Gherkin, BDD y Cucumber son componentes interrelacionados dentro de la metodología de **Desarrollo Guiado por Comportamiento (BDD)** que buscan alinear el negocio con la tecnología mediante pruebas automatizadas legibles.

### BDD - Gherkin: Desarrollo dirigido por comportamiento
BDD es una evolución de TDD (*Test Driven Development* o Desarrollo Dirigido por Pruebas). De hecho, el concepto de BDD fue inicialmente introducido por Dan North como respuesta a los problemas que surgían al enseñar TDD.

TDD está basado en 2 prácticas: **Escribir las pruebas primero** y **Refactorizar después**:

* **Técnica TDD:**
  1. En primer lugar, se escribe una prueba y se verifica que las pruebas fallan.
  2. A continuación, se implementa el código que hace que la prueba pase satisfactoriamente.
  3. Seguidamente se refactoriza el código escrito.

En BDD también vamos a escribir las pruebas antes de escribir el código fuente, pero en lugar de pruebas unitarias, lo que haremos será *"Escribir pruebas que verifiquen que el comportamiento del código es correcto desde el punto de vista de negocio"*. Tras escribir las pruebas, escribimos el código fuente de la funcionalidad que haga que estas pruebas pasen correctamente. Después refactorizamos el código fuente.

### Técnica BDD
Partiremos de historias de usuario, siguiendo el modelo:
```text
As a: [rol stakeholder]
I want: [a feature]
so that: [need / problem to solve].

```

A partir de aquí, en lugar de describir en "lenguaje natural" lo que tiene que hacer esa nueva funcionalidad, vamos a usar un lenguaje que nos va a permitir describir todas nuestras funcionalidades de una misma forma: un lenguaje específico para BDD, como **Gherkin**.

#### Ejemplo de User Story

> **Plantilla:**
> * **As a** `<role>`
> * **I want** `<goal>`
> * **so that** `<benefit>`
> * **Acceptance criteria** *(Conditions of Satisfaction)*
> 
> 

> **Caso práctico:**
> * **As an** Account Manager
> * **I want** a sales report of my account to be sent to my inbox daily
> * **So that** I can monitor the sales progress of my customer portfolio
> * **Acceptance criteria:**
> 1. The report is sent daily to my inbox
> 2. The report contains the following sales details: ...
> 3. The report is in csv format.
> 
> 
> 
> 

La **Metodología BDD** es un proceso de desarrollo de software que trata de combinar los aspectos puramente técnicos y los de negocio, de manera que tengamos un marco de trabajo y un marco de pruebas en el que los requisitos de negocio formen parte del proceso de desarrollo.

---

## codeLAB

### 1. Análisis

Lo normal es probar distintos escenarios para comprobar una determinada funcionalidad. De esta forma vamos a pasar de nuestras historias de usuario a pruebas de comportamiento automatizables. Para automatizar estas pruebas necesitamos apoyarnos en herramientas, como **Cucumber**, y otras que entienden perfectamente el lenguaje Gherkin.

**Archivo Gherkin (`retiro-cash.feature`):**

```gherkin
Feature: Retiro de efectivo
  Regla: Los clientes no pueden retirar más de su saldo

Feature: Withdrawing cash
  Rule: Los clientes no pueden retirar más de su saldo

  Scenario: Retiro exitoso dentro del saldo.
    Given Alice tiene 234.56 en su cuenta.
    When Alice intenta retirar 200.00
    Then el retiro se realiza con éxito.

  Scenario: Retiro rechazado por exceder el saldo.
    Given Hamza tiene 198.76 en su cuenta.
    When Hamza intenta retirar 200.00
    Then el retiro se rechaza.

```

#### Caso de Uso: Realizar Retiro de Efectivo (versión 1.2.0)

* **ID:** UC-001
* **Nombre:** Realizar Retiro de Efectivo
* **Actor Principal:** Cliente (Titular de la cuenta)
* **Actores Secundarios:** Sistema Bancario (Core), Cajero Automático (ATM) / App Móvil (Canal)
* **Nivel:** Objetivo de Usuario

##### 1. Descripción

Este caso de uso permite al Cliente solicitar una cantidad de efectivo de su cuenta bancaria. El sistema valida si la cuenta tiene fondos suficientes y, en caso afirmativo, procesa la transacción descontando el monto. Si los fondos son insuficientes, el sistema rechaza la solicitud y notifica al cliente sin realizar ningún movimiento financiero.

##### 2. Precondiciones

* El Cliente está autenticado en el sistema (ATM, Banca Online o App).
* El Cliente tiene una cuenta activa con saldo disponible.
* El dispositivo (ATM) tiene efectivo disponible o la conexión a la red está operativa.

#### User Stories y Reglas de Negocio (OCL) — Caso de Uso "Retirar efectivo"

##### 1. User Stories

* **US-01: Retiro exitoso** *(escenario: Successful withdrawal within balance)*
* **As a:** cliente autenticado
* **I want:** retirar efectivo, indicando un monto
* **so that:** disponer de mis ahorros y ver mi nuevo saldo
* **Acceptance criteria:** Dado que estoy autenticado, mi cuenta está activa y el ATM tiene efectivo, cuando ingreso un monto válido menor o igual a mi saldo, entonces el sistema descuenta el monto, registra la transacción como exitosa, dispensa el efectivo y muestra el nuevo saldo.


* **US-02: Retiro rechazado por saldo insuficiente** *(escenario: Declined withdrawal in excess of balance)*
* **As a:** cliente en sesión
* **I want:** recibir un mensaje claro cuando mi saldo no alcanza
* **so that:** elegir un monto menor o cancelar sin perder dinero.
* **Acceptance criteria:** Cuando el monto solicitado es mayor que mi saldo, el sistema muestra *"Saldo insuficiente. Intente con un monto menor"*, no descuenta nada y me permite ingresar otro monto o cancelar.


* **US-03: Validación del monto**
* **As a:** cliente en sesión
* **I want:** que el sistema valide el monto ingresado
* **so that:** corregir errores antes de operar.
* **Acceptance criteria:** Si el monto es negativo, cero, no es múltiplo de los billetes disponibles o excede el límite diario, el sistema muestra el error específico y vuelve a solicitar el monto.


* **US-04: Integridad de la transacción**
* **As a:** banco (servicio ATM)
* **I want:** que cada retiro sea atómico y quede registrado
* **so that:** garantizar consistencia contable y auditoría.
* **Acceptance criteria:** Una transacción se completa por entero o no se hace. Todo intento deja rastro en el historial (los fallidos según la política de auditoría).



#### Narrativas de Casos de Uso

##### 3. Flujo Principal (Escenario de Éxito)

*(Corresponde al escenario: "Successful withdrawal within balance")*

1. El Cliente selecciona la opción `"Retirar efectivo"` en el menú.
2. El Sistema solicita ingresar el monto a retirar.
3. El Cliente ingresa el monto (ej. `200.00`) y confirma.
4. El Sistema verifica que el monto sea válido (mayor a 0 y múltiplo de billetes disponibles).
5. El Sistema consulta el saldo actual de la cuenta del Cliente.
6. El Sistema comprueba que el saldo sea mayor o igual al monto solicitado.
7. El Sistema descuenta el monto de la cuenta.
8. El Sistema registra la transacción como exitosa en el historial.
9. El Sistema dispensa el efectivo (si es ATM) o muestra el mensaje de éxito (si es digital).
10. El Sistema muestra el nuevo saldo disponible.
11. El Caso de Uso finaliza con éxito.

##### 4. Flujos Alternativos (Escenario de Error)

* **4.a. Saldo Insuficiente** *(Corresponde al escenario: "Declined withdrawal in excess of balance")*
En el paso 6 del flujo principal:
* El Sistema detecta que el saldo disponible es menor que el monto solicitado.
* El Sistema cancela la operación y **NO** descuenta ningún dinero.
* El Sistema muestra un mensaje de error: `"Saldo insuficiente. Intente con un monto menor"`.
* El Sistema ofrece la opción de ingresar un nuevo monto o cancelar.
* Si el Cliente elige cancelar, el Caso de Uso termina. Si elige otro monto, vuelve al paso 3.


* **4.b. Monto Inválido o Límites Excedidos**
En el paso 4 del flujo principal:
* El Sistema detecta que el monto es negativo, cero, no es múltiplo de billetes o excede el límite diario.
* El Sistema muestra un mensaje de validación indicando el error específico.
* El Sistema solicita al Cliente que ingrese un nuevo monto.
* Vuelve al paso 3.



##### 5. Postcondiciones

* **Si fue exitoso:** El saldo de la cuenta se ha reducido en la cantidad solicitada y se ha generado un registro de transacción.
* **Si fue rechazado:** El saldo de la cuenta permanece inalterado y se ha registrado un intento fallido (opcional, según política de auditoría).

##### 6. Reglas de Negocio Asociadas

* **RN-01:** No se permite realizar un retiro si el monto solicitado es mayor al saldo disponible.
* **RN-02:** No se permiten saldos negativos por retiros de efectivo (sin sobregiro autorizado).
* **RN-03:** Las transacciones deben ser atómicas (o se completa todo o no se hace nada).

---

### 2. Diseño (DSL / Gherkin-BDD-Notation)

**Archivo:** `características/retiro-cash.feature`

```gherkin
Feature: Withdrawing cash
  Rule: Customers cannot withdraw more than their balance

  Scenario: Successful withdrawal within balance
    Given Alice has 234.56 in their account
    When Alice tries to withdraw 200.00
    Then the withdrawal is successful

  Scenario: Declined withdrawal in excess of balance
    Given Hamza has 198.76 in their account
    When Hamza tries to withdraw 200.00
    Then the withdrawal is declined

```

#### Modelo de Dominio

```text
[Cliente] 1..1 ──── 1..* [Cuenta] 1..1 ──── 1..* [Transaccion]
                            │
                         (canal)
                            │
                         [Cajero]

```

#### Modelo de contexto (base para el OCL)

```ocl
Context: Cuenta
  self.saldo: real
  self.activa: logico
  self.limiteDiario: entero
  self.retiradoHoy: fecha

Context: Cliente
  self.autenticado: logico

Context: Cajero
  self.efectivoDisponible: real
  self.operativo: logico
  self.denominacionMinima: entero

Context: Transaccion
  self.tipo: {RETIRO}
  self.monto: real
  self.estado: {EXITOSA, RECHAZADA}

```

> **Nota:** Los montos utilizan `cast(entero, monto (unidad mínima))` para evitar errores de redondeo.

---

### 3. Implementación: Go v.1.28

**Prueba de validación e Implementación:** Valida e implementa en lenguaje Go el Feature Gherkin: `"Feature: Withdrawing cash"`.

**Observación:** Versión 1.2.0 del proyecto, adaptada para aceptar argumentos por línea de comandos (CLI) y con validaciones avanzadas para probar múltiples combinaciones (retiros correctos, insuficiencia de fondos, montos negativos, etc.).

Se mantiene la estructura similar al proyecto original (`prj_withdrawal`), pero se han refactorizado las funciones de test para que sean dinámicas y se ha actualizado `main.go` para usar el paquete `flag` de Go (que es más robusto y profesional que leer `os.Args` manualmente).

#### Referencia: Estructura y ejecución en C++ (Base previa)

```text
prj_gherkin/
└── withdraw/
    ├── retiro-cash.feature
    └── app_retiro-cash.cpp

```

**Terminal (asumiendo G++ instalado):**

* **Paso 1: Compilar**
```bash
g++ app-retiro-cash.cpp -o app-retiro-cash

```


* **Paso 2: Ejecutar**
```bash
./app-retiro-cash

```



#### a) Estructura del proyecto en Go

```text
prj_withdraw/
├── main.go
├── go.mod
└── headers/
    ├── class.go      // Entidad Cuenta
    └── test_bdd.go   // Validación de escenarios Gherkin

```

#### b) Lenguaje Go

##### 1. `headers/class.go` (Entidad Cuenta)

Se mantiene igual, ya que la lógica del negocio no cambia, solo la forma en que la invocamos.

```go
/*
source file: class.go
Cuenta: Representa la entidad de negocio descrita en el Feature Gherkin.
Date: nov-2026
Author: J.Medianero. A (adaptado a Go)
*/
package headers

import (
	"errors"
	"fmt"
)

// Cuenta representa la entidad de negocio descrita en el Feature Gherkin.
type Cuenta struct {
	saldo  float64
	nombre string
}

// NewCuenta es el constructor equivalente a Cuenta(nombre, saldoInicial).
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

// Nombre devuelve el nombre del titular (útil para logs).
func (c *Cuenta) Nombre() string {
	return c.nombre
}

/*
Retirar intenta retirar una cantidad.
Regla: No se puede retirar más de lo que hay en la cuenta.
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

// String implementa fmt.Stringer para depuración.
func (c *Cuenta) String() string {
	return fmt.Sprintf("Cuenta{%s, saldo \$=%.2f}", c.nombre, c.saldo)
}

```

##### 2. `headers/test_bdd.go` (Validación Dinámica)

Aquí está el cambio clave: las funciones ahora reciben parámetros para que puedan ser probadas con cualquier dato.

```go
/*
Source file: test_bdd.go
Representa la entidad de validación de escenarios descrita en el Feature Gherkin.
Date: nov-2026
Author: J.Medianero. A (adaptado a Go)
*/
package headers

import (
	"fmt"
	"math"
)

// imprimirSaldo imprime el saldo formateado a 2 decimales.
func imprimirSaldo(cuenta *Cuenta) {
	fmt.Printf("Saldo de %.2f €\n", cuenta.ObtenerSaldo())
}

// -----------------------------------------------------------------
// Validación de Escenarios Gherkin Dinámicos
// -----------------------------------------------------------------

// TestEscenario1Exito: Simula un retiro exitoso con datos dinámicos
func TestEscenario1Exito(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 1: Successful withdrawal within balance (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	// Then
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
			fmt.Println("ERROR: El saldo restante no coincide.\n")
		}
	} else {
		fmt.Println("Resultado: FALLO. La retirada fue rechazada (incorrecto).\n")
	}
}

// TestEscenario2Fallo: Simula un retiro rechazado con datos dinámicos
func TestEscenario2Fallo(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 2: Declined withdrawal in excess of balance (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	// Then
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
			fmt.Println("ERROR: El saldo se modificó incorrectamente.\n")
		}
	} else {
		fmt.Println("Resultado: FALLO. La retirada fue procesada (debería haber sido rechazada).\n")
	}
}

// TestEscenario3Error: Simula un error de validación (ej. cantidad negativa)
func TestEscenario3Error(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 3: Validation Error (%s) ---\n", nombre)
	cuenta := NewCuenta(nombre, saldoInicial)
	exito, err := cuenta.Retirar(cantidadRetirar)

	if err != nil {
		fmt.Printf("Resultado: ÉXITO. El sistema atrapó el error correctamente: %v\n", err)
		imprimirSaldo(cuenta)
		fmt.Println()
	} else if !exito {
		fmt.Println("Resultado: RECHAZADO. Fondos insuficientes (no es un error de validación).\n")
	} else {
		fmt.Println("Resultado: FALLO. Se procesó un retiro que debía dar error.\n")
	}
}

```

##### 3. `main.go` (Ejecutable CLI)

Este es el archivo principal modificado. Usa el paquete `flag` para leer los parámetros.

```go
/*
Source file: main.go [aplicación main CLI v.1.2.0]
Date: nov 2026
Author: J.Medianero.A
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"prj_withdrawal/headers"
)

func main() {
	// Definir banderas (flags) para los argumentos de línea de comandos
	nombre := flag.String("nombre", "", "Nombre del titular de la cuenta")
	saldo := flag.Float64("saldo", 0.0, "Saldo inicial de la cuenta")
	retiro := flag.Float64("retiro", 0.0, "Cantidad a retirar")

	// Bandera especial para forzar el escenario de test
	escenario := flag.String("test", "", "Tipo de test a ejecutar: 'exito', 'fallo', 'error'")

	flag.Parse()

	// Si no se pasan argumentos, ejecutamos los escenarios por defecto del examen
	if len(os.Args) == 1 {
		fmt.Println("=== EJECUTANDO PRUEBAS POR DEFECTO (Gherkin v1.0.0) ===")
		headers.TestEscenario1Exito("Alice", 234.56, 200.00)
		headers.TestEscenario2Fallo("Hamza", 198.76, 200.00)
		os.Exit(0)
	}

	// Validar que se hayan ingresado los datos necesarios
	if *nombre == "" {
		fmt.Println("Error: Debes proporcionar al menos un nombre. Usa -h para ayuda.")
		os.Exit(1)
	}

	// Enrutamiento de escenarios
	switch *escenario {
	case "exito":
		headers.TestEscenario1Exito(*nombre, *saldo, *retiro)
	case "fallo":
		headers.TestEscenario2Fallo(*nombre, *saldo, *retiro)
	case "error":
		headers.TestEscenario3Error(*nombre, *saldo, *retiro)
	default:
		// Si no se especifica escenario, deducimos qué pasará según los datos
		fmt.Println("=== EJECUTANDO TEST DINÁMICO AUTOMÁTICO ===")
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

```

---

### Configurar Proyecto y Ejecutar (v.1.2.0)

En tu terminal, sigue estos pasos (directorio de trabajo `prj_withdrawal`):

#### Configurar proyecto: Compilación y ejecución

```bash
mkdir prj_withdrawal && cd prj_withdrawal
mkdir headers
# (Copia y pega los archivos class.go, test_bdd.go y main.go en sus carpetas)

# Inicializar el módulo Go
go mod init prj_withdrawal

# Compilar el ejecutable
go build -o prj_retiro.exe main.go

# Alternativa de ejecución en Go:
# go.mod necesario para compilar
# go mod init prj_retiro
# go run main.go

```

#### Pruebas CLI Go

Ejemplos de ejecución en línea de comandos (`CLI go`) usando parámetros `flag`:

1. **Ejecución por defecto (los escenarios del examen original):**
```bash
go run main.go

```


2. **Prueba dinámica exitosa (Carlos retira 150 de 500):**
*(El sistema detectará automáticamente que es un éxito)*
```bash
go run main.go -nombre "Carlos" -saldo 500 -retiro 150

```


3. **Prueba forzada de Fondos Insuficientes (María quiere retirar 100 teniendo 50):**
```bash
go run main.go -nombre "María" -saldo 50 -retiro 100 -test fallo

```


4. **Prueba forzada de Error (Pedro intenta retirar un monto negativo, simulando "sin clave" o error de sistema):**
```bash
go run main.go -nombre "Pedro" -saldo 300 -retiro -20 -test error

```


5. **Ver la ayuda de la herramienta CLI:**
```bash
go run main.go -h

```



#### ✔ Mejoras incluidas en la v.2.1.0:

* **Uso de `flag`:** Es el estándar idiomático de Go para CLI (*Command Line Interface*). Permite usar `-nombre="Juan" -saldo=100.5`.
* **Modo Automático:** Si no usas la bandera `-test`, el código evalúa los datos matemáticamente. Si pides retirar más de lo que tienes, automáticamente ejecuta el Escenario 2 (*Declined*). Si pides negativo, ejecuta el Escenario 3 (*Error*).
* **Cero recompilación:** Puedes simular miles de cajeros automáticos distintos con usuarios y saldos diferentes ejecutando el mismo binario.

---

*Autor: Jorge Medianero Acosta — Ing. de Sistemas y Cómputo (Págs. 1 - 10)*

```

*(Nota: En los bloques de código Go y comandos de terminal se restauraron operadores y caracteres como `:=`, `!=`, `-=`, `-o` y `-h` que aparecían cortados por la extracción del PDF original para que el código sea completamente funcional).*[cite: 2, 3, 4, 5, 6, 7, 8, 9, 10, 11]