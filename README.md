# Proyecto Blockchain en Go

Este proyecto es una implementación simple de una cadena de bloques (blockchain) escrita en Go.

Esta rama feature/transactions implementa lógica de transacciones diferente a la versión en la rama dev/main. 

Implementado para presentarse en la asignatura de Criptografía y Ciberseguridad en la UA.  

## Características

- Implementación de Prueba de Trabajo (Proof of Work)
- Persistencia de datos con BadgerDB
- Interfaz de línea de comandos (CLI)

## Estructura del Proyecto

El proyecto se divide en varios archivos:

- `block.go`: Define la estructura de un bloque y proporciona funciones para crear y deserializar bloques.
- `transacction.go`: Implementa la estructura y lógica para el manejo de las transacciones
- `blockchain.go`: Implementa la lógica de la cadena, funciones para agregarle bloques y recorrerlos.
- `proof.go`: Implementa la Prueba de Trabajo (Proof of Work) y proporciona funciones para validar la prueba.
- `main.go`: Contiene la función main del programa y la interfaz de línea de comandos (CLI).

## Cómo usar

Se puede ejecutar el programa con `go run main.go`

También puede compilarse comúnmente con `go build`. Y luego llamar su ejecutable como corresponda.

El programa se puede utilizar a través de la interfaz de la CLI con estos comandos:

- `printchain`: Muestra en pantalla un resumen de la cadena de bloques.
- `createblockchain -address ADDRESS`: crea la cadena y agrega el bloque génesis
- `send -from FROM -to TO -amount AMOUNT`: envía AMOUNT desde FROM hacia TO
- `getbalance -address ADDRESS`: permite ver el balance para ADDRESS

## Requisitos

- Go versión 1.16 o superior
- BadgerDB