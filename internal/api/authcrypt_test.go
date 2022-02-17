package api

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCrypt_EncodeDecodeUUID(t *testing.T) {
	id := uuid.New()

	crypt := NewAuthCrypt()

	enc, err := crypt.EncodeUUID(id)
	require.NoError(t, err)

	dec, err := crypt.DecodeToken(enc)
	require.NoError(t, err)
	require.Equal(t, id, dec, "Закодированный uuid %v не равен раскодированому %v", id.String(), dec.String())
}
