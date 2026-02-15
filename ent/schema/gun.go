package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Gun holds the schema definition for the Gun entity.
type Gun struct {
	ent.Schema
}

// Annotations of the Gun.
func (Gun) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "guns"},
	}
}

// Fields of the Gun.
func (Gun) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("gun_id"),
		field.String("gun_name").MaxLen(255),
		field.Int("year").Optional().Nillable(),
		field.Int("condition").Optional().Nillable(),
		field.String("serial_number").Optional().MaxLen(255),
		field.String("description").Optional().MaxLen(255),
		field.Bytes("image").Optional().Nillable(),
		field.String("misc_attachments").Optional().MaxLen(255),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			StorageKey("createdAt"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StorageKey("updatedAt"),
	}
}

// Edges of the Gun.
func (Gun) Edges() []ent.Edge {
	return nil
}
