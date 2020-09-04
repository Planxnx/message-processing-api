package connection

import (
	"log"
)

type Connections struct {
	MockCon *Mock
}

func InitializeConnection() *Connections {
	log.Println("Initialize Conections")

	return &Connections{
		MockCon: InitializeMock(),
	}
}
