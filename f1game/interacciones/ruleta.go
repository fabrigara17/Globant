package interacciones

import (
	"errors"
	"math/rand"

	"f1game/entidades"
)

// Categoria representa qué se está sorteando en cada tirada.
type Categoria int

const (
	CategoriaPiloto Categoria = iota
	CategoriaAuto
)

// EstadoRuleta modela el flujo de 2 tiradas descrito en el brief:
//  1. El jugador ELIGE qué categoría tirar primero (piloto o auto).
//  2. Se sortea un resultado random de esa lista.
//  3. La segunda tirada es automática y sortea la categoría restante.
type EstadoRuleta struct {
	primeraCategoriaElegida bool
	primeraCategoria        Categoria
	PilotoAsignado          *entidades.Piloto
	AutoAsignado            *entidades.Auto
	rng                     *rand.Rand
}

func NuevaRuleta(seed int64) *EstadoRuleta {
	return &EstadoRuleta{rng: rand.New(rand.NewSource(seed))}
}

// ElegirPrimeraCategoria: acción del jugador al arrancar. Debe llamarse
// una sola vez, antes de la primera tirada.
func (e *EstadoRuleta) ElegirPrimeraCategoria(c Categoria) error {
	if e.primeraCategoriaElegida {
		return errors.New("la primera categoría ya fue elegida")
	}
	e.primeraCategoriaElegida = true
	e.primeraCategoria = c
	return nil
}

// Tirar ejecuta una tirada de ruleta. La primera vez sortea sobre la
// categoría elegida por ElegirPrimeraCategoria; la segunda vez sortea
// automáticamente sobre la categoría restante. Devuelve error si ya se
// completaron ambas tiradas o si no se eligió categoría inicial.
func (e *EstadoRuleta) Tirar() error {
	if !e.primeraCategoriaElegida {
		return errors.New("primero hay que elegir piloto o auto con ElegirPrimeraCategoria")
	}
	if e.PilotoAsignado != nil && e.AutoAsignado != nil {
		return errors.New("ya se completaron las 2 tiradas")
	}

	// Determinar qué categoría toca en esta tirada.
	var categoriaActual Categoria
	switch {
	case e.PilotoAsignado == nil && e.AutoAsignado == nil:
		categoriaActual = e.primeraCategoria
	case e.PilotoAsignado == nil:
		categoriaActual = CategoriaPiloto
	default:
		categoriaActual = CategoriaAuto
	}

	if categoriaActual == CategoriaPiloto {
		idx := e.rng.Intn(len(entidades.DatasetPilotos))
		p := entidades.DatasetPilotos[idx]
		e.PilotoAsignado = &p
	} else {
		idx := e.rng.Intn(len(entidades.DatasetAutos))
		a := entidades.DatasetAutos[idx]
		e.AutoAsignado = &a
	}
	return nil
}

// Completa indica si ya se sortearon piloto y auto.
func (e *EstadoRuleta) Completa() bool {
	return e.PilotoAsignado != nil && e.AutoAsignado != nil
}
