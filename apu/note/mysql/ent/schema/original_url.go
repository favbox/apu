package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type OriginalUrl struct {
	ent.Schema
}

func (OriginalUrl) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").Values("unset", "doc", "img").Default("unset"),
		field.String("url").MaxLen(512).Unique(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Optional().UpdateDefault(time.Now),
	}
}

func (OriginalUrl) Edges() []ent.Edge {
	return nil
}
