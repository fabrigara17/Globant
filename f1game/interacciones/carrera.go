package interacciones

import (
	"math/rand"
	"sort"

	"f1game/entidades"
)

// Carrera orquesta la simulación de una carrera puntual: avanza a
// jugador y bots cada frame, y mantiene el orden de posiciones
// actualizado.
type Carrera struct {
	Circuito   entidades.Circuito
	Corredores []entidades.Corredor
	rng        *rand.Rand
	Finalizada bool
	llegados   int // cuántos corredores ya cruzaron la meta (orden de llegada)
}

func NuevaCarrera(circuito entidades.Circuito, pilotoJugador entidades.Piloto, autoJugador entidades.Auto, seed int64) *Carrera {
	rng := rand.New(rand.NewSource(seed))
	return &Carrera{
		Circuito:   circuito,
		Corredores: entidades.GenerarGrid(pilotoJugador, autoJugador, rng),
		rng:        rng,
	}
}

// Update avanza la simulación dt segundos. acelerando/frenando son el
// input del jugador leído en vistas.
func (car *Carrera) Update(dt float64, acelerando, frenando, girandoIzq, girandoDer bool) {
	if car.Finalizada {
		return
	}
	for i := range car.Corredores {
		c := &car.Corredores[i]
		if c.EsJugador {
			ActualizarJugador(c, car.Circuito, dt, acelerando, frenando, girandoIzq, girandoDer)
		} else {
			ActualizarBot(c, car.Circuito, dt, car.rng)
		}

		// Registrar el orden de llegada la primera vez que termina.
		if c.Terminado && c.OrdenLlegada == 0 {
			car.llegados++
			c.OrdenLlegada = car.llegados
		}
	}
	car.recalcularPosiciones()

	if j := car.Jugador(); j != nil && j.Terminado {
		car.Finalizada = true
	}
}

func (car *Carrera) recalcularPosiciones() {
	orden := make([]int, len(car.Corredores))
	for i := range orden {
		orden[i] = i
	}
	progresoTotal := func(c *entidades.Corredor) float64 {
		return float64(c.VueltaActual) + c.Progreso
	}
	sort.SliceStable(orden, func(i, j int) bool {
		a := &car.Corredores[orden[i]]
		b := &car.Corredores[orden[j]]

		// Los que ya terminaron van primero, según el orden en que llegaron.
		if a.Terminado && b.Terminado {
			return a.OrdenLlegada < b.OrdenLlegada
		}
		if a.Terminado != b.Terminado {
			return a.Terminado
		}
		// Los que siguen en carrera se ordenan por cuánto recorrieron.
		return progresoTotal(a) > progresoTotal(b)
	})
	for pos, idx := range orden {
		car.Corredores[idx].Posicion = pos + 1
	}
}

func (car *Carrera) Jugador() *entidades.Corredor {
	for i := range car.Corredores {
		if car.Corredores[i].EsJugador {
			return &car.Corredores[i]
		}
	}
	return nil
}

// ResultadoFinal devuelve a los corredores ordenados por posición final
// (para la pantalla de resultados y para alimentar TablaCampeonato).
func (car *Carrera) ResultadoFinal() []entidades.Corredor {
	res := make([]entidades.Corredor, len(car.Corredores))
	copy(res, car.Corredores)
	sort.Slice(res, func(i, j int) bool { return res[i].Posicion < res[j].Posicion })
	return res
}
