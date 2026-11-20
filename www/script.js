/**
 * prj_withdrawal — Lógica de Dominio Refactorizada (Web App)
 * Implementa: Retiro (US-01..03), Abono (US-01b) e Integridad/Log JSON (US-04)
 */

const STORAGE_KEY = "prj_withdrawal_state_v1";

// Respaldo incrustado equivalente a datos.json para ejecución directa vía file://
const DATOS_BASE_INCRUSTADOS = {
  proyecto: "prj_withdrawal",
  version: "1.3.0",
  moneda: "EUR",
  simbolo: "€",
  cuentas: [
    { id: "ACC-001", nombre: "Alice", saldo: 234.56, activa: true },
    { id: "ACC-002", nombre: "Hamza", saldo: 198.76, activa: true },
    { id: "ACC-003", nombre: "Bob", saldo: 100.00, activa: true },
    { id: "ACC-004", nombre: "Carlos", saldo: 500.00, activa: true },
    { id: "ACC-005", nombre: "Pedro", saldo: 300.00, activa: true }
  ],
  transacciones: []
};

const round2 = (n) => Math.round(n * 100) / 100;

class Cuenta {
  constructor({ id, nombre, saldo, activa = true }) {
    this.id = id;
    this.nombre = nombre;
    this.saldo = round2(Number(saldo));
    this.activa = activa;
  }

  /**
   * RN-01 / RN-02: No se permite retirar más del saldo ni montos negativos o cero.
   */
  retirar(cantidad) {
    if (Number.isNaN(cantidad) || cantidad <= 0) {
      throw new Error("No se pueden retirar cantidades negativas o cero");
    }
    if (cantidad > this.saldo) {
      return false; // Escenario: Declined withdrawal (FALLO)
    }
    this.saldo = round2(this.saldo - cantidad);
    return true;
  }

  /**
   * US-01b: El monto a abonar debe ser estrictamente mayor a 0.
   */
  abonar(cantidad) {
    if (Number.isNaN(cantidad) || cantidad <= 0) {
      throw new Error("El monto a abonar debe ser mayor a cero");
    }
    this.saldo = round2(this.saldo + cantidad);
    return true;
  }

  /**
   * US-04 / RN-03: Ejecución atómica con rollback y registro de auditoría.
   */
  ejecutarTransaccion(tipo, monto, seq) {
    const saldoAnterior = round2(this.saldo);
    const montoLimpio = Number.isNaN(monto) ? 0 : round2(monto);
    const tx = {
      id: `TXN-${Date.now()}-${String(seq).padStart(3, "0")}`,
      timestamp: new Date().toISOString(),
      cliente: this.nombre,
      tipo,
      monto: montoLimpio,
      saldo_anterior: saldoAnterior,
      saldo_nuevo: saldoAnterior,
      estado: "ERROR",
      mensaje: "",
      atomica: true
    };

    try {
      if (tipo === "RETIRO") {
        const ok = this.retirar(monto);
        if (!ok) {
          this.saldo = saldoAnterior; // Rollback explícito
          tx.estado = "FALLO";
          tx.mensaje = "Saldo insuficiente. Intente con un monto menor";
        } else {
          tx.saldo_nuevo = this.saldo;
          tx.estado = "EXITO";
          tx.mensaje = "Retiro procesado correctamente";
        }
      } else if (tipo === "ABONO") {
        this.abonar(monto);
        tx.saldo_nuevo = this.saldo;
        tx.estado = "EXITO";
        tx.mensaje = "Abono procesado correctamente";
      } else {
        throw new Error(`Tipo de operación no soportado: ${tipo}`);
      }
    } catch (err) {
      this.saldo = saldoAnterior; // Garantía de atomicidad
      tx.estado = "ERROR";
      tx.mensaje = err.message;
    }

    return tx;
  }
}

class GestorTransacciones {
  constructor() {
    this.proyecto = "prj_withdrawal";
    this.moneda = "EUR";
    this.simbolo = "€";
    this.cuentas = [];
    this.transacciones = [];
    this.fuente = "datos.json";
  }

  async inicializar(forzarReinicio = false) {
    if (!forzarReinicio) {
      const guardado = localStorage.getItem(STORAGE_KEY);
      if (guardado) {
        try {
          const parsed = JSON.parse(guardado);
          this.cargarDesdeObjeto(parsed, "localStorage (datos.json)");
          return;
        } catch {
          // Si localStorage está corrupto, recargar desde datos.json
        }
      }
    }

    try {
      const resp = await fetch("datos.json", { cache: "no-store" });
      if (!resp.ok) throw new Error("HTTP error");
      const data = await resp.json();
      this.cargarDesdeObjeto(data, "datos.json (fetch)");
    } catch {
      const copia = JSON.parse(JSON.stringify(DATOS_BASE_INCRUSTADOS));
      this.cargarDesdeObjeto(copia, "datos.json (embebido)");
    }
    this.persistir();
  }

