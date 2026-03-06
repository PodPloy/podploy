package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"

	uuuid "github.com/PodPloy/podploy/pkg/uuid"
)

// Organization holds the schema definition for the Organization entity.
type Organization struct {
	ent.Schema
}

// Fields of the Organization.
func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuuid.DefaultUUIDV7).Immutable(),
		field.String("name").NotEmpty(),
		field.String("slug").Unique().NotEmpty(),
		field.Bool("active").Default(true),

		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Organization.
func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", Role.Type),
		edge.To("memberships", OrganizationMember.Type),
	}
}
