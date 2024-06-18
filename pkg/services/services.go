package services

import "github.com/moonen-home-automation/go-hasocket/internal"

type ServiceTarget struct {
	EntityId string `json:"entity_id"`
	DeviceId string `json:"device_id"`
	AreaId   string `json:"area_id"`
	LabelId  string `json:"label_id"`
}

type ServiceRequest struct {
	Id             int64          `json:"id"`
	RequestType    string         `json:"type"`
	Domain         string         `json:"domain"`
	Service        string         `json:"service"`
	ServiceData    map[string]any `json:"service_data,omitempty"`
	Target         ServiceTarget  `json:"target,omitempty"`
	ReturnResponse bool           `json:"return_response"`
}

func NewServiceRequest() ServiceRequest {
	id := internal.GetId()
	sr := ServiceRequest{
		Id:          id,
		RequestType: "call_service",
	}
	return sr
}

type ServiceResult struct {
	Id          int64  `json:"id"`
	RequestType string `json:"type"`
	Success     bool   `json:"success"`
	Result      struct {
		Context  map[string]any `json:"context"`
		Response map[string]any `json:"response"`
	} `json:"result"`
}
