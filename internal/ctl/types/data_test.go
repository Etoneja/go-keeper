package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginData_Validate(t *testing.T) {
	t.Run("valid login data", func(t *testing.T) {
		data := LoginData{
			Username: "user123",
			Password: "pass123",
			URL:      "https://example.com",
		}
		assert.NoError(t, data.Validate())
	})

	t.Run("valid login data without URL", func(t *testing.T) {
		data := LoginData{
			Username: "user123",
			Password: "pass123",
		}
		assert.NoError(t, data.Validate())
	})

	t.Run("empty username", func(t *testing.T) {
		data := LoginData{Username: "", Password: "pass123"}
		assert.Error(t, data.Validate())
	})

	t.Run("empty password", func(t *testing.T) {
		data := LoginData{Username: "user123", Password: ""}
		assert.Error(t, data.Validate())
	})

	t.Run("whitespace username", func(t *testing.T) {
		data := LoginData{Username: "   ", Password: "pass123"}
		assert.Error(t, data.Validate())
	})

	t.Run("whitespace password", func(t *testing.T) {
		data := LoginData{Username: "user123", Password: "   "}
		assert.Error(t, data.Validate())
	})
}

func TestTextData_Validate(t *testing.T) {
	t.Run("valid text data", func(t *testing.T) {
		data := TextData{Content: "valid content"}
		assert.NoError(t, data.Validate())
	})

	t.Run("empty content", func(t *testing.T) {
		data := TextData{Content: ""}
		assert.Error(t, data.Validate())
	})

	t.Run("whitespace content", func(t *testing.T) {
		data := TextData{Content: "   "}
		assert.Error(t, data.Validate())
	})
}

func TestFileData_Validate(t *testing.T) {
	t.Run("valid file data", func(t *testing.T) {
		data := FileData{
			FileName: "file.txt",
			FileSize: 100,
			Content:  "content",
		}
		assert.NoError(t, data.Validate())
	})

	t.Run("empty file name", func(t *testing.T) {
		data := FileData{FileName: "", Content: "content"}
		assert.Error(t, data.Validate())
	})

	t.Run("empty content", func(t *testing.T) {
		data := FileData{FileName: "file.txt", Content: ""}
		assert.Error(t, data.Validate())
	})

	t.Run("whitespace file name", func(t *testing.T) {
		data := FileData{FileName: "   ", Content: "content"}
		assert.Error(t, data.Validate())
	})

	t.Run("whitespace content", func(t *testing.T) {
		data := FileData{FileName: "file.txt", Content: "   "}
		assert.Error(t, data.Validate())
	})
}

func TestCardData_Validate(t *testing.T) {
	t.Run("valid card data with 3-digit CVV", func(t *testing.T) {
		data := CardData{
			Number: "4111111111111111",
			Holder: "John Doe",
			Expiry: "12/25",
			CVV:    "123",
		}
		assert.NoError(t, data.Validate())
	})

	t.Run("valid card data with 4-digit CVV", func(t *testing.T) {
		data := CardData{
			Number: "4111111111111111",
			Holder: "John Doe",
			Expiry: "12/2025",
			CVV:    "1234",
		}
		assert.NoError(t, data.Validate())
	})

	t.Run("valid card data with spaces in number", func(t *testing.T) {
		data := CardData{
			Number: "4111 1111 1111 1111",
			Holder: "John Doe",
			Expiry: "12/25",
			CVV:    "123",
		}
		assert.NoError(t, data.Validate())
	})

	t.Run("invalid card number", func(t *testing.T) {
		data := CardData{
			Number: "invalid",
			Holder: "John Doe",
			Expiry: "12/25",
			CVV:    "123",
		}
		assert.Error(t, data.Validate())
	})

	t.Run("empty card holder", func(t *testing.T) {
		data := CardData{
			Number: "4111111111111111",
			Holder: "",
			Expiry: "12/25",
			CVV:    "123",
		}
		assert.Error(t, data.Validate())
	})

	t.Run("invalid expiry format", func(t *testing.T) {
		data := CardData{
			Number: "4111111111111111",
			Holder: "John Doe",
			Expiry: "invalid",
			CVV:    "123",
		}
		assert.Error(t, data.Validate())
	})

	t.Run("empty CVV", func(t *testing.T) {
		data := CardData{
			Number: "4111111111111111",
			Holder: "John Doe",
			Expiry: "12/25",
			CVV:    "",
		}
		assert.Error(t, data.Validate())
	})

	t.Run("invalid CVV format", func(t *testing.T) {
		data := CardData{
			Number: "4111111111111111",
			Holder: "John Doe",
			Expiry: "12/25",
			CVV:    "abc",
		}
		assert.Error(t, data.Validate())
	})

	t.Run("invalid CVV length", func(t *testing.T) {
		data := CardData{
			Number: "4111111111111111",
			Holder: "John Doe",
			Expiry: "12/25",
			CVV:    "12",
		}
		assert.Error(t, data.Validate())
	})
}
