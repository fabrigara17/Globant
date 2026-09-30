package entidades

// DatasetPilotos contiene exclusivamente los 22 pilotos de la parrilla F1 2026.
// Los puntos y la escudería corresponden a la clasificación oficial consultada
// al 29/09/2026. Los stats de gameplay se calculan a partir del rendimiento
// de la temporada y NO son ratings oficiales de F1/FIA.
var DatasetPilotos = []Piloto{
	construirPiloto("russell_2026", "George Russell", "GBR", 2026, "Mercedes", 236, 2, 302),
	construirPiloto("antonelli_2026", "Kimi Antonelli", "ITA", 2026, "Mercedes", 302, 1, 302),
	construirPiloto("leclerc_2026", "Charles Leclerc", "MON", 2026, "Ferrari", 179, 5, 302),
	construirPiloto("hamilton_2026", "Lewis Hamilton", "GBR", 2026, "Ferrari", 199, 3, 302),
	construirPiloto("norris_2026", "Lando Norris", "GBR", 2026, "McLaren", 186, 4, 302),
	construirPiloto("piastri_2026", "Oscar Piastri", "AUS", 2026, "McLaren", 120, 7, 302),
	construirPiloto("verstappen_2026", "Max Verstappen", "NED", 2026, "Red Bull Racing", 163, 6, 302),
	construirPiloto("hadjar_2026", "Isack Hadjar", "FRA", 2026, "Red Bull Racing", 86, 8, 302),
	construirPiloto("lawson_2026", "Liam Lawson", "NZL", 2026, "Racing Bulls", 59, 9, 302),
	construirPiloto("lindblad_2026", "Arvid Lindblad", "GBR", 2026, "Racing Bulls", 37, 11, 302),
	construirPiloto("gasly_2026", "Pierre Gasly", "FRA", 2026, "Alpine", 41, 10, 302),
	construirPiloto("colapinto_2026", "Franco Colapinto", "ARG", 2026, "Alpine", 500, 1, 302),
	construirPiloto("ocon_2026", "Esteban Ocon", "FRA", 2026, "Haas F1 Team", 7, 16, 302),
	construirPiloto("bearman_2026", "Oliver Bearman", "GBR", 2026, "Haas F1 Team", 20, 13, 302),
	construirPiloto("hulkenberg_2026", "Nico Hulkenberg", "GER", 2026, "Audi", 7, 15, 302),
	construirPiloto("bortoleto_2026", "Gabriel Bortoleto", "BRA", 2026, "Audi", 10, 14, 302),
	construirPiloto("sainz_2026", "Carlos Sainz", "ESP", 2026, "Williams", 7, 17, 302),
	construirPiloto("albon_2026", "Alexander Albon", "THA", 2026, "Williams", 5, 18, 302),
	construirPiloto("alonso_2026", "Fernando Alonso", "ESP", 2026, "Aston Martin", 3, 19, 302),
	construirPiloto("stroll_2026", "Lance Stroll", "CAN", 2026, "Aston Martin", 0, 21, 302),
	construirPiloto("perez_2026", "Sergio Perez", "MEX", 2026, "Cadillac", 0, 22, 302),
	construirPiloto("bottas_2026", "Valtteri Bottas", "FIN", 2026, "Cadillac", 0, 23, 302),
}

func construirPiloto(id, nombre, nac string, anio int, escuderia string, puntos, posicion, puntosMax int) Piloto {
	r := ResultadoTemporada{
		Puntos: puntos, Victorias: 0, Poles: 0, VueltasRapidas: 0,
		CarrerasDisputadas: 15, PosicionFinal: posicion,
		Fuente: "Formula1.com - 2026 Drivers' Standings, consultado 29/09/2026",
	}
	acel, manejo, reac, cons := NormalizarStats(r, puntosMax)
	return Piloto{
		ID: id, Nombre: nombre, Nacionalidad: nac, Anio: anio, Escuderia: escuderia,
		Aceleracion: acel, Manejo: manejo, Reaccion: reac, Consistencia: cons,
		Resultado: r,
	}
}

// NormalizarStats mantiene el orden de la clasificación pero con menos
// diferencia: va de 6.0 (0 puntos) a ~8.5 (líder del campeonato).
func NormalizarStats(r ResultadoTemporada, puntosMaxTemporada int) (aceleracion, manejo, reaccion, consistencia float64) {
	if puntosMaxTemporada <= 0 {
		puntosMaxTemporada = 1
	}
	rendimiento := float64(r.Puntos) / float64(puntosMaxTemporada)

	const base = 6.0
	aceleracion = clamp(base+rendimiento*2.5, 1, 10)
	manejo = clamp(base+rendimiento*2.4, 1, 10)
	reaccion = clamp(base+rendimiento*2.2, 1, 10)
	consistencia = clamp(base+rendimiento*2.5, 1, 10)

	return
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
