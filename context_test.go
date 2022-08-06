package autho

import (
	"context"
	"errors"
	"testing"
)

func TestContextWithError(t *testing.T) {
	expectedErr := errors.New("expected error")
	ctx := ContextWithError(context.Background(), expectedErr)
	if err := ErrorFromContext(ctx); err.Error() != expectedErr.Error() {
		t.Fatalf("expected expectedErr: expected error but got %s", err.Error())
	}
}

func TestContextWithUser(t *testing.T) {
	type userType struct {
		name string
	}
	expectedUser := &userType{
		name: "testing",
	}
	ctx := ContextWithUser(context.Background(), expectedUser)
	user, err := UserFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	gotUser, ok := user.(*userType)
	if !ok {
		t.Fatal("expected user to be set")
	}
	if gotUser.name != expectedUser.name {
		t.Fatalf("expected user name to be testing but got %s", gotUser.name)
	}
}
