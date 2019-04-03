package wsService

type IService interface {
	Start()
}

func NewClientService() IService {
	return &clientService{}
}

func NewServerService() IService {
	return &serverService{}
}
