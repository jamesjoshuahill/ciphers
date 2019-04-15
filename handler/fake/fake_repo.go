package fake

import (
	"github.com/jamesjoshuahill/ciphers/repository"
)

type FakeRepo struct {
	StoreCall struct {
		Received struct {
			Secret repository.Secret
		}
		Returns struct {
			Error error
		}
	}
	FindByResourceIDCall struct {
		Received struct {
			ID string
		}
		Returns struct {
			Secret repository.Secret
			Error  error
		}
	}
}

func (r *FakeRepo) Store(cipher repository.Secret) error {
	r.StoreCall.Received.Secret = cipher
	return r.StoreCall.Returns.Error
}

func (r *FakeRepo) FindByID(id string) (repository.Secret, error) {
	r.FindByResourceIDCall.Received.ID = id
	return r.FindByResourceIDCall.Returns.Secret, r.FindByResourceIDCall.Returns.Error
}
