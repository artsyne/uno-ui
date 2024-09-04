package tests

import (
    "testing"
	"fmt"

    "github.com/biuren/hebe/pkg"
)

func setup() {
    // 初始化 GlobalHashID
    salt := "M5o9Q3s7Y4v2D8z6"
    minLength := 8
    pkg.GlobalHashID = pkg.NewHashID(salt, minLength)
    fmt.Println("GlobalHashID initialized with salt:", salt, "and minLength:", minLength)
}

func TestGlobalHashIDEncodeDecode(t *testing.T) {
    setup()

    id := 12345744645
    fmt.Println("Testing Encode and Decode with ID:", id)

    encoded, err := pkg.GlobalHashID.Encode(id)
    if err != nil {
        t.Fatalf("Encode failed: %v", err)
    }
    fmt.Println("Encoded ID:", encoded)

    decoded, err := pkg.GlobalHashID.Decode(encoded)
    if err != nil {
        t.Fatalf("Decode failed: %v", err)
    }
    fmt.Println("Decoded ID:", decoded)

    if decoded != id {
        t.Errorf("Expected decoded ID %d, got %d", id, decoded)
    }
}

func TestGlobalHashIDDecodeInvalid(t *testing.T) {
    setup()

    _, err := pkg.GlobalHashID.Decode("invalid")
    if err == nil {
        t.Fatalf("Expected error for invalid input, got nil")
    }
}