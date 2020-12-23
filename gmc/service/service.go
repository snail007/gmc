package gservice

import (
	gconfig "github.com/snail007/gmc/config"
	gcore "github.com/snail007/gmc/core"
)



func (s *ServiceItem) AfterInit() func(srv gcore.ServiceItem) (err error) {
	return s.afterInit
}

func (s *ServiceItem) BeforeInit() func(srv gcore.Service, cfg *gconfig.Config) (err error) {
	return s.beforeInit
}

func (s *ServiceItem) SetBeforeInit(f func(srv gcore.Service, cfg *gconfig.Config) (err error)) gcore.ServiceItem {
	s.beforeInit = f
	return s
}

func (s *ServiceItem) SetAfterInit(f func(srv gcore.ServiceItem) (err error)) gcore.ServiceItem {
	s.afterInit = f
	return s
}

func (s *ServiceItem) GetService() gcore.Service {
	return s.service
}

func (s *ServiceItem) GetConfigID() string {
	return s.configID
}

func NewServiceItem(configID string, s gcore.Service) gcore.ServiceItem {
	return &ServiceItem{
		service:  s,
		configID: configID,
	}
}
