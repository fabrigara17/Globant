package vistas

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"f1game/entidades"
	"f1game/interacciones"
)

type EstadoJuego int

const (
	EstadoMenuPrincipal EstadoJuego = iota
	EstadoRuletaEleccion
	EstadoRuletaGirando
	EstadoCarrera
	EstadoPausa
	EstadoResultadoCarrera
	EstadoResultadoTemporada
)

const (
	AnchoPantalla = 960
	AltoPantalla  = 540
)

// Paleta simple, look "F1 arcade" (cajas oscuras + acento rojo).
var (
	colorFondo    = color.RGBA{12, 12, 15, 255}
	colorCaja     = color.RGBA{28, 28, 33, 235}
	colorAcento   = color.RGBA{225, 6, 0, 255}
	colorTextoDeb = color.RGBA{235, 235, 235, 255}
)

// Paleta de pista.
var (
	colorPasto      = color.RGBA{20, 60, 25, 255}
	colorAsfaltoA   = color.RGBA{55, 55, 60, 255}
	colorAsfaltoB   = color.RGBA{65, 65, 70, 255}
	colorKerbRojo   = color.RGBA{200, 30, 30, 255}
	colorKerbBlanco = color.RGBA{230, 230, 230, 255}
)

const framesSpin = 40 // ~0.66s a 60 TPS: duración visual de cada tirada

type Hitbox struct {
	X     float64
	Y     float64
	Ancho float64
	Alto  float64
}

func (h Hitbox) Colisiona(otra Hitbox) bool {
	return h.X < otra.X+otra.Ancho &&
		h.X+h.Ancho > otra.X &&
		h.Y < otra.Y+otra.Alto &&
		h.Y+h.Alto > otra.Y
}

type Juego struct {
	Estado EstadoJuego

	Ruleta *interacciones.EstadoRuleta

	CarreraActualNro int
	Campeonato       *interacciones.TablaCampeonato
	MonedaJugador    int

	circuito        *entidades.Circuito
	carreraActiva   *interacciones.Carrera
	resultado       []entidades.Corredor
	tiempoCarrera   float64
	cuentaRegresiva float64 // segundos que faltan para largar
	mensajeYa       float64 // segundos que se muestra el "YA!" al largar

	// Cartel de fin de carrera (según la posición del jugador).
	cartelFinal  float64 // segundos que quedan de cartel
	cartelTitulo string
	cartelSub    string
	cartelColor  color.RGBA

	// Estado de la animación de la ruleta.
	girando       bool
	frameGiro     int
	nombreGiro    string
	esperandoAuto bool    // true si ya salió el piloto y falta tirar el auto (o viceversa)
	esperaInicio  float64 // segundos de pausa tras la segunda tirada, antes de la carrera
}

func NuevoJuego() *Juego {
	return &Juego{
		Estado:     EstadoMenuPrincipal,
		Campeonato: interacciones.NuevaTablaCampeonato(),
	}
}

