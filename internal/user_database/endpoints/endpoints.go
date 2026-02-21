package endpoints

import "github.com/vwency/microservices_golang/internal/user_database/service"

type Endpoints struct {
	AddUser      *AddUserEndpoint
	GetUser      *GetUserEndpoint
	UpdateUser   *UpdateUserEndpoint
	DeleteUser   *DeleteUserEndpoint
	GetUserByID  *GetUserByIDEndpoint
	InitDatabase *InitDatabaseEndpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		AddUser:      &AddUserEndpoint{svc: s},
		GetUser:      &GetUserEndpoint{svc: s},
		UpdateUser:   &UpdateUserEndpoint{svc: s},
		DeleteUser:   &DeleteUserEndpoint{svc: s},
		GetUserByID:  &GetUserByIDEndpoint{svc: s},
		InitDatabase: &InitDatabaseEndpoint{svc: s},
	}
}
