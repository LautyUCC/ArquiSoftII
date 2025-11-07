package utils

import (
	"fmt"
	"strings"
)

// ValidatePropertyType valida que el tipo de propiedad sea válido
func ValidatePropertyType(propertyType string) error {
	validTypes := []string{"casa", "apartamento", "terreno", "local", "oficina"}
	propertyType = strings.ToLower(strings.TrimSpace(propertyType))

	for _, validType := range validTypes {
		if propertyType == validType {
			return nil
		}
	}

	return fmt.Errorf("tipo de propiedad inválido. Tipos válidos: %s", strings.Join(validTypes, ", "))
}

// ValidatePropertyStatus valida que el estado de la propiedad sea válido
func ValidatePropertyStatus(status string) error {
	validStatuses := []string{"disponible", "vendido", "alquilado", "reservado"}
	status = strings.ToLower(strings.TrimSpace(status))

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return fmt.Errorf("estado de propiedad inválido. Estados válidos: %s", strings.Join(validStatuses, ", "))
}

// ValidatePrice valida que el precio sea válido
func ValidatePrice(price float64) error {
	if price <= 0 {
		return fmt.Errorf("el precio debe ser mayor a 0")
	}
	return nil
}

// ValidateArea valida que el área sea válida
func ValidateArea(area float64) error {
	if area <= 0 {
		return fmt.Errorf("el área debe ser mayor a 0")
	}
	return nil
}

