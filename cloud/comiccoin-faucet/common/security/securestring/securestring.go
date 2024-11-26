package securestring

import (
	"errors"
	"fmt"

	"github.com/awnumar/memguard"
)

// SecureString is used to store a string securely in memory.
type SecureString struct {
	buffer *memguard.LockedBuffer
}

// NewSecureString creates a new SecureString instance from the given string.
func NewSecureString(s string) (*SecureString, error) {
	if len(s) == 0 {
		return nil, errors.New("string cannot be empty")
	}

	buffer := memguard.NewBuffer(len(s))

	copy(buffer.Bytes(), s)

	buffer.Lock()

	return &SecureString{buffer: buffer}, nil
}

// String returns the securely stored string.
func (ss *SecureString) String() string {
	return string(ss.buffer.Bytes())
}

// String returns the securely stored string.
func (ss *SecureString) Bytes() []byte {
	return ss.buffer.Bytes()
}

// Wipe removes the string from memory and makes it unrecoverable.
func (ss *SecureString) Wipe() error {
	fmt.Println("SecureString: Wipe(): Starting...")
	ss.buffer.Wipe()
	ss.buffer = nil // Make sure to set the buffer to nil after wiping.
	fmt.Println("SecureString: Wipe(): Finished")
	return nil
}
