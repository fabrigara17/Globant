package entidades

import "math/rand"

// Corredor es un participante de una carrera puntual: combina un Piloto
// y un Auto (ambos ya elegidos) con su estado dinámico durante la carrera.
type Corredor struct {
	Nombre    string
	Piloto    Piloto
	Auto      Auto
	EsJugador bool

	Progreso     float64
	VueltaActual int
	VelocidadKmh float64
	Posicion     int
	Terminado    bool
	OrdenLlegada int     // 0 = todavía no terminó; 1 = primero en llegar, etc.
	PosLateral   float64 // -1 (izquierda) a 1 (derecha), solo se usa para el jugador
}

// RitmoBase combina los stats normalizados de piloto+auto en un factor
// relativo de ritmo (1.0 ~ promedio de la parrilla). Cuanto más alto,
// más rápido corre esa combinación piloto+auto — así un Mercedes 2026
// con Verstappen 2026 corre claramente distinto a un Cadillac 2026
// con Hülkenberg.
func (c *Corredor) RitmoBase() float64 {
	statsPiloto := (c.Piloto.Aceleracion + c.Piloto.Manejo + c.Piloto.Reaccion + c.Piloto.Consistencia) / 4.0
	statsAuto := (c.Auto.VelocidadMax + c.Auto.Aceleracion + c.Auto.Manejo + c.Auto.Fiabilidad) / 4.0
	promedio := (statsPiloto + statsAuto) / 2.0 // rango ~1..10
	return promedio / 6.5                       // 6.5 ~ promedio del dataset; da factor ~0.6-1.5
}

// GenerarGrid arma los 10 corredores de una carrera: el jugador (con su
// piloto/auto elegidos en la ruleta) + 9 bots sacados al azar del
// dataset general, evitando repetir la misma combinación piloto+auto
// que ya está en la grilla.
func GenerarGrid(pilotoJugador Piloto, autoJugador Auto, rng *rand.Rand) []Corredor {
	grid := make([]Corredor, 0, 10)
	grid = append(grid, Corredor{
		Nombre: pilotoJugador.Nombre, Piloto: pilotoJugador, Auto: autoJugador,
		EsJugador: true, VueltaActual: 1,
	})

	// En 2026 cada piloto corre con el auto de su propia escudería.
	autosPorEscuderia := make(map[string]Auto, len(DatasetAutos))
	for _, a := range DatasetAutos {
		autosPorEscuderia[a.Escuderia] = a
	}

	usadosPilotos := map[string]bool{pilotoJugador.ID: true}
	intentos := 0
	for len(grid) < 10 && intentos < 500 {
		intentos++
		p := DatasetPilotos[rng.Intn(len(DatasetPilotos))]
		if usadosPilotos[p.ID] {
			continue
		}

		a, ok := autosPorEscuderia[p.Escuderia]
		if !ok {
			continue
		}

		usadosPilotos[p.ID] = true
		grid = append(grid, Corredor{
			Nombre: p.Nombre, Piloto: p, Auto: a,
			EsJugador: false, VueltaActual: 1,
		})
	}
	return grid
}
