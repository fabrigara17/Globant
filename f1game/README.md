# F1 Game (Go + Ebiten v2)

## Estado actual

El proyecto es un juego de F1 arcade en Go + Ebiten.

### Cambios de esta versión

- **Vista cenital**: la carrera se muestra como un mapa visto desde arriba.
- **Circuito Monza**: trazado procedural inspirado en el Autodromo Nazionale Monza,
  con recta principal, chicanas, Lesmo, Ascari y Parabólica.
- **11 equipos de F1 2026** y **22 pilotos de la parrilla actual**.
- Los pilotos se asignan a su auto/equipo correspondiente; ya no se combinan
  pilotos y monoplazas de temporadas distintas.
- Sprites de autos cenitales diferenciados por escudería.
- La IA y el jugador reducen su ritmo en las curvas según la geometría del circuito.
- Las hitboxes ahora se calculan alrededor de la posición real del auto sobre el
  trazado, en lugar de usar una carretera pseudo-3D.
- Los valores de velocidad, aceleración, manejo y fiabilidad son **estadísticas de
  gameplay**, no ratings oficiales de F1/FIA.

## Parrilla 2026

Los 11 equipos y los 22 pilotos se basan en la parrilla oficial de Formula 1 para
2026. Los puntos usados para normalizar el rendimiento se actualizan con la
clasificación oficial consultada el 29/09/2026.

## Circuito

El trazado jugable está inspirado en Monza. El circuito real mide 5,793 m, tiene
7 curvas a derecha y 4 a izquierda y se recorre en sentido horario. El juego no
reproduce el circuito a escala: usa una versión simplificada y adaptada a una
pantalla de 960x540.

## Estadísticas

Los stats no deben interpretarse como mediciones oficiales. Se calculan para
crear diferencias de gameplay entre pilotos y equipos a partir de su rendimiento
actual.

## Cómo correr

```bash
go mod tidy
go run ./cmd/f1game
```

Para navegador se mantienen las opciones de WASM del proyecto original.

### Controles

- `Enter`: jugar
- `P`: elegir piloto como primera categoría de ruleta
- `A`: elegir auto como primera categoría
- `Espacio`: tirar la ruleta
- `W`: acelerar
- `S`: frenar
- `A / D`: moverse lateralmente
- `Esc`: pausa

## Fuentes de referencia

- Formula 1 — pilotos 2026.
- Formula 1 — equipos 2026.
- Formula 1 — clasificación de pilotos 2026.
- Formula 1 — clasificación de equipos 2026.
- Autodromo Nazionale Monza — información oficial del circuito.

El uso de nombres, marcas, libreas y otros elementos identificables de F1 puede
estar sujeto a derechos de sus respectivos titulares. Para publicación o
comercialización conviene revisar las condiciones aplicables.
