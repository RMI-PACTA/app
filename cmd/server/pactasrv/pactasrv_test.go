package pactasrv

import (
	"context"
	"testing"

	"github.com/RMI/pacta/openapi/pacta"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func TestPets(t *testing.T) {
	srv := &Server{}

	ctx := context.Background()
	tkn := jwt.New()
	tkn.Set("sub", "user123")
	ctx = jwtauth.NewContext(ctx, tkn, nil)

	{
		got, err := srv.FindPets(ctx, pacta.FindPetsRequestObject{
			Params: pacta.FindPetsParams{},
		})
		if err != nil {
			t.Fatalf("srv.CreateRun: %v", err)
		}

		want := pacta.FindPets200JSONResponse{}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected create run response (-want +got)\n%s", diff)
		}
	}

	{
		got, err := srv.AddPet(ctx, pacta.AddPetRequestObject{
			Body: &pacta.NewPet{
				Name: "Scruffles",
				Tag:  ptr("good boy"),
			},
		})
		if err != nil {
			t.Fatalf("srv.CreateRun: %v", err)
		}

		want := pacta.AddPet200JSONResponse{
			Id:   1,
			Name: "Scruffles",
			Tag:  ptr("good boy"),
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected create run response (-want +got)\n%s", diff)
		}
	}

	{
		got, err := srv.FindPets(ctx, pacta.FindPetsRequestObject{
			Params: pacta.FindPetsParams{},
		})
		if err != nil {
			t.Fatalf("srv.CreateRun: %v", err)
		}

		want := pacta.FindPets200JSONResponse{
			pacta.Pet{
				Id:   1,
				Name: "Scruffles",
				Tag:  ptr("good boy"),
			},
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected create run response (-want +got)\n%s", diff)
		}
	}
}

func ptr[T any](in T) *T {
	return &in
}
