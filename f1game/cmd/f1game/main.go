package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"f1game/vistas"
)

func main() {
	ebiten.SetWindowSize(vistas.AnchoPantalla, vistas.AltoPantalla)
	ebiten.SetWindowTitle("F1 Game")

	juego := vistas.NuevoJuego()
	if err := ebiten.RunGame(juego); err != nil {
		log.Fatal(err)
	}
}
