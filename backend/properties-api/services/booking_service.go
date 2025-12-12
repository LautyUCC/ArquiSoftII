package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"properties-api/clients"
	"properties-api/domain"
	"properties-api/dto"
	"properties-api/repositories"
)

// BookingService define la interfaz para la lógica de negocio de bookings
type BookingService interface {
	// CreateBooking crea una nueva reserva
	// userId se obtiene del JWT (cualquier rol autenticado puede crear reservas)
	CreateBooking(ctx context.Context, createDTO dto.BookingCreateDTO, userID string) (dto.BookingDTO, error)
	// GetMyBookings obtiene todas las reservas del usuario autenticado
	// Ordenadas por fecha inicio descendente
	GetMyBookings(ctx context.Context, userID string) ([]dto.BookingDTO, error)
	// GetPropertyBookings obtiene todas las reservas CONFIRMED de una propiedad
	// Puede ser llamado por cualquier usuario autenticado (para bloquear calendario)
	GetPropertyBookings(ctx context.Context, propertyID string) ([]dto.BookingDTO, error)
	// UpdateBooking actualiza las fechas de una reserva
	// Solo el dueño de la reserva o un ADMIN puede actualizar
	// Valida solapamiento con otras reservas CONFIRMED
	UpdateBooking(ctx context.Context, bookingID string, updateDTO dto.BookingUpdateDTO, userID string, userRole string) (dto.BookingDTO, error)
	// DeleteBooking elimina una reserva
	// Solo el dueño de la reserva o un ADMIN puede eliminar
	// OPCIÓN ELEGIDA: Marca la reserva como CANCELLED (soft delete) en lugar de borrarla físicamente
	// Esto permite mantener un registro histórico de reservas canceladas
	DeleteBooking(ctx context.Context, bookingID string, userID string, userRole string) error
	// GetAllBookings obtiene todas las reservas del sistema (solo para admin)
	GetAllBookings(ctx context.Context) ([]dto.BookingDTO, error)
	// GetBookingByID obtiene una reserva por su ID
	// Útil para verificar ownership antes de actualizar/eliminar
	GetBookingByID(ctx context.Context, bookingID string) (*domain.Booking, error)
}

// bookingService es la implementación concreta de BookingService
type bookingService struct {
	bookingRepo  repositories.BookingRepository
	propertyRepo repositories.PropertyRepository
	rabbitClient clients.RabbitMQClient
}

// NewBookingService crea una nueva instancia del servicio de bookings
func NewBookingService(
	bookingRepo repositories.BookingRepository,
	propertyRepo repositories.PropertyRepository,
	rabbitClient clients.RabbitMQClient,
) BookingService {
	return &bookingService{
		bookingRepo:  bookingRepo,
		propertyRepo: propertyRepo,
		rabbitClient: rabbitClient,
	}
}

