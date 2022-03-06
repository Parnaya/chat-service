package entity

var database []*Entity

type Entity struct {
	id   string
	text []*string
	tags []string
}

func Create(entity Entity) {
	database = append(database, &entity)
}

func Update(entity Entity) {
	with(entity.id, func(i int, found *Entity) {
		found.text = entity.text
		found.tags = entity.tags
	})
}

func Delete(entity Entity) {
	with(entity.id, func(i int, _ *Entity) {
		database[i] = nil
	})
}

func with(id string, block func(i int, entity *Entity)) {
	for i, v := range database {
		if v.id == id {
			block(i, v)
		}
	}
}
