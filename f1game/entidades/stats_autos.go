package entidades

// DatasetAutos contiene los 11 monoplazas/equipos de la parrilla 2026.
// Los autos tienen un rendimiento casi igual entre sí: la diferencia real
// en carrera la ponen los pilotos. Los puntos del constructor (clasificación
// al 29/09/2026) solo agregan una ventaja mínima.
var DatasetAutos = []Auto{
	construirAuto("mercedes_2026", "Mercedes", 2026, "W17", 538, 1, 538, 8.5),
	construirAuto("ferrari_2026", "Ferrari", 2026, "SF-26", 378, 2, 538, 8.5),
	construirAuto("mclaren_2026", "McLaren", 2026, "MCL40", 306, 3, 538, 8.5),
	construirAuto("red_bull_2026", "Red Bull Racing", 2026, "RB22", 263, 4, 538, 8.5),
	construirAuto("racing_bulls_2026", "Racing Bulls", 2026, "VCARB 03", 83, 5, 538, 8.5),
	construirAuto("alpine_2026", "Alpine", 2026, "A526", 68, 6, 538, 8.5),
	construirAuto("haas_2026", "Haas F1 Team", 2026, "VF-26", 27, 7, 538, 8.5),
	construirAuto("audi_2026", "Audi", 2026, "R26", 17, 8, 538, 8.5),
	construirAuto("williams_2026", "Williams", 2026, "FW48", 12, 9, 538, 8.5),
	construirAuto("aston_martin_2026", "Aston Martin", 2026, "AMR26", 3, 10, 538, 8.5),
	construirAuto("cadillac_2026", "Cadillac", 2026, "MAC-26", 0, 11, 538, 8.5),
}

func construirAuto(id, escuderia string, anio int, modelo string, puntos, posicion, puntosMax int, fiabilidad float64) Auto {
	r := ResultadoConstructor{
		PuntosConstructor:       puntos,
		VictoriasConstructor:    0,
		Abandonos:               0,
		CarrerasDisputadas:      15,
		PosicionFinalCampeonato: posicion,
		Fuente:                  "Formula1.com - 2026 Teams' Standings, consultado 29/09/2026",
	}
	velMax, acel, manejo, _ := NormalizarStatsAuto(r, puntosMax)
	return Auto{
		ID: id, Escuderia: escuderia, Anio: anio, Modelo: modelo,
		VelocidadMax: velMax, Aceleracion: acel, Manejo: manejo, Fiabilidad: fiabilidad,
		Resultado: r,
	}
}

// NormalizarStatsAuto deja a todos los autos entre 7.0 y 7.8: la diferencia
// entre el mejor y el peor equipo es de menos de 1 punto.
func NormalizarStatsAuto(r ResultadoConstructor, puntosMaxConstructorTemporada int) (velocidadMax, aceleracion, manejo, fiabilidad float64) {
	if puntosMaxConstructorTemporada <= 0 {
		puntosMaxConstructorTemporada = 1
	}
	rendimiento := float64(r.PuntosConstructor) / float64(puntosMaxConstructorTemporada)

	velocidadMax = clamp(7.0+rendimiento*0.80, 1, 10)
	aceleracion = clamp(7.0+rendimiento*0.76, 1, 10)
	manejo = clamp(7.0+rendimiento*0.72, 1, 10)
	fiabilidad = 8.5
	return
}
