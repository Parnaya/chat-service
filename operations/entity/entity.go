package entity

var Database []*Entity

type Entity struct {
	Id   string   `json:"id"`
	Data string   `json:"data"`
	Tags []string `json:"tags"` //here there is a tag - a type of message (common, voice, image...)
}

type Tag struct {
	Id   string
	Type string
	Data string
}

func Create(entity Entity) {
	Database = append(Database, &entity)
}

func Update(entity Entity) {
	with(entity.Id, func(i int, found *Entity) {
		found.Data = entity.Data
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
