package entity

var database []Entity

type Entity struct {
	id   string
	text []*string
	Tags []string
}

func Create(entity Entity) {
	database = append(database, entity)
}
