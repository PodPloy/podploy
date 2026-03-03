package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// OrganizationMember holds the schema definition for the join table.
type OrganizationMember struct {
	ent.Schema
}

// Fields of the OrganizationMember.
func (OrganizationMember) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("organization_id", uuid.UUID{}),
		field.UUID("role_id", uuid.UUID{}),

		field.Time("joined_at").Default(time.Now).Immutable(),
	}
}

// Edges of the OrganizationMember.
func (OrganizationMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", Users.Type).Ref("memberships").Unique().Field("user_id").Required(),
		edge.From("organization", Organization.Type).Ref("memberships").Unique().Field("organization_id").Required(),
		edge.From("role", Role.Type).Ref("memberships").Unique().Field("role_id").Required(),
	}
}

// Indexes of the OrganizationMember.
func (OrganizationMember) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "organization_id").Unique(),
	}
}
