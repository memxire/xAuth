package tests

import (
	"testing"
	"xauth/tests/suite"

	"github.com/brianvoe/gofakeit/v6"
	ssov1 "github.com/memxire/protobuf/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	userID     = 2
	emptyValue = 0
	notExists  = 99999
)

func TestGetUserByID_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    gofakeit.Email(),
		Password: randomFakePassword(),
		Username: gofakeit.Username(),
	})
	require.NoError(t, err)
	userID := respReg.GetUserId()
	require.NotEmpty(t, userID)

	respGet, err := st.AuthClient.GetUser(ctx, &ssov1.GetUserRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	assert.Equal(t, userID, respGet.GetUserId())
	assert.NotEmpty(t, respGet.GetEmail())
	assert.NotEmpty(t, respGet.GetUsername())
}

func TestGetUserByID_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	respGet, err := st.AuthClient.GetUser(ctx, &ssov1.GetUserRequest{
		UserId: emptyValue,
	})
	require.Error(t, err)
	assert.Empty(t, respGet.GetUserId())

	respGet, err = st.AuthClient.GetUser(ctx, &ssov1.GetUserRequest{
		UserId: notExists,
	})
	require.Error(t, err)
	assert.Empty(t, respGet.GetUserId())
}