func (j *Juego) Update() error {
	dt := 1.0 / 60.0

	switch j.Estado {
	case EstadoMenuPrincipal:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			j.Ruleta = interacciones.NuevaRuleta(time.Now().UnixNano())
			j.Estado = EstadoRuletaEleccion
		}

	case EstadoRuletaEleccion:
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			j.Ruleta.ElegirPrimeraCategoria(interacciones.CategoriaPiloto)
			j.Estado = EstadoRuletaGirando
		} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			j.Ruleta.ElegirPrimeraCategoria(interacciones.CategoriaAuto)
			j.Estado = EstadoRuletaGirando
		}

	case EstadoRuletaGirando:
		j.actualizarRuletaGirando()

	case EstadoCarrera:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			j.Estado = EstadoPausa
			return nil
		}
		// Cartel de fin de carrera: 2 segundos y después a los resultados.
		if j.cartelFinal > 0 {
			j.cartelFinal -= dt
			if j.cartelFinal <= 0 {
				j.cartelFinal = 0
				j.Estado = EstadoResultadoCarrera
			}
			return nil
		}

		// Cuenta regresiva: nadie se mueve hasta que llega a 0.
		if j.cuentaRegresiva > 0 {
			j.cuentaRegresiva -= dt
			if j.cuentaRegresiva <= 0 {
				j.cuentaRegresiva = 0
				j.mensajeYa = 1.0
			}
			return nil
		}
		if j.mensajeYa > 0 {
			j.mensajeYa -= dt
		}

		acelerando := ebiten.IsKeyPressed(ebiten.KeyW)
		frenando := ebiten.IsKeyPressed(ebiten.KeyS)
		girandoIzq := ebiten.IsKeyPressed(ebiten.KeyA)
		girandoDer := ebiten.IsKeyPressed(ebiten.KeyD)
		j.carreraActiva.Update(dt, acelerando, frenando, girandoIzq, girandoDer)
		j.tiempoCarrera += dt

		if jugador := j.carreraActiva.Jugador(); jugador != nil && j.detectarColisiones(jugador) {
			jugador.VelocidadKmh *= 0.8 // penalización simple al chocar: -20% de velocidad
		}

		if j.carreraActiva.Finalizada {
			j.resultado = j.carreraActiva.ResultadoFinal()
			j.registrarResultadoCarrera()
			if jugador := j.carreraActiva.Jugador(); jugador != nil {
				j.cartelTitulo, j.cartelSub, j.cartelColor = cartelPorPosicion(jugador.Posicion)
			}
			j.cartelFinal = 2.0 // el cartel dura 2 segundos
		}

	case EstadoPausa:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			j.Estado = EstadoCarrera
		}

	case EstadoResultadoCarrera:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			// TODO Fase 4: acá va el paso a la siguiente carrera de la
			// temporada (10 en total) en vez de volver directo al menú.
			j.Estado = EstadoMenuPrincipal
		}
	}
	return nil
}

