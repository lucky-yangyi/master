package services

import (
	"fenqi_v1/models"
	"strconv"
	"sync"
)

func ComputeValue(s *models.Salesman, salesmanMap map[string]models.SysOrganization, finish *sync.WaitGroup) {
	s.IsLinkAllotment, _ = models.GetIslinkAlltomentCount(s.Id)
	s.NotIsLinkAllotment, _ = models.GetNotIslinkAlltomentCount(s.Id)
	if s.DepId != 0 {
		s.Place, s.OperaDep, s.Region, _ = FindOrgTypeDepartment(salesmanMap, strconv.Itoa(s.DepId))
	} else if s.PlaceId != 0 {
		s.Region, s.Place, _ = FindOrgTypePlace(salesmanMap, strconv.Itoa(s.PlaceId))
	} else if s.RegionId != 0 {
		_, s.Region, _ = FindOrgTypePlace(salesmanMap, strconv.Itoa(s.RegionId))
	}
	finish.Done()
}
