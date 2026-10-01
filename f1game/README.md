Juego de Fórmula 1 
https://github.com/fabrigara17/Globant.git

Integrantes del grupo: Fabricio Garavaglia, Luciano Christiansen, Emmanuel Christiansen, Micaela Origlio y Alejo Rodriguez

INSTRUCCIONES PARA DESCARGAR EL JUEGO:

1. Instalar Go manualmente desde 
https://go.dev/dl/ (paso manual, no hay comando para esto)

2. Cerrar y reabrir la terminal 
go version 

3. Clonar y entrar al proyecto 
cd C:\Users\fabri\Proyectos (este es mi caso, solo hay q seleccionar una carpeta)
git clone https://github.com/fabrigara17/Globant.git
cd Globant\f1game 

4. Dependencias y ejecución 
go mod tidy 
go run ./cmd/f1game 


Para este proyecto se usó una cuenta principal de claude y a medida que nos íbamos quedando sin tokens rotamos entre otra cuenta de claude, chat gpt y gemini aquí están los prompts en orden de cuenta no en orden de promt

CLAUDE 1
Spec detallada del juego: menú con ruleta (piloto + auto de listas separadas, pilotos de F1 desde 1990, autos de cada escudería desde 1990), stats diferenciados por piloto/auto basados en rendimiento real, carrera con 9 bots + jugador (10 total), 3 vueltas de ~25-30s cada una, sistema de moneda por posición + puntos oficiales de F1, temporada de 10 carreras con campeón final, circuitos con forma real.
Ampliar el dataset de pilotos y autos primero, y preparar el juego para correr en localhost desde el navegador, antes de seguir con los bots.
(Subiendo un zip) "Eso es lo último que tengo del juego" — cargarle más pilotos y autos con estadísticas, y preparlo para correr en localhost desde el navegador.
(Subiendo otro zip) Ayudar a armar la estructura faltante del juego, pero pidiendo que le pase los códigos como texto para pegarlos él mismo (sin crear los archivos yo), y el comando para correrlo desde el cmd.
Pegó el error de PowerShell: go no reconocido como comando.
Repitió el mismo error de PowerShell.
Aclaración: "VS Code me tira eso, pero el PowerShell me tira la versión" (o sea, funcionaba en una terminal sí y en otra no).
Pegó error de go mod tidy: versión inválida de gomobile.
Pegó error tras go clean -modcache / go get: rutas de módulo mal formadas (por pegar varios comandos juntos).
(Imagen) El juego se ve con pantalla negra y HUD; preguntó por qué al iniciar solo tira una ruleta en vez de dos, y pidió corregir las vistas y el funcionamiento de la ruleta.
(Imagen) "El juego va mejor, ahora hay que arreglar la vista de la carrera, que se sigue viendo así" (sin fondo de pista).
(Imagen) "Se sigue viendo así" (mismo problema).
Respondió que la terminal no tiró ningún error de compilación.
(6 imágenes: captura del juego + 5 de un video de referencia) "Ahora se ve así" — pidió agregar movimiento con las flechas A/D, y mejorar las vistas para que se parezcan a las imágenes de referencia.
(Pegó el código de juego.go) "Decime si ahí quedó bien."
Pidió que le tire todo el código de juego.go corregido.
(Imagen) "El juego se sigue viendo así, ¿hay alguna forma de asemejarlo a las imágenes que te pasé?"
Dijo que le gustaba la Opción 1 (sprites en Ebiten), pero preguntó qué tan complicado sería migrar a un motor 3D (Opción 3).
Contó que el video de referencia salió de ChatGPT, y pidió seguir con la Opción 1, dejando la Opción 3 como posible intento si quedaba feo.
Pegó error de compilación: patrón de //go:embed inválido (../assets/sprites/*.png).
(Subiendo zip) "Le cambié un par de cosas, quiero seguir con el fondo para los menús y eso."
(Imagen) "Me tira ese error en la línea 305 de juego."
Aclaró que el error aparece como subrayado rojo en el editor.
(Imagen con el mensaje de error completo) "¿Ahí se llega a ver bien?"
(Imagen) Error al guardar: "Failed to save... the content of the file is newer."
Preguntó si el conflicto de guardado fue porque descomprimió los archivos de fondo dentro de esa carpeta.
(Imagen del diff de "resolve save conflict") "Eso me tira."
Pegó el archivo juego.go completo (versión con vista cenital de Monza), sin pregunta adicional, en respuesta a mi pedido de sincronizar.
Pegó otra versión de juego.go (la vieja, con bugs de dibujarFondoCompleto mal puesto) + "Tengo eso ahora, no sé qué toqué, pasame el código que te había pasado antes con las correcciones ya hechas."
Avisó que tira errores "undefined" en las líneas de spriteParaEscuderia.
Pegó el error de terminal completo con los undefined: dibujarFondoCompleto / imgFondoMenu, etc.
Confirmó que assets.go sí tenía la función dibujarFondoCompleto en el editor.
"Me tira el error abajo cuando intento guardar" (sin adjuntar nada la primera vez).
(Imagen) Mismo error de guardado, con el código visible al lado — "Ese error, perdón."
Avisó que ya lo arregló, y pidió achicar los sprites de los autos y cambiar el circuito a uno ovalado tipo Indianápolis 500.
Pidió que los autos rivales tengan colores diferentes entre sí, no todos rojos.
Aclaró: no quiere más sprites nuevos, quiere que en el código se asigne un sprite distinto a cada bot.
Retomó el pedido del circuito ovalado tipo Indianápolis 500.
Eligió: reemplazar Monza por el óvalo (en vez de agregarlo como circuito nuevo).
Pegó el contenido completo de entidades/circuito.go.
(Este mensaje) Pidió esta lista de todos los prompts del chat.



