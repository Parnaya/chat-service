package entity

var Database []*Entity

type Entity struct {
	Id   string    `json:"id"`
	Text []*string `json:"text"`
	Tags []string  `json:"tags"`
}

func Create(entity Entity) {
	Database = append(Database, &entity)
}

func Update(entity Entity) {
	with(entity.Id, func(i int, found *Entity) {
		found.Text = entity.Text
		found.Tags = entity.Tags
	})
}

func Delete(entity Entity) {
	with(entity.Id, func(i int, _ *Entity) {
		Database[i] = nil
	})
}

func with(id string, block func(i int, entity *Entity)) {
	for i, v := range Database {
		if v.Id == id {
			block(i, v)
		}
	}
}
