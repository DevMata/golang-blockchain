# Proyecto Blockchain en Go

Este proyecto es una implementación simple de una cadena de bloques (blockchain) escrita en Go.

Implementado para presentarse en la asignatura de Criptografía y Ciberseguridad en la UA.  

## Características

- Implementación de Prueba de Trabajo (Proof of Work)
- Persistencia de datos con BadgerDB
- Interfaz de línea de comandos (CLI)

## Estructura del Proyecto

El proyecto se divide en varios archivos:

- `block.go`: Define la estructura de un bloque y proporciona funciones para crear y deserializar bloques.
- `blockchain.go`: Implementa la lógica de la cadena, funciones para agregarle bloques y recorrerlos.
- `proof.go`: Implementa la Prueba de Trabajo (Proof of Work) y proporciona funciones para validar la prueba.
- `main.go`: Contiene la función main del programa y la interfaz de línea de comandos (CLI).

## Cómo usar

Se puede ejecutar el programa con `go run main.go`

También puede compilarse comúnmente con `go build`. Y luego llamar su ejecutable como corresponda.

El programa se puede utilizar a través de la interfaz de la CLI con estos comandos:

- `add -block BLOCK_DATA`: Agrega un bloque a la cadena con los datos proporcionados.
- `print`: Muestra en pantalla el contenido de la cadena de bloques.

## Requisitos

- Go versión 1.16 o superior
- BadgerDB