// actualizarRuletaGirando maneja tanto el disparo de cada tirada (Espacio)
// como la animación de "giro" que dura framesSpin frames antes de fijar
// el resultado real (que sigue viniendo de interacciones.Tirar(), no se
// tocó esa lógica).
func (j *Juego) actualizarRuletaGirando() {
	// Pausa de 2 segundos después de la segunda tirada.
	if j.esperaInicio > 0 {
		j.esperaInicio -= 1.0 / 60.0
		if j.esperaInicio <= 0 {
			j.esperaInicio = 0
			j.iniciarCarrera()
		}
		return
	}

	if j.girando {
		j.frameGiro++
		// Cada 4 frames cambia el nombre mostrado, para dar sensación de giro.
		if j.frameGiro%4 == 0 {
			j.nombreGiro = j.nombreAlAzarParaAnimacion()
		}
		if j.frameGiro >= framesSpin {
			j.girando = false
			j.Ruleta.Tirar() // acá recién se fija el resultado real
			if j.Ruleta.Completa() {
				j.esperaInicio = 2.0 // en vez de arrancar de golpe
			}
		}
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && !j.Ruleta.Completa() {
		j.girando = true
		j.frameGiro = 0
		j.nombreGiro = j.nombreAlAzarParaAnimacion()
	}
}

// nombreAlAzarParaAnimacion elige un nombre random del dataset que
// corresponda (solo para el efecto visual de "está girando"; no influye
// en el resultado real, que sale de interacciones.EstadoRuleta).
func (j *Juego) posicionEnPista(c *entidades.Corredor) (x, y, angulo float64) {
	if j.circuito == nil {
		return AnchoPantalla / 2, AltoPantalla / 2, 0
	}

	x, y, tx, ty := j.circuito.PosicionEnProgreso(c.Progreso)

	// Desplazamiento lateral respecto del eje del circuito.
	normalX, normalY := -ty, tx
	lateral := c.PosLateral * (j.circuito.AnchoPista * 0.32)
	x += normalX * lateral
	y += normalY * lateral

	// Los sprites apuntan hacia arriba; se rotan para seguir la tangente.
	angulo = math.Atan2(ty, tx) + math.Pi/2
	return
}

func (j *Juego) hitboxJugador(jugador *entidades.Corredor) Hitbox {
	x, y, _ := j.posicionEnPista(jugador)
	return Hitbox{X: x - 30, Y: y - 42, Ancho: 60, Alto: 84}
}

func (j *Juego) hitboxOponente(_ *entidades.Corredor, rival *entidades.Corredor) (Hitbox, bool) {
	x, y, _ := j.posicionEnPista(rival)
	return Hitbox{X: x - 30, Y: y - 42, Ancho: 60, Alto: 84}, true
}

func (j *Juego) detectarColisiones(jugador *entidades.Corredor) bool {
	hitboxJugador := j.hitboxJugador(jugador)

	for i := range j.carreraActiva.Corredores {
		rival := &j.carreraActiva.Corredores[i]
		if rival.EsJugador {
			continue
		}

		hitboxRival, visible := j.hitboxOponente(jugador, rival)
		if visible && hitboxJugador.Colisiona(hitboxRival) {
			return true
		}
	}

	return false
}

func (j *Juego) nombreAlAzarParaAnimacion() string {
	fallaPiloto := j.Ruleta.PilotoAsignado == nil
	if fallaPiloto {
		p := entidades.DatasetPilotos[j.frameGiro%len(entidades.DatasetPilotos)]
		return fmt.Sprintf("%s (%d)", p.Nombre, p.Anio)
	}
	a := entidades.DatasetAutos[j.frameGiro%len(entidades.DatasetAutos)]
	return fmt.Sprintf("%s %s (%d)", a.Escuderia, a.Modelo, a.Anio)
}

func (j *Juego) iniciarCarrera() {
	c := entidades.CircuitosDisponibles[0]
	j.circuito = &c
	j.tiempoCarrera = 0
	j.cuentaRegresiva = 5.0
	j.mensajeYa = 0
	j.cartelFinal = 0
	j.carreraActiva = interacciones.NuevaCarrera(
		*j.circuito, *j.Ruleta.PilotoAsignado, *j.Ruleta.AutoAsignado,
		time.Now().UnixNano(),
	)
	j.Estado = EstadoCarrera
}

func (j *Juego) registrarResultadoCarrera() {
	for _, c := range j.resultado {
		j.Campeonato.RegistrarResultado(c.Nombre, c.Posicion)
		if c.EsJugador {
			j.MonedaJugador += interacciones.MonedaPorPosicion(c.Posicion)
		}
	}
}

// ---------- DIBUJO ----------

func caja(screen *ebiten.Image, x, y, w, h float32) {
	vector.DrawFilledRect(screen, x, y, w, h, colorCaja, false)
	vector.StrokeRect(screen, x, y, w, h, 2, colorAcento, false)
}

func (j *Juego) Draw(screen *ebiten.Image) {
	screen.Fill(colorFondo)

	switch j.Estado {
	case EstadoMenuPrincipal:
		j.dibujarMenu(screen)
	case EstadoRuletaEleccion:
		j.dibujarRuletaEleccion(screen)
	case EstadoRuletaGirando:
		j.dibujarRuletaGirando(screen)
	case EstadoCarrera:
		j.dibujarCarrera(screen)
	case EstadoPausa:
		j.dibujarPausa(screen)
	case EstadoResultadoCarrera:
		j.dibujarResultado(screen)
	}
}

func (j *Juego) dibujarMenu(screen *ebiten.Image) {
	dibujarFondoCompleto(screen, imgFondoMenu)
	caja(screen, 40, 40, 300, 60)
	ebitenutil.DebugPrintAt(screen, "F1 GAME", 55, 60)
	caja(screen, 40, 220, 260, 140)
	ebitenutil.DebugPrintAt(screen, "[Enter] Jugar", 55, 235)
	ebitenutil.DebugPrintAt(screen, "Pilotos", 55, 270)
	ebitenutil.DebugPrintAt(screen, "Monoplazas", 55, 295)
	ebitenutil.DebugPrintAt(screen, "Configuración", 55, 320)
}

func (j *Juego) dibujarRuletaEleccion(screen *ebiten.Image) {
	dibujarFondoCompleto(screen, imgFondoRuleta)
	caja(screen, 40, 40, 500, 50)
	ebitenutil.DebugPrintAt(screen, "Elegí qué tirar primero:", 55, 58)
	caja(screen, 40, 110, 230, 90)
	ebitenutil.DebugPrintAt(screen, "[P] Piloto", 60, 145)
	caja(screen, 300, 110, 230, 90)
	ebitenutil.DebugPrintAt(screen, "[A] Auto", 320, 145)
}

func (j *Juego) dibujarRuletaGirando(screen *ebiten.Image) {
	dibujarFondoCompleto(screen, imgFondoRuleta)
	// Caja Piloto
	caja(screen, 40, 40, 400, 180)
	ebitenutil.DebugPrintAt(screen, "PILOTO", 55, 55)
	switch {
	case j.girando && j.Ruleta.PilotoAsignado == nil:
		ebitenutil.DebugPrintAt(screen, "Girando: "+j.nombreGiro, 55, 100)
	case j.Ruleta.PilotoAsignado != nil:
		p := j.Ruleta.PilotoAsignado
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s", p.Nombre), 55, 100)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s - %d", p.Escuderia, p.Anio), 55, 125)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Aceleracion %.1f  Manejo %.1f", p.Aceleracion, p.Manejo), 55, 150)
	default:
		ebitenutil.DebugPrintAt(screen, "( sin tirar )", 55, 100)
	}

	// Caja Auto
	caja(screen, 480, 40, 400, 180)
	ebitenutil.DebugPrintAt(screen, "AUTO", 495, 55)
	switch {
	case j.girando && j.Ruleta.PilotoAsignado != nil && j.Ruleta.AutoAsignado == nil:
		ebitenutil.DebugPrintAt(screen, "Girando: "+j.nombreGiro, 495, 100)
	case j.Ruleta.AutoAsignado != nil:
		a := j.Ruleta.AutoAsignado
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s %s", a.Escuderia, a.Modelo), 495, 100)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", a.Anio), 495, 125)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("VelMax %.1f  Aceleracion %.1f", a.VelocidadMax, a.Aceleracion), 495, 150)
	default:
		ebitenutil.DebugPrintAt(screen, "( sin tirar )", 495, 100)
	}

	// Prompt inferior
	caja(screen, 40, 250, 840, 50)
	switch {
	case j.girando:
		ebitenutil.DebugPrintAt(screen, "Girando...", 55, 268)
	case j.Ruleta.Completa():
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("¡Listo! Iniciando carrera en %d...", int(math.Ceil(j.esperaInicio))), 55, 268)
	default:
		ebitenutil.DebugPrintAt(screen, "[Espacio] Tirar", 55, 268)
	}
}

