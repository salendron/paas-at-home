package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type StorageInterface interface {
	SaveService(service *Service) error
	GetServiceByPathIdentifiers(name string, version string, address string) *Service
	GetAllServices() []*Service
	DeleteService(service *Service) (bool, error)
}

type CachedSQLiteStorage struct {
	SQLiteDBFile string
	Cache        []ServiceModel
}

func (s *CachedSQLiteStorage) Initialize(sqliteDBFile string) {
	s.SQLiteDBFile = sqliteDBFile
	s.syncCache()
}

func (s *CachedSQLiteStorage) connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(s.SQLiteDBFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&ServiceModel{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *CachedSQLiteStorage) syncCache() {
	serviceModels, err := s.getAllServiceModels()
	if err != nil {
		panic(err)
	}

	s.Cache = serviceModels
}

func (s *CachedSQLiteStorage) getAllServiceModels() ([]ServiceModel, error) {
	var serviceModels []ServiceModel

	db, err := s.connect()
	if err != nil {
		return serviceModels, err
	}

	result := db.Find(&serviceModels)
	if result.Error != nil {
		return serviceModels, result.Error
	}

	return serviceModels, nil
}

func (s *CachedSQLiteStorage) getServiceModelByPathIdentifiers(name string, version string, address string) *ServiceModel {
	for i := 0; i < len(s.Cache); i++ {
		item := s.Cache[i]
		if item.Name == name && item.Version == version && item.Address == address {
			return &item
		}
	}

	return nil
}

func (s *CachedSQLiteStorage) SaveService(service *Service) error {
	db, err := s.connect()
	if err != nil {
		return err
	}

	existingServiceModel := s.getServiceModelByPathIdentifiers(*service.Name, *service.Version, *service.Address)
	if existingServiceModel != nil {
		existingServiceModel.IsHealthy = true
		tx := db.Save(existingServiceModel)
		if tx.Error != nil {
			return tx.Error
		}
	} else {
		serviceModel := service.ToServiceModel()
		tx := db.Save(&serviceModel)
		if tx.Error != nil {
			return tx.Error
		}
	}

	go s.syncCache()

	return nil
}

func (s *CachedSQLiteStorage) GetServiceByPathIdentifiers(name string, version string, address string) *Service {
	sm := s.getServiceModelByPathIdentifiers(name, version, address)
	if sm == nil {
		return nil
	}

	service := sm.ToService()
	return &service
}

func (s *CachedSQLiteStorage) GetAllServices() []*Service {
	services := make([]*Service, len(s.Cache))

	for i := 0; i < len(s.Cache); i++ {
		service := s.Cache[i].ToService()
		services[i] = &service
	}

	return services
}

func (s *CachedSQLiteStorage) DeleteService(service *Service) (bool, error) {
	db, err := s.connect()
	if err != nil {
		return false, err
	}

	existingServiceModel := s.getServiceModelByPathIdentifiers(*service.Name, *service.Version, *service.Address)
	if existingServiceModel == nil {
		return false, nil
	}

	tx := db.Delete(&existingServiceModel)
	if tx.Error != nil {
		return true, tx.Error
	}

	return true, nil
}

type ServiceModel struct {
	gorm.Model
	DefinitionPath string
	Name           string
	Version        string
	Address        string
	IsHealthy      bool
	Latency        int
}

func (sm ServiceModel) ToService() Service {
	s := Service{}
	s.DefinitionPath = &sm.DefinitionPath
	s.Name = &sm.Name
	s.Version = &sm.Version
	s.Address = &sm.Address
	s.IsHealthy = &sm.IsHealthy
	s.Latency = &sm.Latency

	return s
}

func (s Service) ToServiceModel() ServiceModel {
	sm := ServiceModel{}
	sm.DefinitionPath = *s.DefinitionPath
	sm.Name = *s.Name
	sm.Version = *s.Version
	sm.Address = *s.Address
	sm.IsHealthy = *s.IsHealthy
	sm.Latency = *s.Latency

	return sm
}
