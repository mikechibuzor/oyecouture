package database

import "context"

type Pinger interface {
	Ping(context.Context) error
}