func (j *Juego) dibujarCarrera(screen *ebiten.Image) {
	jugador := j.carreraActiva.Jugador()
	if jugador == nil {
		return
	}

	// Vista cenital: todo el trazado está visible como un mapa.
	// (dibujarPista ya pinta el césped de fondo — no hace falta
	// dibujarFondoCompleto acá, quedaría tapado igual.)
	j.dibujarPista(screen)
	j.dibujarOponentes(screen)
	j.dibujarAutoJugador(screen, jugador)

	// HUD compacto para no tapar el trazado.
	caja(screen, 18, 15, 185, 42)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s · VISTA CENITAL", j.circuito.Nombre), 30, 29)

	caja(screen, AnchoPantalla-205, 15, 185, 42)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("POS %d/10", jugador.Posicion), AnchoPantalla-190, 29)

	caja(screen, 18, AltoPantalla-70, 210, 50)
	ebitenutil.DebugPrintAt(screen, jugador.Piloto.Nombre, 30, AltoPantalla-58)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s · %s", jugador.Auto.Escuderia, jugador.Auto.Modelo), 30, AltoPantalla-40)

	caja(screen, AnchoPantalla-210, AltoPantalla-70, 190, 50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("VUELTA %d/%d", jugador.VueltaActual, j.circuito.VueltasCarrera), AnchoPantalla-195, AltoPantalla-58)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f km/h", jugador.VelocidadKmh), AnchoPantalla-195, AltoPantalla-40)

	// Cuenta regresiva y "YA!".
	if j.cuentaRegresiva > 0 {
		vector.DrawFilledRect(screen, 0, 0, AnchoPantalla, AltoPantalla, color.RGBA{0, 0, 0, 120}, false)
		n := int(math.Ceil(j.cuentaRegresiva))
		dibujarTextoGrande(screen, fmt.Sprintf("%d", n), AnchoPantalla/2, AltoPantalla/2-20, 14)
		dibujarTextoGrande(screen, "PREPARATE...", AnchoPantalla/2, AltoPantalla/2+70, 3)
	} else if j.mensajeYa > 0 {
		dibujarTextoGrande(screen, "YA!", AnchoPantalla/2, AltoPantalla/2-20, 14)
	}

	// Cartel de fin de carrera según la posición.
	if j.cartelFinal > 0 {
		j.dibujarCartelFinal(screen)
	}
}

