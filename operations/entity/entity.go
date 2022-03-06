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
	for _, v := range database {
		if entity.id == v.id {
			v.text = entity.text
			v.tags = entity.tags
		}
	}
}

func Delete(entity Entity) {
	for i, v := range database {
		if entity.id == v.id {
			database[i] = nil
		}
	}
}