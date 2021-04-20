package db

type Driver interface {
	GetName() string
	GetVersion() string
}