// cartelPorPosicion devuelve el título, el subtítulo y el color del cartel
// que se muestra al terminar la carrera. Solo usa caracteres ASCII porque
// la fuente de debug de Ebiten no dibuja bien tildes ni signos como "¡".
func cartelPorPosicion(pos int) (titulo, sub string, col color.RGBA) {
	sub = fmt.Sprintf("Puesto %d de 10", pos)
	switch {
	case pos == 1:
		return "VICTORIA!", "Ganaste la carrera", color.RGBA{255, 200, 40, 255}
	case pos == 2:
		return "PODIO!", "Segundo puesto", color.RGBA{200, 205, 215, 255}
	case pos == 3:
		return "PODIO!", "Tercer puesto", color.RGBA{205, 127, 50, 255}
	case pos <= 5:
		return "CASI EN EL PODIO", sub, color.RGBA{80, 170, 255, 255}
	case pos <= 8:
		return "BUENA CARRERA", sub, color.RGBA{120, 220, 120, 255}
	default:
		return "SEGUI INTENTANDO", sub, color.RGBA{225, 6, 0, 255}
	}
}

func (j *Juego) dibujarCartelFinal(screen *ebiten.Image) {
	// Velo oscuro sobre la pista y caja con el color del cartel.
	vector.DrawFilledRect(screen, 0, 0, AnchoPantalla, AltoPantalla, color.RGBA{0, 0, 0, 140}, false)
	vector.DrawFilledRect(screen, 160, 170, 640, 200, colorCaja, false)
	vector.StrokeRect(screen, 160, 170, 640, 200, 4, j.cartelColor, false)

	// El título se achica solo si es muy largo para entrar en la caja.
	ancho := float64(len(j.cartelTitulo)*6 + 2)
	escala := 8.0
	if ancho*escala > 580 {
		escala = 580 / ancho
	}
	dibujarTextoGrandeColor(screen, j.cartelTitulo, AnchoPantalla/2, 245, escala, j.cartelColor)
	dibujarTextoGrandeColor(screen, j.cartelSub, AnchoPantalla/2, 330, 2, color.RGBA{235, 235, 235, 255})
}

// dibujarTextoGrande escribe texto blanco con la fuente de debug agrandada,
// centrado en (cx, cy).
func dibujarTextoGrande(screen *ebiten.Image, texto string, cx, cy, escala float64) {
	dibujarTextoGrandeColor(screen, texto, cx, cy, escala, color.White)
}

// dibujarTextoGrandeColor es igual que dibujarTextoGrande pero con color.
func dibujarTextoGrandeColor(screen *ebiten.Image, texto string, cx, cy, escala float64, col color.Color) {
	w, h := len(texto)*6+2, 16 // la fuente de debug usa ~6x16 px por caracter
	tmp := ebiten.NewImage(w, h)
	ebitenutil.DebugPrint(tmp, texto)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(escala, escala)
	op.GeoM.Translate(cx, cy)
	op.ColorScale.ScaleWithColor(col)
	screen.DrawImage(tmp, op)
	tmp.Dispose()
}