  cargarDesdeObjeto(data, fuente) {
    this.proyecto = data.proyecto || "prj_withdrawal";
    this.moneda = data.moneda || "EUR";
    this.simbolo = data.simbolo || "€";
    this.cuentas = (data.cuentas || []).map((c) => new Cuenta(c));
    this.transacciones = Array.isArray(data.transacciones) ? data.transacciones : [];
    this.fuente = fuente;
  }

  persistir() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(this.toJSON()));
  }

  obtenerCuenta(id) {
    return this.cuentas.find((c) => c.id === id) || this.cuentas[0];
  }

  registrarOperacion(cuentaId, tipo, monto) {
    const cuenta = this.obtenerCuenta(cuentaId);
    const tx = cuenta.ejecutarTransaccion(tipo, monto, this.transacciones.length + 1);
    this.transacciones.unshift(tx);
    this.persistir();
    return { cuenta, tx };
  }

  toJSON() {
    return {
      proyecto: this.proyecto,
      moneda: this.moneda,
      simbolo: this.simbolo,
      total_eventos: this.transacciones.length,
      cuentas: this.cuentas.map((c) => ({
        id: c.id,
        nombre: c.nombre,
        saldo: c.saldo,
        activa: c.activa
      })),
      transacciones: this.transacciones
    };
  }
}

// --- Controlador de Interfaz ---
const gestor = new GestorTransacciones();
let operacionActual = "RETIRO";

const el = {
  selectCuenta: document.getElementById("select-cuenta"),
  saldoActual: document.getElementById("saldo-actual"),
  monedaSimbolo: document.getElementById("moneda-simbolo"),
  inputMonto: document.getElementById("input-monto"),
  txForm: document.getElementById("tx-form"),
  btnEjecutar: document.getElementById("btn-ejecutar"),
  feedback: document.getElementById("feedback-banner"),
  filterEstado: document.getElementById("filter-estado"),
  txTbody: document.getElementById("tx-tbody"),
  jsonOutput: document.getElementById("json-log-output"),
  statTotal: document.getElementById("stat-total"),
  statExito: document.getElementById("stat-exito"),
  statFallo: document.getElementById("stat-fallo"),
  statError: document.getElementById("stat-error"),
  sourceBadge: document.getElementById("data-source-badge"),
  btnReset: document.getElementById("btn-reset"),
  btnExport: document.getElementById("btn-export-json")
};

function renderCuentas(seleccionId) {
  const prev = seleccionId || el.selectCuenta.value;
  el.selectCuenta.innerHTML = gestor.cuentas
    .map((c) => `<option value="${c.id}">${c.nombre} (${c.id}) — ${c.saldo.toFixed(2)} ${gestor.simbolo}</option>`)
    .join("");
  if (prev && gestor.cuentas.some((c) => c.id === prev)) {
    el.selectCuenta.value = prev;
  }
  actualizarSaldoVisible();
}

function actualizarSaldoVisible() {
  const cuenta = gestor.obtenerCuenta(el.selectCuenta.value);
  if (cuenta) {
    el.saldoActual.textContent = `${cuenta.saldo.toFixed(2)} ${gestor.simbolo}`;
    el.saldoActual.className = `balance-value ${cuenta.saldo > 0 ? "amount-pos" : "amount-neg"}`;
  }
  el.monedaSimbolo.textContent = gestor.simbolo;
  el.sourceBadge.textContent = `Fuente: ${gestor.fuente}`;
}

function mostrarFeedback(tx) {
  const claseMonto = tx.estado === "EXITO" && tx.tipo === "ABONO" ? "amount-pos" : (tx.estado === "EXITO" ? "amount-pos" : "amount-neg");
  el.feedback.className = "feedback";
  el.feedback.innerHTML = `[${tx.estado}] ${tx.cliente} (${tx.tipo} <span class="${claseMonto}">${tx.monto.toFixed(2)} ${gestor.simbolo}</span>): ${tx.mensaje}. Saldo actual: <span class="${tx.saldo_nuevo > 0 ? "amount-pos" : "amount-neg"}">${tx.saldo_nuevo.toFixed(2)} ${gestor.simbolo}</span>`;
}

