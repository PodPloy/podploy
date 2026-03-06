//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("../graphql/ent.graphql"),
		entgql.WithConfigPath("../graphql/gqlgen.yml"),
		entgql.WithWhereInputs(true),
	)
	if err != nil {
		log.Fatalf("create extension entgql: %v", err)
	}

	cfg := &gen.Config{
		Target:  "../ent",
		Package: "github.com/PodPloy/podploy/internal/domain/ent",
	}

	err = entc.Generate("../schemas", cfg, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("corriendo ent codegen: %v", err)
	}
}
