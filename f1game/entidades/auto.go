package entidades

// Auto representa el monoplaza de una escudería en una temporada dada.
// Igual que Piloto: los stats de gameplay se derivan del rendimiento
// oficial del CONSTRUCTOR (no del piloto), para separar "qué tan bueno
// era el auto" de "qué tan bueno era el piloto que lo manejaba".
type Auto struct {
	ID        string // ej: "red_bull_2026"
	Escuderia string
	Anio      int
	Modelo    string // ej: "RB19"

	VelocidadMax float64 // 0-10
	Aceleracion  float64 // 0-10
	Manejo       float64 // 0-10 (downforce/aero)
	Fiabilidad   float64 // 0-10 (abandonos por fallas mecánicas)

	Resultado ResultadoConstructor
}

type ResultadoConstructor struct {
	PuntosConstructor       int
	VictoriasConstructor    int
	Abandonos               int // por falla mecánica, no por choque
	CarrerasDisputadas      int
	PosicionFinalCampeonato int
	Fuente                  string
}