function renderLogs() {
  const filtro = el.filterEstado.value;
  const lista = filtro === "ALL"
    ? gestor.transacciones
    : gestor.transacciones.filter((t) => t.estado === filtro);

  el.statTotal.textContent = gestor.transacciones.length;
  el.statExito.textContent = gestor.transacciones.filter((t) => t.estado === "EXITO").length;
  el.statFallo.textContent = gestor.transacciones.filter((t) => t.estado === "FALLO").length;
  el.statError.textContent = gestor.transacciones.filter((t) => t.estado === "ERROR").length;

  if (lista.length === 0) {
    el.txTbody.innerHTML = `<tr><td colspan="7" class="empty-row">Sin transacciones para el filtro seleccionado.</td></tr>`;
  } else {
    el.txTbody.innerHTML = lista
      .map((t) => {
        const esIngreso = t.tipo === "ABONO" && t.estado === "EXITO";
        const claseMonto = esIngreso ? "amount-pos" : "amount-neg";
        const signo = t.estado === "EXITO" ? (t.tipo === "ABONO" ? "+" : "-") : "";
        return `
      <tr>
        <td><code>${t.id}</code></td>
        <td>${t.cliente}</td>
        <td>${t.tipo}</td>
        <td class="${claseMonto}">${signo}${t.monto.toFixed(2)} ${gestor.simbolo}</td>
        <td>${t.saldo_anterior.toFixed(2)} ${gestor.simbolo}</td>
        <td>${t.saldo_nuevo.toFixed(2)} ${gestor.simbolo}</td>
        <td><span class="badge">${t.estado}</span></td>
      </tr>`;
      })
      .join("");
  }

  el.jsonOutput.textContent = JSON.stringify(
    {
      proyecto: gestor.proyecto,
      total_eventos: gestor.transacciones.length,
      transacciones: lista
    },
    null,
    2
  );
}

function cambiarOperacion(op) {
  operacionActual = op;
  document.querySelectorAll(".op-tab").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.op === op);
  });
  el.btnEjecutar.textContent = op === "RETIRO" ? "Ejecutar Retiro" : "Ejecutar Abono";
}

function ejecutarYActualizar(cuentaId, tipo, monto) {
  const { cuenta, tx } = gestor.registrarOperacion(cuentaId, tipo, monto);
  renderCuentas(cuenta.id);
  mostrarFeedback(tx);
  renderLogs();
}

// Eventos
el.selectCuenta.addEventListener("change", actualizarSaldoVisible);
el.filterEstado.addEventListener("change", renderLogs);

document.querySelectorAll(".op-tab").forEach((tab) => {
  tab.addEventListener("click", () => cambiarOperacion(tab.dataset.op));
});

document.querySelectorAll(".chip").forEach((chip) => {
  chip.addEventListener("click", () => {
    el.inputMonto.value = chip.dataset.amount;
  });
});

el.txForm.addEventListener("submit", (e) => {
  e.preventDefault();
  const monto = parseFloat(el.inputMonto.value);
  ejecutarYActualizar(el.selectCuenta.value, operacionActual, monto);
});

const ESCENARIOS_BDD = {
  "alice-retiro": { nombre: "Alice", tipo: "RETIRO", monto: 200.0 },
  "hamza-retiro": { nombre: "Hamza", tipo: "RETIRO", monto: 200.0 },
  "pedro-error": { nombre: "Pedro", tipo: "RETIRO", monto: -20.0 },
  "bob-abono": { nombre: "Bob", tipo: "ABONO", monto: 50.0 },
  "bob-abono-fallo": { nombre: "Bob", tipo: "ABONO", monto: -20.0 }
};

document.querySelectorAll("[data-bdd]").forEach((btn) => {
  btn.addEventListener("click", () => {
    const esc = ESCENARIOS_BDD[btn.dataset.bdd];
    const cuenta = gestor.cuentas.find((c) => c.nombre === esc.nombre) || gestor.cuentas[0];
    cambiarOperacion(esc.tipo);
    el.inputMonto.value = esc.monto;
    ejecutarYActualizar(cuenta.id, esc.tipo, esc.monto);
  });
});

el.btnReset.addEventListener("click", async () => {
  localStorage.removeItem(STORAGE_KEY);
  await gestor.inicializar(true);
  el.feedback.className = "feedback hidden";
  renderCuentas();
  renderLogs();
});

el.btnExport.addEventListener("click", () => {
  const blob = new Blob([JSON.stringify(gestor.toJSON(), null, 2)], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "transacciones.json";
  a.click();
  URL.revokeObjectURL(url);
});

gestor.inicializar().then(() => {
  renderCuentas();
  renderLogs();
});
