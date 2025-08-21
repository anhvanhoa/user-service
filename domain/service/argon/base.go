package argon

type Argon interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) (bool, error)
	SetParams(
		memory uint32,
		iterations uint32,
		parallelism uint8,
		saltLength uint32,
		keyLength uint32,
	) Argon
	SetMemory(memory uint32) Argon
	SetIterations(iterations uint32) Argon
	SetParallelism(parallelism uint8) Argon
	SetSaltLength(saltLength uint32) Argon
	SetKeyLength(keyLength uint32) Argon
}
