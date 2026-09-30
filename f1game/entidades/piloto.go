package entidades

// Piloto representa a un piloto de F1 en una temporada específica.
// Los campos de stats (0-10) NO son datos oficiales de la FIA/F1 (esa
// info no existe publicada): se DERIVAN de resultados oficiales de esa
// temporada mediante NormalizarStats(). Así cualquier valor es auditable
// contra la fuente en ResultadoTemporada.
type Piloto struct {
	ID           string // ej: "verstappen_2026"
	Nombre       string
	Nacionalidad string
	Anio         int
	Escuderia    string

	// Stats de gameplay (0.0 - 10.0), calculados por NormalizarStats.
	Aceleracion  float64
	Manejo       float64 // consistencia en curva / control
	Reaccion     float64 // salidas de largada (derivado de posiciones ganadas en vuelta 1, aproximado)
	Consistencia float64 // derivado de terminar carreras / puntuar seguido

	Resultado ResultadoTemporada
}

// ResultadoTemporada son los datos oficiales/verificables de la temporada
// (fuente: FIA / Formula1.com / archivos históricos de resultados).
type ResultadoTemporada struct {
	Puntos             int
	Victorias          int
	Poles              int
	VueltasRapidas     int
	CarrerasDisputadas int
	PosicionFinal      int    // posición en el campeonato de pilotos
	Fuente             string // URL o referencia de donde salió el dato
}
