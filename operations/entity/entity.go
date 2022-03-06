package entity

var database []*Entity

type Entity struct {
	Id   string    `json:"id"`
	Text []*string `json:"text"`
	Tags []string  `json:"tags"`
}

func Create(entity Entity) {
	database = append(database, &entity)
}

func Update(entity Entity) {
	with(entity.Id, func(i int, found *Entity) {
		found.Text = entity.Text
		found.Tags = entity.Tags
	})
}

func Delete(entity Entity) {
	with(entity.Id, func(i int, _ *Entity) {
		database[i] = nil
	})
}

func with(id string, block func(i int, entity *Entity)) {
	for i, v := range database {
		if v.Id == id {
			block(i, v)
		}
	}
}