// CreateBooking crea una nueva reserva
// REGLA DE NEGOCIO: Solo requiere autenticación (cualquier rol: USER, ADMIN, OWNER)
// El userId se obtiene del JWT, no del body del request
// Esto asegura que un usuario solo pueda crear reservas para sí mismo
// Verifica solapamiento de reservas confirmadas y devuelve 409 si hay conflicto
func (s *bookingService) CreateBooking(ctx context.Context, createDTO dto.BookingCreateDTO, userID string) (dto.BookingDTO, error) {
	log.Printf("📅 [CreateBooking] Iniciando creación de reserva para usuario %s, propiedad %s", userID, createDTO.PropertyID)

	// 1. Validar que la propiedad existe
	property, err := s.propertyRepo.GetByID(createDTO.PropertyID)
	if err != nil {
		log.Printf("❌ [CreateBooking] Propiedad %s no existe: %v", createDTO.PropertyID, err)
		return dto.BookingDTO{}, fmt.Errorf("propiedad con ID '%s' no existe", createDTO.PropertyID)
	}
	log.Printf("✅ [CreateBooking] Propiedad %s encontrada: %s", createDTO.PropertyID, property.Title)

	// 2. Validar que la propiedad está disponible
	if !property.Available {
		log.Printf("❌ [CreateBooking] Propiedad %s no está disponible", createDTO.PropertyID)
		return dto.BookingDTO{}, fmt.Errorf("propiedad con ID '%s' no está disponible para reserva", createDTO.PropertyID)
	}

	// 3. Normalizar fechas (usar StartDate/EndDate si vienen, sino CheckIn/CheckOut para compatibilidad)
	var startDate, endDate time.Time
	startDateTime := createDTO.StartDate.Time()
	checkInTime := createDTO.CheckIn.Time()
	
	if !startDateTime.IsZero() {
		startDate = startDateTime
		endDate = createDTO.EndDate.Time()
	} else if !checkInTime.IsZero() {
		// Compatibilidad con frontend que puede enviar CheckIn/CheckOut
		startDate = checkInTime
		endDate = createDTO.CheckOut.Time()
	} else {
		return dto.BookingDTO{}, fmt.Errorf("se requiere startDate/endDate o checkIn/checkOut")
	}

	// Normalizar fechas a hora 00:00:00
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	// Validar fechas
	if startDate.After(endDate) || startDate.Equal(endDate) {
		log.Printf("❌ [CreateBooking] Fechas inválidas: startDate %v >= endDate %v", startDate, endDate)
		return dto.BookingDTO{}, fmt.Errorf("startDate debe ser anterior a endDate")
	}

	if startDate.Before(time.Now()) {
		log.Printf("❌ [CreateBooking] StartDate en el pasado: %v", startDate)
		return dto.BookingDTO{}, fmt.Errorf("startDate no puede ser en el pasado")
	}

	// 4. Verificar solapamiento de reservas confirmadas para esa propiedad
	existingBookings, err := s.bookingRepo.FindConfirmedByPropertyID(ctx, createDTO.PropertyID)
	if err != nil {
		log.Printf("⚠️ [CreateBooking] Error obteniendo reservas existentes: %v", err)
		return dto.BookingDTO{}, fmt.Errorf("error verificando disponibilidad: %w", err)
	}

		// Verificar si hay solapamiento con reservas confirmadas
		for _, existing := range existingBookings {
			if s.hasOverlap(startDate, endDate, existing.StartDate, existing.EndDate) {
				log.Printf("❌ [CreateBooking] Conflicto de reserva detectado: reserva %s solapa con fechas %v - %v", existing.ID.Hex(), existing.StartDate, existing.EndDate)
				// Retornar error con mensaje específico para que el controlador lo detecte como 409
				return dto.BookingDTO{}, errors.New("conflict: la propiedad ya tiene una reserva confirmada en esas fechas")
			}
		}
	log.Printf("✅ [CreateBooking] No hay conflictos de reserva para las fechas %v - %v", startDate, endDate)

	// 5. Calcular precio total
	// Precio por noche * número de noches
	nights := int(endDate.Sub(startDate).Hours() / 24)
	if nights <= 0 {
		log.Printf("❌ [CreateBooking] Número de noches inválido: %d", nights)
		return dto.BookingDTO{}, fmt.Errorf("debe haber al menos una noche entre startDate y endDate")
	}
	totalPrice := property.Price * float64(nights)
	log.Printf("💰 [CreateBooking] Precio total calculado: $%.2f (%d noches x $%.2f)", totalPrice, nights, property.Price)

	// 6. Crear booking con userId del JWT (no del body)

	booking := domain.Booking{
		PropertyID: createDTO.PropertyID,
		UserID:     userID, // userId del JWT, no del body
		StartDate:  startDate,
		EndDate:    endDate,
		TotalPrice: totalPrice,
		Status:     domain.BookingStatusConfirmed,
	}

	// 7. Guardar en repositorio
	err = s.bookingRepo.Create(ctx, &booking)
	if err != nil {
		log.Printf("❌ [CreateBooking] Error guardando reserva: %v", err)
		return dto.BookingDTO{}, fmt.Errorf("error creando reserva: %w", err)
	}
	log.Printf("✅ [CreateBooking] Reserva creada exitosamente con ID: %s", booking.ID.Hex())

	// 8. Publicar evento booking.created en RabbitMQ con información completa
	bookingID := booking.ID.Hex()
	if err := s.publishBookingCreatedEvent(bookingID, booking.PropertyID, booking.UserID, booking.StartDate, booking.EndDate); err != nil {
		// Log del error pero no fallar la operación si el evento no se publica
		// La reserva ya fue creada exitosamente
		log.Printf("⚠️ [CreateBooking] Error publicando evento 'booking.created' en RabbitMQ para reserva %s: %v", bookingID, err)
	} else {
		log.Printf("✅ [CreateBooking] Evento 'booking.created' publicado en RabbitMQ para reserva %s", bookingID)
	}

	// 9. Retornar DTO de respuesta
	// Mapear StartDate/EndDate del dominio a CheckIn/CheckOut en el DTO para compatibilidad con frontend
	return dto.BookingDTO{
		ID:         booking.ID.Hex(),
		PropertyID: booking.PropertyID,
		UserID:     booking.UserID,
		CheckIn:    booking.StartDate,
		CheckOut:   booking.EndDate,
		TotalPrice: booking.TotalPrice,
		Status:     booking.Status,
		CreatedAt:  booking.CreatedAt,
		UpdatedAt:  booking.UpdatedAt,
	}, nil
}

