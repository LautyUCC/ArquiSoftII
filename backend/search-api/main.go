package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"search-api/config"
	"search-api/consumers"
	"search-api/controllers"
	"search-api/repositories"
	"search-api/services"
)

func main() {
	log.Println("🚀 Iniciando Search API...")

	// ============================================
	// SECCIÓN 1: CARGAR CONFIGURACIÓN
	// ============================================
	log.Println("📋 Cargando configuración...")
	cfg := config.LoadConfig()
	log.Printf("✅ Configuración cargada:")
	log.Printf("   - Solr URL: %s", cfg.SolrURL)
	log.Printf("   - Memcached Host: %s", cfg.MemcachedHost)
	log.Printf("   - RabbitMQ URL: %s", cfg.RabbitMQURL)
	log.Printf("   - Properties API URL: %s", cfg.PropertiesAPIURL)
	log.Printf("   - Port: %s", cfg.Port)

	// ============================================
	// SECCIÓN 2: INICIALIZAR REPOSITORIOS
	// ============================================
	log.Println("📦 Inicializando repositorios...")

	// Inicializar repositorio de Solr
	solrRepo := repositories.NewSolrRepository(cfg.SolrURL)
	log.Println("✅ Repositorio de Solr inicializado")

	// Inicializar repositorio de caché
	cacheRepo := repositories.NewCacheRepository(cfg.MemcachedHost)
	log.Println("✅ Repositorio de caché inicializado")

	// ============================================
	// SECCIÓN 3: INICIALIZAR SERVICIO
	// ============================================
	log.Println("🔧 Inicializando servicio...")
	searchService := services.NewSearchService(solrRepo, cacheRepo, cfg.PropertiesAPIURL)
	log.Println("✅ Servicio de búsqueda inicializado")

	// ============================================
	// SECCIÓN 4: INICIALIZAR CONTROLADOR
	// ============================================
	log.Println("🎮 Inicializando controlador...")
	searchController := controllers.NewSearchController(searchService)
	log.Println("✅ Controlador de búsqueda inicializado")

	// ============================================
	// SECCIÓN 5: INICIALIZAR Y ARRANCAR CONSUMIDOR DE RABBITMQ
	// ============================================
	log.Println("🐰 Inicializando consumidor de RabbitMQ...")
	consumer, err := consumers.NewRabbitMQConsumer(cfg.RabbitMQURL, "property_events", searchService)
	if err != nil {
		log.Fatalf("❌ Error creando consumidor de RabbitMQ: %v", err)
	}
	defer func() {
		log.Println("🔌 Cerrando consumidor de RabbitMQ...")
		if err := consumer.Close(); err != nil {
			log.Printf("⚠️ Error cerrando consumidor de RabbitMQ: %v", err)
		}
	}()

	// Arrancar consumidor en una goroutine
	go func() {
		if err := consumer.Start(); err != nil {
			log.Fatalf("❌ Error iniciando consumidor de RabbitMQ: %v", err)
		}
	}()
	log.Println("✅ Consumidor de RabbitMQ iniciado en goroutine")

	// ============================================
	// SECCIÓN 6: CONFIGURAR ROUTER HTTP
	// ============================================
	log.Println("🛣️ Configurando rutas HTTP...")

	// Crear mux para las rutas
	mux := http.NewServeMux()

	// Registrar rutas
	mux.HandleFunc("/search", searchController.Search)
	mux.HandleFunc("/health", healthHandler)

	log.Println("✅ Rutas configuradas:")
	log.Println("   - GET /search")
	log.Println("   - GET /health")

	// ============================================
	// SECCIÓN 7: CONFIGURAR MIDDLEWARE DE CORS
	// ============================================
	log.Println("🌐 Configurando middleware de CORS...")

	// Handler con middleware de CORS
	handler := corsMiddleware(mux)

	// ============================================
	// SECCIÓN 8: CONFIGURAR SERVIDOR HTTP
	// ============================================
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ============================================
	// SECCIÓN 9: MANEJAR GRACEFUL SHUTDOWN
	// ============================================
	// Canal para recibir señales del sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		log.Println("🚀 =======================================")
		log.Printf("🚀 Search API corriendo en puerto %s", cfg.Port)
		log.Println("🚀 =======================================")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error iniciando servidor: %v", err)
		}
	}()

	// Esperar señal de terminación
	sig := <-sigChan
	log.Printf("📨 Señal recibida: %v. Iniciando shutdown graceful...", sig)

	// Crear contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Cerrar servidor gracefulmente
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Error durante shutdown del servidor: %v", err)
	} else {
		log.Println("✅ Servidor cerrado exitosamente")
	}

	log.Println("👋 Search API finalizada")
}

// healthHandler maneja las peticiones GET /health
// Retorna un JSON con el estado del servicio
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Solo permitir método GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Configurar headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Escribir respuesta JSON
	response := map[string]string{
		"status": "ok",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("⚠️ Error escribiendo respuesta de health: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// corsMiddleware agrega headers CORS a todas las respuestas
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Manejar preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r)
	})
}
