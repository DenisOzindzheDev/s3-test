package models

import "io"

type ImageUnit struct {
	User
	Payload     io.Reader
	PayloadName string
	PayloadSize int64
}
