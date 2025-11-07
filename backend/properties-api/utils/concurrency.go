package utils

import (
	"sync"
)

// CalculatePriceWithConcurrency calcula el precio total de una propiedad usando goroutines
// Divide el cálculo en 3 partes que se ejecutan en paralelo para mejorar el rendimiento:
// 1. Precio base con impuestos (21%)
// 2. Costo adicional por amenidades ($50 cada una)
// 3. Costo adicional por capacidad ($30 por persona)
// Finalmente suma todos los resultados para obtener el precio total
func CalculatePriceWithConcurrency(basePrice float64, amenities []string, capacity int) float64 {
	// WaitGroup permite esperar a que todas las goroutines terminen
	// Se usa para sincronizar las operaciones concurrentes
	var wg sync.WaitGroup

	// Channel para recibir los resultados de cada goroutine
	// Permite comunicación segura entre goroutines
	resultsChan := make(chan float64, 3) // Buffer de 3 para evitar bloqueos

	// Goroutine 1: Calcula el precio base con impuestos del 21%
	// Aplica el impuesto al precio base y envía el resultado al channel
	wg.Add(1) // Incrementar el contador del WaitGroup
	go func() {
		defer wg.Done() // Decrementar el contador cuando termine la goroutine

		// Calcular precio base con impuestos (21% = 0.21)
		// Precio con impuestos = precio base * 1.21
		priceWithTaxes := basePrice * 1.21

		// Enviar resultado al channel
		resultsChan <- priceWithTaxes
	}()

	// Goroutine 2: Calcula el costo adicional por amenidades
	// Cada amenidad tiene un costo de $50
	// Suma todos los costos de las amenidades y envía el resultado al channel
	wg.Add(1) // Incrementar el contador del WaitGroup
	go func() {
		defer wg.Done() // Decrementar el contador cuando termine la goroutine

		// Calcular costo total de amenidades
		// $50 por cada amenidad en la lista
		amenityCost := 50.0
		totalAmenityCost := float64(len(amenities)) * amenityCost

		// Enviar resultado al channel
		resultsChan <- totalAmenityCost
	}()

	// Goroutine 3: Calcula el costo adicional por capacidad
	// Cada persona de capacidad tiene un costo de $30
	// Multiplica la capacidad por el costo unitario y envía el resultado al channel
	wg.Add(1) // Incrementar el contador del WaitGroup
	go func() {
		defer wg.Done() // Decrementar el contador cuando termine la goroutine

		// Calcular costo total por capacidad
		// $30 por cada persona de capacidad
		capacityCost := 30.0
		totalCapacityCost := float64(capacity) * capacityCost

		// Enviar resultado al channel
		resultsChan <- totalCapacityCost
	}()

	// Esperar a que todas las goroutines terminen
	// Wait() bloquea hasta que el contador del WaitGroup llegue a 0
	// Esto garantiza que todas las goroutines hayan enviado sus resultados al channel
	wg.Wait()

	// Cerrar el channel cuando todas las goroutines hayan terminado
	// Esto es importante para poder hacer un range sobre el channel
	// Sin cerrar el channel, el range nunca terminaría
	close(resultsChan)

	// Sumar todos los resultados del channel
	// Range sobre el channel lee todos los valores hasta que se cierre
	totalPrice := 0.0
	for result := range resultsChan {
		totalPrice += result
	}

	// Retornar el precio total calculado
	return totalPrice
}

