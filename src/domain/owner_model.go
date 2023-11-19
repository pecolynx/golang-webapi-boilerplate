package domain

type OwnerModel interface {
	AppUserModel
}

type ownerModel struct {
	AppUserModel
}

func NewOwnerModel(appUser AppUserModel) (OwnerModel, error) {
	return &ownerModel{
		AppUserModel: appUser,
	}, nil
}
