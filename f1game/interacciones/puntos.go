package interacciones

// TablaPuntosF1 replica el sistema de puntuación oficial vigente desde
// 2010 (top 10 puntúan). Fuente: reglamento deportivo de la FIA.
// Índice 0 = 1er puesto.
var TablaPuntosF1 = map[int]int{
	1: 25, 2: 18, 3: 15, 4: 12, 5: 10,
	6: 8, 7: 6, 8: 4, 9: 2, 10: 1,
}

// TablaMonedaCarrera: moneda in-game otorgada por posición final, según
// el brief (10 desde el 1° hasta 1 en el último, en una carrera de 10
// pilotos). Es diseño de juego, no un dato de F1 real.
var TablaMonedaCarrera = map[int]int{
	1: 10, 2: 9, 3: 8, 4: 7, 5: 6,
	6: 5, 7: 4, 8: 3, 9: 2, 10: 1,
}

// PuntosF1PorPosicion devuelve 0 si la posición no puntúa (11°-20°).
func PuntosF1PorPosicion(posicion int) int {
	return TablaPuntosF1[posicion]
}

// MonedaPorPosicion devuelve 0 si la posición está fuera de rango (1-10).
func MonedaPorPosicion(posicion int) int {
	return TablaMonedaCarrera[posicion]
}

// ResultadoCarrera es lo que se acumula por piloto a lo largo de la
// temporada de 10 carreras.
type ResultadoCarrera struct {
	NombrePiloto string
	Posicion     int
	PuntosF1     int
	Moneda       int
}

// TablaCampeonato acumula puntos F1 de cada piloto (humano + bots) a lo
// largo de las 10 carreras. Al finalizar la última, el de más puntos es
// el campeón (con desempate simple por cantidad de victorias, como en
// el reglamento real).
type EntradaCampeonato struct {
	NombrePiloto  string
	PuntosTotales int
	Victorias     int
}

type TablaCampeonato struct {
	Entradas map[string]*EntradaCampeonato
}

func NuevaTablaCampeonato() *TablaCampeonato {
	return &TablaCampeonato{Entradas: make(map[string]*EntradaCampeonato)}
}

func (t *TablaCampeonato) RegistrarResultado(nombre string, posicion int) {
	e, ok := t.Entradas[nombre]
	if !ok {
		e = &EntradaCampeonato{NombrePiloto: nombre}
		t.Entradas[nombre] = e
	}
	e.PuntosTotales += PuntosF1PorPosicion(posicion)
	if posicion == 1 {
		e.Victorias++
	}
}

// Campeon devuelve el nombre del piloto líder tras la última carrera.
// Desempate: más victorias gana; si persiste el empate, se devuelve el
// primero encontrado (ampliar con más criterios FIA si hace falta:
// 2dos puestos, 3ros, etc. — dejado como TODO documentado).
func (t *TablaCampeonato) Campeon() string {
	var mejor *EntradaCampeonato
	for _, e := range t.Entradas {
		if mejor == nil ||
			e.PuntosTotales > mejor.PuntosTotales ||
			(e.PuntosTotales == mejor.PuntosTotales && e.Victorias > mejor.Victorias) {
			mejor = e
		}
	}
	if mejor == nil {
		return ""
	}
	return mejor.NombrePiloto
}
