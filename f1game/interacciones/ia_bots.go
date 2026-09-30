package interacciones

import (
	"math/rand"

	"f1game/entidades"
)

// ActualizarBot avanza el progreso de un corredor manejado por IA. Usa
// su ritmo base (derivado de piloto+auto) + jitter aleatorio frame a
// frame para simular variaciones de ritmo real, y una pequeña chance de
// "bache" ligada a la fiabilidad del auto (autos poco fiables sufren
// microcaídas de ritmo más seguido — no es un abandono real, eso queda
// para una fase futura).
func ActualizarBot(c *entidades.Corredor, circuito entidades.Circuito, dt float64, rng *rand.Rand) {
	if c.Terminado {
		return
	}

	ritmo := c.RitmoBase()

	// Jitter: +/-8% de variación aleatoria por frame.
	jitter := 1.0 + (rng.Float64()*0.16 - 0.08)

	// Bache de fiabilidad: cuanto más baja la fiabilidad del auto, más
	// frecuente la microcaída de ritmo.
	factorFiabilidad := 1.0
	probabilidadBache := (10.0 - c.Auto.Fiabilidad) / 10.0 * 0.02 // hasta 2% de chance por frame
	if rng.Float64() < probabilidadBache {
		factorFiabilidad = 0.6
	}

	factorCurva := circuito.FactorCurva(c.Progreso)
	avanceVuelta := (dt / circuito.DuracionVueltaS) * ritmo * jitter * factorFiabilidad * factorCurva
	c.Progreso += avanceVuelta
	c.VelocidadKmh = 120 + ritmo*75*jitter*factorCurva // HUD; no representa una velocidad real de F1

	for c.Progreso >= 1.0 {
		c.Progreso -= 1.0
		c.VueltaActual++
	}
	if c.VueltaActual > circuito.VueltasCarrera {
		c.Terminado = true
		c.Progreso = 1.0
		c.VueltaActual = circuito.VueltasCarrera
	}
}

// ActualizarJugador aplica el input de acelerar/frenar sobre el ritmo
// base del jugador. El techo y el piso de velocidad quedan definidos
// por el ritmo base de su propio piloto+auto.
func ActualizarJugador(c *entidades.Corredor, circuito entidades.Circuito, dt float64, acelerando, frenando, girandoIzq, girandoDer bool) {
	if c.Terminado {
		return
	}

	ritmo := c.RitmoBase()

	// Sin W el auto queda quieto (tampoco se mueve de costado).
	if !acelerando {
		c.VelocidadKmh = 0
		return
	}

	factorInput := 1.15
	if frenando {
		factorInput = 0.5 // W + S: avanza más lento
	}

	factorCurva := circuito.FactorCurva(c.Progreso)
	avanceVuelta := (dt / circuito.DuracionVueltaS) * ritmo * factorInput * factorCurva
	c.Progreso += avanceVuelta
	c.VelocidadKmh = 120 + ritmo*75*factorInput*factorCurva

	// Movimiento lateral: A/D mueven el auto entre -1 (borde izquierdo)
	// y 1 (borde derecho) de la pista, con una velocidad de giro fija.
	const velocidadGiro = 2.5 // unidades por segundo
	switch {
	case girandoIzq && !girandoDer:
		c.PosLateral -= velocidadGiro * dt
	case girandoDer && !girandoIzq:
		c.PosLateral += velocidadGiro * dt
	}
	if c.PosLateral < -1 {
		c.PosLateral = -1
	}
	if c.PosLateral > 1 {
		c.PosLateral = 1
	}

	for c.Progreso >= 1.0 {
		c.Progreso -= 1.0
		c.VueltaActual++
	}
	if c.VueltaActual > circuito.VueltasCarrera {
		c.Terminado = true
		c.Progreso = 1.0
		c.VueltaActual = circuito.VueltasCarrera
	}
}