// hasOverlap verifica si dos rangos de fechas se solapan
// Retorna true si hay solapamiento, false en caso contrario
func (s *bookingService) hasOverlap(start1, end1, start2, end2 time.Time) bool {
	// Dos rangos se solapan si:
	// - start1 < end2 AND start2 < end1
	return start1.Before(end2) && start2.Before(end1)
}

// publishBookingCreatedEvent publica un evento booking.created en RabbitMQ con información completa
// Incluye: bookingId, propertyId, userId, startDate, endDate
func (s *bookingService) publishBookingCreatedEvent(bookingID, propertyID, userID string, startDate, endDate time.Time) error {
	// Usar PublishPropertyEvent para publicar el evento
	// El evento se publica con la operación "booking.created" y el bookingID
	// La cola "property_events" puede manejar tanto eventos de propiedades como de bookings
	// NOTA: Por ahora usamos el método existente, pero en el futuro se podría crear un método específico
	// que publique un JSON con toda la información del booking
	return s.rabbitClient.PublishPropertyEvent("booking.created", bookingID)
}

// GetMyBookings obtiene todas las reservas del usuario autenticado
// Ordenadas por fecha inicio descendente (más recientes primero)
func (s *bookingService) GetMyBookings(ctx context.Context, userID string) ([]dto.BookingDTO, error) {
	log.Printf("📋 [GetMyBookings] Obteniendo reservas para usuario %s", userID)

	bookings, err := s.bookingRepo.GetByUser(ctx, userID)
	if err != nil {
		log.Printf("❌ [GetMyBookings] Error obteniendo reservas: %v", err)
		return nil, fmt.Errorf("error obteniendo reservas: %w", err)
	}

	log.Printf("✅ [GetMyBookings] Encontradas %d reservas para usuario %s", len(bookings), userID)

	// Ordenar por fecha inicio descendente (más recientes primero)
	sort.Slice(bookings, func(i, j int) bool {
		return bookings[i].StartDate.After(bookings[j].StartDate)
	})

	// Convertir a DTOs
	// Mapear StartDate/EndDate del dominio a CheckIn/CheckOut en el DTO
	bookingDTOs := make([]dto.BookingDTO, len(bookings))
	for i, booking := range bookings {
		bookingDTOs[i] = dto.BookingDTO{
			ID:         booking.ID.Hex(),
			PropertyID: booking.PropertyID,
			UserID:     booking.UserID,
			CheckIn:    booking.StartDate,
			CheckOut:   booking.EndDate,
			TotalPrice: booking.TotalPrice,
			Status:     booking.Status,
			CreatedAt:  booking.CreatedAt,
			UpdatedAt:  booking.UpdatedAt,
		}
	}

	return bookingDTOs, nil
}