GEMINI 
Lista de Prompts
Prompt 1: explicame este problema de mi proyecto de go, las imagen q supuestamente no puede econtrar esta en la carpeta q dice (junto con el comando de ejecución inicial).
Prompt 2: y ahora explicame este otro errorgo: cannot find main module, but found .git/config in C:\Users\fabri\Proyectos\Globant...
Prompt 3: me vuelve a tirar este error Run 'go help mod init' for more information...
Prompt 4: fijate, asi tengo la estructura en vscode, pero sigue sin dejarme correr el proyecto (acompañado de la captura de pantalla de VS Code).
Prompt 5: pero hablame en español
Prompt 6: eso en q archivo tendria q ir?
Prompt 7: El bloque de código del archivo assets.go acompañado de la frase el archivo exe tiene este contenido, poneselo vos....
Prompt 8: El bloque de código y error con la frase ahora me tira este errorpackage vistas....
Prompt 9: ahora me tira este otro errorvistas\assets.go:15:12: pattern ../assets/sprites/*.png: invalid pattern syntax...
Prompt 10: package vistas... el codigo esta asi, pero sigue tirando este error... (con el código modificado previamente).
Prompt 11: fijate ese archivo, ayudame a agregarle hitbox a los autos (junto con el archivo adjunto f1game.zip).
Prompt 12: arma una lista con todos los promts q te tire, incluyendo este mismo (este mensaje actual).

CHAT GPT
Sprites de los autos
necesito q recrees esos sprites de autos de f1 para un juego, dejalos similares a un auto de f1 visto de arriba
Modificar el juego / vista / circuito / temporada
quiero q se vea desde arriba, como mirando un mapa, agregar curvas basadas en Monza, cambiar pilotos y autos para que sean solo de la temporada actual y cambiar las características de los pilotos y autos.
Agregar hitboxes a los rivales
fijate ese archivo, ayudame a agregarle hit boxs a los autos rivales
Revisar la implementación de las hitboxes
fjiate si asi quedo bien, sino corregilo
Error en assets.go / funciones de dibujo
vs code me tira errores al leer las funciones de dibujo, decime xq
Error expected declaration
y el error de vs code es expedted declaration, q nose muy bien q significa, la funcion esta declarada
Unir todos los prompts
necesito q juntes todos los promts q te fui tirando y los juntos en un solo texto
Aclaración de lo que quería recopilar
nono, quiero q hagas una lista con todos los promts q te tire yo, no lo q me respondiste vos, agrega este promtt tambien

CLAUDE 2
1. mira ese archivo  de un proyecto mio, hay q cambiarle la forma del circuito al juego, dame los codigos para hacerlo, no hagas un archico

2. interacciones\ia_bots.go:33:26: circuito.FactorCurva undefined (type entidades.Circuito has no field or method FactorCurva)
interacciones\ia_bots.go:67:26: circuito.FactorCurva undefined (type entidades.Circuito has no field or method FactorCurva)
PS C:\Users\fabri\Proyectos\Globant\f1game>

me tira ese error

3. ahora fijate el archivo de stats autos, dejalos todos casi iguales, para q tengan mas o menos el mismo rendimiento, fijate en la clasificacion actual de la f1 para asignarles unas stats acordes a cada piloto y auto

4. bien, ahora ayudame a cambiarle las estadisticas a los pilotos tambien

5. bien y ahora agregale algo q haga q el jugador no se mueva a menos q toque la tecla w

6. ahora dale un contador de 5 segundos al momento de inicar la carrera, xq sino el cambio es muy brusco entre eleccion de piloto y auto y la carrera en si

7. bien, ahora cambiale una cosa mas al circuito, quiero q no se vea la linea blanca y roja al medio de la pista, sino en los bordes

8. bien, ya casi terminamos, agregale ahora un tiempito chiquito de 2 segundos luego de tirar la segunda ruleta

9. cambialos vos y pasame el codigo de juego con los cambios

10. listo, y ahora necesito q agregues un cartel dependiendo de la posicion del jugador, agregalo cuando el juagdor finaliza la carrera y dale una duracion de 2 segundos, agregaselo a juego o al archivo q tengas q agregarlo y pasame el codigo entero

11. fjiate q me tira q gano la carerra sin importar la posicion donde haya terminado, corrgeilo y pasame el codigo

12. listo, ahora necesito q me separes todos los promts q te fui tirandi en este chat incluyendo este mismo y q los pongas todos en un mismo texto