// normalEntre devuelve el vector perpendicular (de largo 1) a la dirección
// que va de ant a sig. Sirve para correr puntos hacia los bordes de la pista.
func normalEntre(ant, sig entidades.PuntoCircuito) entidades.PuntoCircuito {
	tx, ty := sig.X-ant.X, sig.Y-ant.Y
	n := math.Hypot(tx, ty)
	if n == 0 {
		return entidades.PuntoCircuito{}
	}
	return entidades.PuntoCircuito{X: -ty / n, Y: tx / n}
}

func dibujarLineaAncha(screen *ebiten.Image, puntos []entidades.PuntoCircuito, ancho float32, col color.Color) {
	for i := range puntos {
		a := puntos[i]
		b := puntos[(i+1)%len(puntos)]
		vector.StrokeLine(screen, float32(a.X), float32(a.Y), float32(b.X), float32(b.Y), ancho, col, true)
		vector.DrawFilledCircle(screen, float32(a.X), float32(a.Y), ancho/2, col, true)
	}
}

func (j *Juego) dibujarPista(screen *ebiten.Image) {
	// Base de césped.
	vector.DrawFilledRect(screen, 0, 0, AnchoPantalla, AltoPantalla, colorPasto, false)

	puntos := j.circuito.Puntos
	if len(puntos) < 2 {
		return
	}

	// Bordes de seguridad, piano y asfalto.
	dibujarLineaAncha(screen, puntos, float32(j.circuito.AnchoPista+20), colorKerbBlanco)
	dibujarLineaAncha(screen, puntos, float32(j.circuito.AnchoPista+8), colorAsfaltoB)
	dibujarLineaAncha(screen, puntos, float32(j.circuito.AnchoPista), colorAsfaltoA)

	// Pianos rojos/blancos en los dos bordes de la pista.
	n := len(puntos)
	distBorde := j.circuito.AnchoPista/2 + 2
	for i := 0; i < n; i++ {
		prev := puntos[(i-1+n)%n]
		a := puntos[i]
		b := puntos[(i+1)%n]
		next := puntos[(i+2)%n]

		na := normalEntre(prev, b) // normal en el punto a
		nb := normalEntre(a, next) // normal en el punto b

		col := colorKerbRojo
		if (i/2)%2 == 0 {
			col = colorKerbBlanco
		}
		for _, lado := range []float64{-1, 1} {
			d := distBorde * lado
			vector.StrokeLine(screen,
				float32(a.X+na.X*d), float32(a.Y+na.Y*d),
				float32(b.X+nb.X*d), float32(b.Y+nb.Y*d),
				7, col, true)
		}
	}

	// Línea discontinua muy sutil en el eje, solo como referencia visual.
	for i := 0; i < len(puntos); i += 2 {
		a := puntos[i]
		b := puntos[(i+1)%len(puntos)]
		vector.StrokeLine(screen, float32(a.X), float32(a.Y), float32(b.X), float32(b.Y), 1, color.RGBA{145, 145, 145, 130}, true)
	}

	// Meta.
	x, y, tx, ty := j.circuito.PosicionEnProgreso(0)
	nx, ny := -ty, tx
	vector.StrokeLine(screen,
		float32(x-nx*30), float32(y-ny*30),
		float32(x+nx*30), float32(y+ny*30),
		6, color.White, true)

	ebitenutil.DebugPrintAt(screen, "START / FINISH", 860, 388)
}

// dibujarAutoJugador dibuja el auto del jugador en su posición real
// sobre el trazado, rotado según la tangente de la pista.
const escalaAutos = 0.15 // ajustá este número: más chico = auto más chico

func (j *Juego) dibujarAutoJugador(screen *ebiten.Image, jugador *entidades.Corredor) {
	x, y, angulo := j.posicionEnPista(jugador)
	sprite := spriteParaEscuderia(jugador.Auto.Escuderia)

	w := float64(sprite.Bounds().Dx())
	h := float64(sprite.Bounds().Dy())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(escalaAutos, escalaAutos)
	op.GeoM.Rotate(angulo)
	op.GeoM.Translate(x, y)
	screen.DrawImage(sprite, op)
}