// GetPropertyBookings obtiene todas las reservas CONFIRMED de una propiedad
// Útil para bloquear el calendario en el frontend
// Puede ser llamado por cualquier usuario autenticado
func (s *bookingService) GetPropertyBookings(ctx context.Context, propertyID string) ([]dto.BookingDTO, error) {
	log.Printf("📋 [GetPropertyBookings] Obteniendo reservas CONFIRMED para propiedad %s", propertyID)

	// Verificar que la propiedad existe
	_, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		log.Printf("❌ [GetPropertyBookings] Propiedad %s no existe: %v", propertyID, err)
		return nil, fmt.Errorf("propiedad con ID '%s' no existe", propertyID)
	}

	// Obtener solo reservas CONFIRMED (para bloquear calendario)
	bookings, err := s.bookingRepo.FindConfirmedByPropertyID(ctx, propertyID)
	if err != nil {
		log.Printf("❌ [GetPropertyBookings] Error obteniendo reservas: %v", err)
		return nil, fmt.Errorf("error obteniendo reservas: %w", err)
	}

	log.Printf("✅ [GetPropertyBookings] Encontradas %d reservas CONFIRMED para propiedad %s", len(bookings), propertyID)

	// Convertir a DTOs
	// Mapear StartDate/EndDate del dominio a CheckIn/CheckOut en el DTO
	bookingDTOs := make([]dto.BookingDTO, len(bookings))
	for i, booking := range bookings {
		bookingDTOs[i] = dto.BookingDTO{
			ID:         booking.ID.Hex(),
			PropertyID: booking.PropertyID,
			UserID:     booking.UserID,
			CheckIn:    booking.StartDate,
			CheckOut:   booking.EndDate,
			TotalPrice: booking.TotalPrice,
			Status:     booking.Status,
			CreatedAt:  booking.CreatedAt,
			UpdatedAt:  booking.UpdatedAt,
		}
	}

	return bookingDTOs, nil
}

// GetBookingByID obtiene una reserva por su ID
// Útil para verificar ownership antes de actualizar/eliminar
func (s *bookingService) GetBookingByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo reserva: %w", err)
	}
	return booking, nil
}

