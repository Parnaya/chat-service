package entity

import "chat.service/model"

var Database []*model.Entity

func StubCreate(entity *model.Entity) {
	Database = append(Database, entity)
}

func StubUpdate(entity *model.Entity) {
	with(entity.Id, func(i int, found *model.Entity) {
		//found.Text = entity.Text
		found.Tags = entity.Tags
	})
}

func StubDelete(entity *model.Entity) {
	with(entity.Id, func(i int, _ *model.Entity) {
		Database[i] = nil
	})
}
