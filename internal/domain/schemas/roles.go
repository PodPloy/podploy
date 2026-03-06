package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"

	uuuid "github.com/PodPloy/podploy/pkg/uuid"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuuid.DefaultUUIDV7).Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		field.Bool("is_system").Default(false),

		field.UUID("organization_id", uuid.UUID{}).Optional(),

		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).
			Ref("roles").
			Unique().
			Field("organization_id"),

		edge.To("memberships", OrganizationMember.Type),
		edge.To("permissions", Permission.Type),
	}
}