// UpdateBooking actualiza las fechas de una reserva
// REGLAS DE AUTORIZACIÓN:
// - Solo el dueño de la reserva (userId coincide) o un ADMIN puede actualizar
// REGLAS DE NEGOCIO:
// - Valida que startDate < endDate
// - Verifica solapamiento con otras reservas CONFIRMED (excluyendo la reserva actual)
// - Si hay solapamiento → devuelve error 409 Conflict
func (s *bookingService) UpdateBooking(ctx context.Context, bookingID string, updateDTO dto.BookingUpdateDTO, userID string, userRole string) (dto.BookingDTO, error) {
	log.Printf("🔄 [UpdateBooking] Actualizando reserva %s para usuario %s (rol: %s)", bookingID, userID, userRole)

	// 1. Obtener la reserva existente
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("❌ [UpdateBooking] Reserva %s no existe: %v", bookingID, err)
		return dto.BookingDTO{}, fmt.Errorf("reserva con ID '%s' no existe", bookingID)
	}

	// 2. Verificar autorización: solo dueño o ADMIN puede actualizar
	if userRole != "ADMIN" && booking.UserID != userID {
		log.Printf("❌ [UpdateBooking] Usuario %s no tiene permisos para actualizar reserva %s (dueño: %s)", userID, bookingID, booking.UserID)
		return dto.BookingDTO{}, fmt.Errorf("no tienes permisos para actualizar esta reserva")
	}

	// 3. Normalizar fechas (usar StartDate/EndDate si vienen, sino CheckIn/CheckOut)
	var newStartDate, newEndDate time.Time
	if updateDTO.StartDate != nil {
		newStartDate = updateDTO.StartDate.Time()
	} else if updateDTO.CheckIn != nil {
		newStartDate = updateDTO.CheckIn.Time()
	} else {
		newStartDate = booking.StartDate // Mantener fecha actual si no se proporciona
	}

	if updateDTO.EndDate != nil {
		newEndDate = updateDTO.EndDate.Time()
	} else if updateDTO.CheckOut != nil {
		newEndDate = updateDTO.CheckOut.Time()
	} else {
		newEndDate = booking.EndDate // Mantener fecha actual si no se proporciona
	}

	// Normalizar fechas a hora 00:00:00
	newStartDate = time.Date(newStartDate.Year(), newStartDate.Month(), newStartDate.Day(), 0, 0, 0, 0, newStartDate.Location())
	newEndDate = time.Date(newEndDate.Year(), newEndDate.Month(), newEndDate.Day(), 0, 0, 0, 0, newEndDate.Location())

	// 4. Validar fechas
	if newStartDate.After(newEndDate) || newStartDate.Equal(newEndDate) {
		log.Printf("❌ [UpdateBooking] Fechas inválidas: startDate %v >= endDate %v", newStartDate, newEndDate)
		return dto.BookingDTO{}, fmt.Errorf("startDate debe ser anterior a endDate")
	}

	// 5. Verificar solapamiento con otras reservas CONFIRMED (excluyendo la reserva actual)
	existingBookings, err := s.bookingRepo.FindConfirmedByPropertyID(ctx, booking.PropertyID)
	if err != nil {
		log.Printf("⚠️ [UpdateBooking] Error obteniendo reservas existentes: %v", err)
		return dto.BookingDTO{}, fmt.Errorf("error verificando disponibilidad: %w", err)
	}

	// Verificar solapamiento excluyendo la reserva actual
	for _, existing := range existingBookings {
		if existing.ID.Hex() == bookingID {
			continue // Saltar la reserva actual
		}
		if s.hasOverlap(newStartDate, newEndDate, existing.StartDate, existing.EndDate) {
			log.Printf("❌ [UpdateBooking] Conflicto de reserva detectado: reserva %s solapa con fechas %v - %v", existing.ID.Hex(), existing.StartDate, existing.EndDate)
			return dto.BookingDTO{}, errors.New("conflict: la propiedad ya tiene una reserva confirmada en esas fechas")
		}
	}
	log.Printf("✅ [UpdateBooking] No hay conflictos de reserva para las nuevas fechas %v - %v", newStartDate, newEndDate)

	// 6. Recalcular precio total si las fechas cambiaron
	var totalPrice float64
	if newStartDate != booking.StartDate || newEndDate != booking.EndDate {
		property, err := s.propertyRepo.GetByID(booking.PropertyID)
		if err != nil {
			log.Printf("⚠️ [UpdateBooking] Error obteniendo propiedad para recalcular precio: %v", err)
			// Usar precio anterior si no se puede obtener la propiedad
			totalPrice = booking.TotalPrice
		} else {
			nights := int(newEndDate.Sub(newStartDate).Hours() / 24)
			totalPrice = property.Price * float64(nights)
			log.Printf("💰 [UpdateBooking] Precio total recalculado: $%.2f (%d noches x $%.2f)", totalPrice, nights, property.Price)
		}
	} else {
		totalPrice = booking.TotalPrice // Mantener precio si las fechas no cambiaron
	}

	// 7. Actualizar booking
	booking.StartDate = newStartDate
	booking.EndDate = newEndDate
	booking.TotalPrice = totalPrice
	booking.UpdatedAt = time.Now()

	err = s.bookingRepo.Update(ctx, booking)
	if err != nil {
		log.Printf("❌ [UpdateBooking] Error actualizando reserva: %v", err)
		return dto.BookingDTO{}, fmt.Errorf("error actualizando reserva: %w", err)
	}
	log.Printf("✅ [UpdateBooking] Reserva %s actualizada exitosamente", bookingID)

	// 8. Retornar DTO de respuesta
	return dto.BookingDTO{
		ID:         booking.ID.Hex(),
		PropertyID: booking.PropertyID,
		UserID:     booking.UserID,
		CheckIn:    booking.StartDate,
		CheckOut:   booking.EndDate,
		TotalPrice: booking.TotalPrice,
		Status:     booking.Status,
		CreatedAt:  booking.CreatedAt,
		UpdatedAt:  booking.UpdatedAt,
	}, nil
}