func (j *Juego) dibujarOponentes(screen *ebiten.Image) {
	for idx := range j.carreraActiva.Corredores {
		c := &j.carreraActiva.Corredores[idx]
		if c.EsJugador {
			continue
		}

		c.PosLateral = float64((idx%3)-1) * 0.55
		x, y, angulo := j.posicionEnPista(c)

		// Sprite asignado por índice del bot, no por escudería real —
		// así se garantiza variedad visual entre los 9 rivales aunque
		// el sorteo les haya tocado la misma escudería/auto.
		sprite := imgAutosRivales[idx%len(imgAutosRivales)]

		w := float64(sprite.Bounds().Dx())
		h := float64(sprite.Bounds().Dy())

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Scale(escalaAutos, escalaAutos)
		op.GeoM.Rotate(angulo)
		op.GeoM.Translate(x, y)
		screen.DrawImage(sprite, op)
	}
}

func dibujarHitbox(screen *ebiten.Image, h Hitbox) {
	vector.StrokeRect(
		screen,
		float32(h.X),
		float32(h.Y),
		float32(h.Ancho),
		float32(h.Alto),
		2,
		color.RGBA{255, 0, 0, 255},
		false,
	)
}

func (j *Juego) dibujarIndicadorTeclas(screen *ebiten.Image) {
	teclas := []struct {
		letra    string
		x        float32
		presiona func() bool
	}{
		{"A", 0, func() bool { return ebiten.IsKeyPressed(ebiten.KeyA) }},
		{"W", 30, func() bool { return ebiten.IsKeyPressed(ebiten.KeyW) }},
		{"S", 60, func() bool { return ebiten.IsKeyPressed(ebiten.KeyS) }},
		{"D", 90, func() bool { return ebiten.IsKeyPressed(ebiten.KeyD) }},
	}
	baseX := float32(AnchoPantalla/2 - 60)
	baseY := float32(AltoPantalla - 40)
	for _, t := range teclas {
		x := baseX + t.x
		if t.presiona() {
			vector.DrawFilledRect(screen, x, baseY, 26, 26, colorAcento, false)
		} else {
			vector.StrokeRect(screen, x, baseY, 26, 26, 2, colorAcento, false)
		}
		ebitenutil.DebugPrintAt(screen, t.letra, int(x)+9, int(baseY)+6)
	}
}

func formatearTiempo(segundos float64) string {
	minutos := int(segundos) / 60
	segs := int(segundos) % 60
	milis := int((segundos - float64(int(segundos))) * 1000)
	return fmt.Sprintf("%02d:%02d.%03d", minutos, segs, milis)
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func (j *Juego) dibujarPausa(screen *ebiten.Image) {
	dibujarFondoCompleto(screen, imgFondoPausa)
	caja(screen, AnchoPantalla/2-150, 180, 300, 160)
	ebitenutil.DebugPrintAt(screen, "EN PAUSA", AnchoPantalla/2-40, 210)
	ebitenutil.DebugPrintAt(screen, "[ESC] Continuar", AnchoPantalla/2-70, 250)
}

func (j *Juego) dibujarResultado(screen *ebiten.Image) {
	dibujarFondoCompleto(screen, imgFondoResultado)
	caja(screen, 40, 20, 500, 40)
	ebitenutil.DebugPrintAt(screen, "RESULTADO DE LA CARRERA", 55, 32)

	caja(screen, 40, 70, 500, 220)
	for i, c := range j.resultado {
		linea := fmt.Sprintf("%d. %-20s %s", c.Posicion, c.Nombre, c.Auto.Escuderia)
		ebitenutil.DebugPrintAt(screen, linea, 55, 90+i*20)
	}

	caja(screen, 40, AltoPantalla-90, 500, 60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Moneda acumulada: %d", j.MonedaJugador), 55, AltoPantalla-72)
	ebitenutil.DebugPrintAt(screen, "[Enter] Continuar", 55, AltoPantalla-48)
}

func (j *Juego) Layout(outsideWidth, outsideHeight int) (int, int) {
	return AnchoPantalla, AltoPantalla
}
