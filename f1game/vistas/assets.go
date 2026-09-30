package vistas

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/sprites/*.png
var spritesFS embed.FS

var (
	imgAsfalto *ebiten.Image
	imgPasto   *ebiten.Image

	imgFondoMenu      *ebiten.Image
	imgFondoRuleta    *ebiten.Image
	imgFondoPausa     *ebiten.Image
	imgFondoResultado *ebiten.Image

	imgAutoJugador  *ebiten.Image // fallback, por si una escudería no matchea
	imgAutosRivales []*ebiten.Image
)

var spritesPorEscuderia map[string]*ebiten.Image

func init() {
	imgAsfalto = cargarSprite("asfalto.png")
	imgPasto = cargarSprite("pasto.png")

	imgFondoMenu = cargarSprite("fondo_menu.png")
	imgFondoRuleta = cargarSprite("fondo_ruleta.png")
	imgFondoPausa = cargarSprite("fondo_pausa.png")
	imgFondoResultado = cargarSprite("fondo_resultado.png")

	spritesPorEscuderia = map[string]*ebiten.Image{
		"Mercedes":     cargarSprite("auto_rival_gris.png"),
		"Ferrari":      cargarSprite("auto_jugador.png"),
		"McLaren":      cargarSprite("auto_rival_amarillo.png"),
		"Racing Bulls": cargarSprite("auto_rival_azul.png"),
		"Aston Martin": cargarSprite("auto_rival_verde.png"),
		"Cadillac":     cargarSprite("auto_rival_blanco.png"),
	}

	imgAutosRivales = make([]*ebiten.Image, 0, len(spritesPorEscuderia))
	for _, sprite := range spritesPorEscuderia {
		imgAutosRivales = append(imgAutosRivales, sprite)
	}

	// Sprite de referencia/fallback; durante la carrera se usa el
	// sprite correspondiente a la escudería real de cada corredor.
	imgAutoJugador = spritesPorEscuderia["Ferrari"]
}

func spriteParaEscuderia(escuderia string) *ebiten.Image {
	if sprite, ok := spritesPorEscuderia[strings.TrimSpace(escuderia)]; ok {
		return sprite
	}
	return imgAutoJugador
}

func dibujarFondoCompleto(screen *ebiten.Image, img *ebiten.Image) {
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	iw, ih := img.Bounds().Dx(), img.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(sw)/float64(iw), float64(sh)/float64(ih))
	screen.DrawImage(img, op)
}

func cargarSprite(nombre string) *ebiten.Image {
	datos, err := spritesFS.ReadFile("assets/sprites/" + nombre)
	if err != nil {
		log.Fatalf("no pude leer el sprite %s: %v", nombre, err)
	}
	img, _, err := image.Decode(bytes.NewReader(datos))
	if err != nil {
		log.Fatalf("no pude decodificar el sprite %s: %v", nombre, err)
	}
	return ebiten.NewImageFromImage(img)
}
