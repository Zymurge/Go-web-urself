package types

import (
	"testing"
	"github.com/stretchr/testify/require"
//	"github.com/stretchr/testify/suite"
)

func TestBuildPositive(t *testing.T) {
	expectedID := "-2.-2.4"
	target := Grid{}
	target.Build( 4, 4 )
	result := target.GetLoc(expectedID)
	require.Equal(t, expectedID, result.GetID())
	require.Equal(t, 25, len(target.locs))
	require.Equal(t, -2, target.XMin())
	require.Equal(t,  2, target.XMax())
	require.Equal(t, -2, target.YMin())
	require.Equal(t,  2, target.YMax())
	require.Equal(t, -4, target.ZMin())
	require.Equal(t,  4, target.ZMax())
}