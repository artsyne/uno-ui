package pkg

import (
    "context"
    "testing"
	"fmt"
    "time"

    "github.com/biuren/hebe/pkg"
    core "github.com/biuren/hebe/tests/framework/core"
)

func setup() {
    // 初始化 GlobalHashID
    salt := "M5o9Q3s7Y4v2D8z6"
    minLength := 8
    pkg.GlobalHashID = pkg.NewHashID(salt, minLength)
    fmt.Println("GlobalHashID initialized with salt:", salt, "and minLength:", minLength)
}

func TestGlobalHashIDEncodeDecode(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
    defer cancel()

    select {
    case <-ctx.Done():
        t.Fatal("测试超时")
    default:
    }

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
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
    defer cancel()

    select {
    case <-ctx.Done():
        t.Fatal("测试超时")
    default:
    }

    setup()

    _, err := pkg.GlobalHashID.Decode("invalid")
    if err == nil {
        t.Fatalf("Expected error for invalid input, got nil")
    }
}