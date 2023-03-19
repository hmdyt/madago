package encoder

import "github.com/hmdyt/madago/domain/entities"

type IEncoder interface {
	Write(events []*entities.Event)
}