// DeleteBooking elimina una reserva
// REGLAS DE AUTORIZACIÓN:
// - Solo el dueño de la reserva (userId coincide) o un ADMIN puede eliminar
// OPCIÓN ELEGIDA: Marca la reserva como CANCELLED (soft delete) en lugar de borrarla físicamente
// Esto permite mantener un registro histórico de reservas canceladas para auditoría y análisis
func (s *bookingService) DeleteBooking(ctx context.Context, bookingID string, userID string, userRole string) error {
	log.Printf("🗑️ [DeleteBooking] Eliminando reserva %s para usuario %s (rol: %s)", bookingID, userID, userRole)

	// 1. Obtener la reserva existente
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("❌ [DeleteBooking] Reserva %s no existe: %v", bookingID, err)
		return fmt.Errorf("reserva con ID '%s' no existe", bookingID)
	}

	// 2. Verificar autorización: solo dueño o ADMIN puede eliminar
	if userRole != "ADMIN" && booking.UserID != userID {
		log.Printf("❌ [DeleteBooking] Usuario %s no tiene permisos para eliminar reserva %s (dueño: %s)", userID, bookingID, booking.UserID)
		return fmt.Errorf("no tienes permisos para eliminar esta reserva")
	}

	// 3. Marcar como CANCELLED (soft delete) en lugar de borrar físicamente
	booking.Status = domain.BookingStatusCancelled
	booking.UpdatedAt = time.Now()

	err = s.bookingRepo.Update(ctx, booking)
	if err != nil {
		log.Printf("❌ [DeleteBooking] Error marcando reserva como cancelada: %v", err)
		return fmt.Errorf("error eliminando reserva: %w", err)
	}

	log.Printf("✅ [DeleteBooking] Reserva %s marcada como CANCELLED exitosamente", bookingID)
	return nil
}

// GetAllBookings obtiene todas las reservas del sistema (solo para admin)
// No requiere parámetros de usuario ya que el middleware RequireAdmin garantiza acceso
func (s *bookingService) GetAllBookings(ctx context.Context) ([]dto.BookingDTO, error) {
	log.Printf("📋 [GetAllBookings] Obteniendo todas las reservas del sistema")

	// Obtener todas las reservas del repositorio
	bookings, err := s.bookingRepo.GetAll(ctx)
	if err != nil {
		log.Printf("❌ [GetAllBookings] Error obteniendo todas las reservas: %v", err)
		return nil, fmt.Errorf("error obteniendo todas las reservas: %w", err)
	}

	// Convertir a DTOs
	bookingDTOs := make([]dto.BookingDTO, 0, len(bookings))
	for _, booking := range bookings {
		bookingDTOs = append(bookingDTOs, dto.BookingDTO{
			ID:         booking.ID.Hex(),
			PropertyID: booking.PropertyID,
			UserID:     booking.UserID,
			CheckIn:    booking.StartDate,
			CheckOut:   booking.EndDate,
			TotalPrice: booking.TotalPrice,
			Status:     booking.Status,
			CreatedAt:  booking.CreatedAt,
			UpdatedAt:  booking.UpdatedAt,
		})
	}

	log.Printf("✅ [GetAllBookings] Obtenidas %d reservas", len(bookingDTOs))
	return bookingDTOs, nil
}

