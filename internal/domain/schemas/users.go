package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"

	uuuid "github.com/PodPloy/podploy/pkg/uuid"
)

// Users holds the schema definition for the Users entity.
type Users struct {
	ent.Schema
}

// Fields of the Users.,
func (Users) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuuid.DefaultUUIDV7).Immutable(),
		field.String("email").Unique().NotEmpty(),
		field.String("password_hash").NotEmpty(),
		field.String("full_name").NotEmpty(),
		field.Bool("active").Default(true),

		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Users.
func (Users) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("memberships", OrganizationMember.Type),
	}
}
