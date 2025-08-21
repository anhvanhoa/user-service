package argonS

import (
	"cms-server/domain/service/argon"

	"github.com/alexedwards/argon2id"
)

type agronImpl struct {
	params *argon2id.Params
}

func NewArgon() argon.Argon {
	return &agronImpl{
		params: argon2id.DefaultParams,
	}
}

func (a *agronImpl) HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, a.params)
}

func (a *agronImpl) VerifyPassword(hashedPassword, password string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		return false, err
	}
	return match, nil
}

func (a *agronImpl) SetParams(
	memory uint32,
	iterations uint32,
	parallelism uint8,
	saltLength uint32,
	keyLength uint32,
) argon.Argon {
	a.params = &argon2id.Params{
		Memory:      memory,
		Iterations:  iterations,
		Parallelism: parallelism,
		SaltLength:  saltLength,
		KeyLength:   keyLength,
	}
	return a
}

func (a *agronImpl) SetMemory(memory uint32) argon.Argon {
	a.params.Memory = memory
	return a
}

func (a *agronImpl) SetIterations(iterations uint32) argon.Argon {
	a.params.Iterations = iterations
	return a
}

func (a *agronImpl) SetParallelism(parallelism uint8) argon.Argon {
	a.params.Parallelism = parallelism
	return a
}

func (a *agronImpl) SetSaltLength(saltLength uint32) argon.Argon {
	a.params.SaltLength = saltLength
	return a
}

func (a *agronImpl) SetKeyLength(keyLength uint32) argon.Argon {
	a.params.KeyLength = keyLength
	return a
}

func (a *agronImpl) GetParams() *argon2id.Params {
	return a.params
}
