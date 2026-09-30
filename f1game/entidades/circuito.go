package entidades

import "math"

// PuntoCircuito es un punto del eje central del trazado en coordenadas de pantalla.
type PuntoCircuito struct {
	X, Y float64
}

// Circuito define la geometría y metadata de una pista.
type Circuito struct {
	ID              string
	Nombre          string
	PaisBandera     string
	VueltasCarrera  int
	DuracionVueltaS float64
	ArchivoTrazado  string

	AnchoPista float64
	Puntos     []PuntoCircuito
}

// TrazadoMonza es una versión jugable inspirada en el trazado de Monza:
// recta principal, Primera Variante, Curva Grande, Roggia, Lesmo 1/2,
// Ascari y la Parabólica. No es una reproducción a escala.
// Puntos de control del nuevo circuito. Se suavizan con una spline,
// así que alcanzan pocos puntos para tener curvas redondeadas.
var ControlNuevo = []PuntoCircuito{
	{800, 390}, {800, 290}, {790, 190}, {740, 110}, {650, 90},
	{570, 130}, {520, 200}, {440, 215}, {380, 160}, {300, 100},
	{200, 90}, {120, 140}, {105, 230}, {160, 300}, {250, 305},
	{300, 350}, {285, 400}, {340, 430}, {520, 438}, {690, 430},
}

// FactorCurva devuelve un multiplicador de ritmo aproximado para el tramo.
// Las rectas permiten más velocidad; las curvas reducen el ritmo.
func (c Circuito) FactorCurva(progreso float64) float64 {
	_, _, tx1, ty1 := c.PosicionEnProgreso(progreso)
	_, _, tx2, ty2 := c.PosicionEnProgreso(progreso + 0.018)
	dot := tx1*tx2 + ty1*ty2
	if dot > 1 {
		dot = 1
	}
	if dot < -1 {
		dot = -1
	}
	angulo := math.Acos(dot)
	// 0 rad = recta, valores altos = curva más cerrada.
	factor := 1.0 - (angulo/math.Pi)*0.48
	if factor < 0.52 {
		factor = 0.52
	}
	return factor
}

// PosicionEnProgreso convierte 0..1 en una posición sobre el centro del
// trazado y devuelve también el vector tangente.

// SuavizarCatmullRom convierte puntos de control en un trazado suave y cerrado.
func SuavizarCatmullRom(ctrl []PuntoCircuito, subdiv int) []PuntoCircuito {
	n := len(ctrl)
	out := make([]PuntoCircuito, 0, n*subdiv)
	for i := 0; i < n; i++ {
		p0 := ctrl[(i-1+n)%n]
		p1 := ctrl[i]
		p2 := ctrl[(i+1)%n]
		p3 := ctrl[(i+2)%n]
		for k := 0; k < subdiv; k++ {
			t := float64(k) / float64(subdiv)
			t2, t3 := t*t, t*t*t
			x := 0.5 * (2*p1.X + (-p0.X+p2.X)*t +
				(2*p0.X-5*p1.X+4*p2.X-p3.X)*t2 +
				(-p0.X+3*p1.X-3*p2.X+p3.X)*t3)
			y := 0.5 * (2*p1.Y + (-p0.Y+p2.Y)*t +
				(2*p0.Y-5*p1.Y+4*p2.Y-p3.Y)*t2 +
				(-p0.Y+3*p1.Y-3*p2.Y+p3.Y)*t3)
			out = append(out, PuntoCircuito{X: x, Y: y})
		}
	}
	return out
}

var TrazadoNuevo = SuavizarCatmullRom(ControlNuevo, 8)

var CircuitosDisponibles = []Circuito{
	{
		ID: "nuevo", Nombre: "Circuito Nuevo", PaisBandera: "ARG",
		VueltasCarrera: 3, DuracionVueltaS: 28,
		ArchivoTrazado: "procedural",
		AnchoPista:     64, Puntos: TrazadoNuevo,
	},
}

// PosicionEnProgreso convierte 0..1 en una posición sobre el centro del
// trazado y devuelve también el vector tangente. La geometría se recorre
// por distancia para que las curvas no cambien la velocidad visual.
func (c Circuito) PosicionEnProgreso(progreso float64) (x, y, tx, ty float64) {
	if len(c.Puntos) < 2 {
		return 0, 0, 1, 0
	}

	p := progreso - math.Floor(progreso)
	segLen := make([]float64, len(c.Puntos))
	total := 0.0
	for i := range c.Puntos {
		a := c.Puntos[i]
		b := c.Puntos[(i+1)%len(c.Puntos)]
		d := math.Hypot(b.X-a.X, b.Y-a.Y)
		segLen[i] = d
		total += d
	}
	if total == 0 {
		return c.Puntos[0].X, c.Puntos[0].Y, 1, 0
	}

	dist := p * total
	for i, d := range segLen {
		if dist <= d || i == len(segLen)-1 {
			a := c.Puntos[i]
			b := c.Puntos[(i+1)%len(c.Puntos)]
			t := dist / d
			if d == 0 {
				t = 0
			}
			tx, ty = b.X-a.X, b.Y-a.Y
			n := math.Hypot(tx, ty)
			if n > 0 {
				tx, ty = tx/n, ty/n
			}
			return a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t, tx, ty
		}
		dist -= d
	}
	return c.Puntos[0].X, c.Puntos[0].Y, 1, 0
}
