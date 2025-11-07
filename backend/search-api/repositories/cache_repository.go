package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/karlseguin/ccache/v3"
	"search-api/domain"
)

// CacheRepository define la interfaz para las operaciones de caché
type CacheRepository interface {
	// Get obtiene datos del caché (properties y total count)
	// Retorna (properties, total, found)
	Get(key string) ([]domain.Property, int, bool)

	// Set guarda datos en el caché con TTL
	Set(key string, properties []domain.Property, total int, ttl time.Duration)

	// Delete elimina datos del caché
	Delete(key string)
}

// cacheRepository es la implementación concreta de CacheRepository
// Implementa un sistema de caché de dos niveles: local (ccache) y distribuido (Memcached)
type cacheRepository struct {
	localCache     *ccache.Cache[string, *cacheData]
	memcachedClient *memcache.Client
}

// cacheData representa los datos almacenados en el caché
type cacheData struct {
	Properties []domain.Property `json:"properties"`
	Total      int               `json:"total"`
}

// NewCacheRepository crea una nueva instancia del repositorio de caché
// Inicializa ccache local y conecta con Memcached
func NewCacheRepository(memcachedHost string) CacheRepository {
	// Inicializar caché local con ccache
	localCache := ccache.New(ccache.Configure[string, *cacheData]().
		MaxSize(1000).
		ItemsToPrune(100))

	// Inicializar cliente de Memcached
	memcachedClient := memcache.New(memcachedHost)
	log.Printf("✅ Cliente de Memcached inicializado para %s", memcachedHost)

	return &cacheRepository{
		localCache:      localCache,
		memcachedClient: memcachedClient,
	}
}

// Get obtiene datos del caché con estrategia de dos niveles
// 1. Busca primero en caché local (ccache)
// 2. Si no está, busca en Memcached
// 3. Si está en Memcached, guarda en caché local
// Retorna (properties, total, found)
func (r *cacheRepository) Get(key string) ([]domain.Property, int, bool) {
	// Nivel 1: Buscar en caché local
	item := r.localCache.Get(key)
	if item != nil && !item.Expired() {
		data := item.Value()
		log.Printf("✅ Cache hit (local) para key: %s", key)
		return data.Properties, data.Total, true
	}

	// Nivel 2: Buscar en Memcached
	memcachedItem, err := r.memcachedClient.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			log.Printf("❌ Cache miss para key: %s", key)
			return nil, 0, false
		}
		log.Printf("⚠️ Error obteniendo de Memcached para key %s: %v", key, err)
		return nil, 0, false
	}

	// Deserializar datos de Memcached
	var data cacheData
	if err := json.Unmarshal(memcachedItem.Value, &data); err != nil {
		log.Printf("⚠️ Error deserializando datos de Memcached para key %s: %v", key, err)
		return nil, 0, false
	}

	// Guardar en caché local para próximas consultas (TTL de 5 minutos)
	r.localCache.Set(key, &data, 5*time.Minute)
	log.Printf("✅ Cache hit (Memcached) para key: %s, guardado en local", key)

	return data.Properties, data.Total, true
}

// Set guarda datos en ambos niveles de caché
// - Caché local: TTL de 5 minutos
// - Memcached: TTL de 15 minutos (o el TTL proporcionado si es mayor)
func (r *cacheRepository) Set(key string, properties []domain.Property, total int, ttl time.Duration) {
	data := &cacheData{
		Properties: properties,
		Total:      total,
	}

	// Guardar en caché local con TTL de 5 minutos
	r.localCache.Set(key, data, 5*time.Minute)
	log.Printf("✅ Datos guardados en caché local para key: %s", key)

	// Serializar para Memcached
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("⚠️ Error serializando datos para Memcached (key %s): %v", key, err)
		return
	}

	// Calcular TTL para Memcached (mínimo 15 minutos)
	memcachedTTL := ttl
	if memcachedTTL < 15*time.Minute {
		memcachedTTL = 15 * time.Minute
	}

	// Guardar en Memcached
	item := &memcache.Item{
		Key:        key,
		Value:      jsonData,
		Expiration: int32(memcachedTTL.Seconds()),
	}

	if err := r.memcachedClient.Set(item); err != nil {
		log.Printf("⚠️ Error guardando en Memcached (key %s): %v", key, err)
		return
	}

	log.Printf("✅ Datos guardados en Memcached para key: %s (TTL: %v)", key, memcachedTTL)
}

// Delete elimina datos de ambos niveles de caché
func (r *cacheRepository) Delete(key string) {
	// Eliminar de caché local
	r.localCache.Delete(key)
	log.Printf("✅ Datos eliminados de caché local para key: %s", key)

	// Eliminar de Memcached
	if err := r.memcachedClient.Delete(key); err != nil {
		if err == memcache.ErrCacheMiss {
			log.Printf("ℹ️ Key %s no existe en Memcached", key)
		} else {
			log.Printf("⚠️ Error eliminando de Memcached (key %s): %v", key, err)
		}
		return
	}

	log.Printf("✅ Datos eliminados de Memcached para key: %s", key)
}

