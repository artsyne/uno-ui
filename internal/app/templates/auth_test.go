package templates

import (
	"testing"

	types "github.com/biuren/hebe/internal/contracts"
)

func TestShowUserName(t *testing.T) {
	// 1. Struct with DisplayName
	type User1 struct {
		DisplayName string
		Handle      string
	}
	u1 := User1{DisplayName: "Alice", Handle: "alice123"}
	if got := showUserName(u1); got != "Alice" {
		t.Errorf("showUserName(User1) = %q, want %q", got, "Alice")
	}

	// 2. Struct with DisplayName empty, Handle set
	u2 := User1{DisplayName: "", Handle: "bob123"}
	if got := showUserName(u2); got != "bob123" {
		t.Errorf("showUserName(User1) = %q, want %q", got, "bob123")
	}

	// 3. Struct with AddUserDisplayName
	type Entry1 struct {
		AddUserDisplayName string
		AddUserHandle      string
	}
	e1 := Entry1{AddUserDisplayName: "Charlie", AddUserHandle: "charlie123"}
	if got := showUserName(e1); got != "Charlie" {
		t.Errorf("showUserName(Entry1) = %q, want %q", got, "Charlie")
	}

	// 4. Struct with AddUserDisplayName empty
	e2 := Entry1{AddUserDisplayName: "", AddUserHandle: "dave123"}
	if got := showUserName(e2); got != "dave123" {
		t.Errorf("showUserName(Entry1) = %q, want %q", got, "dave123")
	}

	// 5. PageData with Context.AuthUser
	pd := &types.PageData[map[string]interface{}]{
		Context: &types.RenderContext{
			AuthUser: &types.AuthUser{
				DisplayName: "Eve",
				Handle:      "eve123",
			},
		},
	}
	if got := showUserName(pd); got != "Eve" {
		t.Errorf("showUserName(PageData) = %q, want %q", got, "Eve")
	}

    // 5b. PageData with Context.AuthUser (fallback to Handle)
    pd2 := &types.PageData[map[string]interface{}]{
        Context: &types.RenderContext{
            AuthUser: &types.AuthUser{
                DisplayName: "",
                Handle:      "frank123",
            },
        },
    }
    if got := showUserName(pd2); got != "frank123" {
        t.Errorf("showUserName(PageData fallback) = %q, want %q", got, "frank123")
    }

	// 6. Nil
	if got := showUserName(nil); got != "" {
		t.Errorf("showUserName(nil) = %q, want %q", got, "")
	}

	// 7. Unknown struct
	type Unknown struct {
		Foo string
	}
	if got := showUserName(Unknown{Foo: "bar"}); got != "" {
		t.Errorf("showUserName(Unknown) = %q, want %q", got, "")
	}
    
    // 8. Pointer to struct
    u3 := &User1{DisplayName: "George", Handle: "george123"}
    if got := showUserName(u3); got != "George" {
        t.Errorf("showUserName(*User1) = %q, want %q", got, "George")
    }
}
