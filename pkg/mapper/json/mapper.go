package json

import (
	"github.com/siyu-nuro/go-example/pkg/entity"
	"github.com/siyu-nuro/go-example/pkg/model"
)

func MapCreateOrderResponse(orderID *entity.OrderID, err error) *model.CreateOrderResponse {
	if err != nil {
		return &model.CreateOrderResponse{
			Error: err,
		}
	}
	orderIDDummy := ""
	return &model.CreateOrderResponse{
		OrderID: &orderIDDummy,
	}
}

// MapModelOrderDetailToInputOrderDetail maps OrderDetail in model package to InputOrderDetail in entity package
func MapModelOrderDetailToInputOrderDetail(orderDetail *model.OrderDetail) *entity.InputOrderDetail {
	return &entity.InputOrderDetail{
		MerchantID:      orderDetail.MerchantID,
		MerchantOrderID: orderDetail.MerchantOrderID,
		PickupInfo: &entity.InputPickupInfo{
			Location:           mapModelLocationToEntityLocation(orderDetail.PickupInfo.Location),
			PhoneNumber:        orderDetail.PickupInfo.PhoneNumber,
			EarliestPickupTime: orderDetail.PickupInfo.EarliestPickupTime,
		},
		DropoffInfo: &entity.InputDropoffInfo{
			Location:            mapModelLocationToEntityLocation(orderDetail.DropoffInfo.Location),
			FirstName:           orderDetail.DropoffInfo.FirstName,
			LastName:            orderDetail.DropoffInfo.LastName,
			Email:               orderDetail.DropoffInfo.Email,
			PhoneNumber:         orderDetail.DropoffInfo.PhoneNumber,
			EarliestDropoffTime: orderDetail.DropoffInfo.EarliestDropoffTime,
			LatestDropoffTime:   orderDetail.DropoffInfo.LatestDropoffTime,
		},
		MerchantNotesOption: orderDetail.MerchantNotes,
	}
}

func mapModelLocationToEntityLocation(location *model.Location) *entity.Location {
	return &entity.Location{
		Geolocation: &entity.Geolocation{
			Latitude:  location.Geolocation.Latitude,
			Longitude: location.Geolocation.Longitude,
		},
		Address: &entity.Address{
			Street1:          location.Address.Street1,
			Street2:          location.Address.Street2,
			City:             location.Address.City,
			State:            location.Address.State,
			PostalCode:       location.Address.PostalCode,
			FormattedAddress: location.Address.FormattedAddress,
		},
	}
}
