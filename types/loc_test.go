package types

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
)

func TestLocCtor(t* testing.T) {
	result := Loc{ "3.6.9", 3, 6, 9, "new" } // TODO: remove the hard coded "new"
	assert.IsType(t, Loc{}, result )
	assert.True( t,
		result.ID == "3.6.9" && result.X == 3 && result.Y == 6 && result.Z == 9,
		"Expected a basic struct to construct properly" )
	assert.True( t,
		result.Status == "new",
		"Expected a new loc to default to Status = new" )
}

func TestLocFromStringPositive(t* testing.T) {
	result, err := LocFromString( "3.6.9" )
	require.NoError( t, err, "Positive test should not throw an error" )
	assert.IsType(t, Loc{}, result, "Should return a Loc struct" )
	assert.True( t,
		result.ID == "3.6.9" && result.X == 3 && result.Y == 6 && result.Z == 9,
		"X,Y,Z values not mapped as expected" )
}

func TestLocFromStringBadDelimiters(t* testing.T) {
	_, err := LocFromString( "3.6*9" )
	assert.Error( t, err, "Negative test should throw an error" )
	assert.True( t, strings.Contains( err.Error(), "x.y.z"), "Expect an error message that mentions the proper format" )
}

func TestLocConvertPositive(t *testing.T) {
	testInput := "25.0.-54"
	x,y,z,err := LocConvert(testInput)
	require.NoError( t, err, "Should not get error in positive test case")
	assert.Equal(t,  25, x, "x should be int of value 25. Instead got type %T of value %v", x, x)
	assert.Equal(t,   0, y, "y should be int of value 0. Instead got type %T of value %v", y, y)
	assert.Equal(t, -54, z, "z should be int of value -54. Instead got type %T of value %v", z, z)
}

func TestLocConvertBadDelimiters(t *testing.T) {
	testInput := "13,14,15"
	_,_,_,err := LocConvert(testInput)
	require.Error( t, err, "Bad delimeters should throw an error")
	assert.True( t, strings.Contains( err.Error(), "x.y.z"), "Expect an error message that mentions the proper format" )
}

func TestLocConvertTooFewCoords(t * testing.T) {
	testInput := "33.105"
	_,_,_,err := LocConvert(testInput)
	require.Error( t, err, "Not enough coords -- must be 3 -- should throw an error")
	assert.True( t, strings.Contains( err.Error(), "x.y.z"), "Expect an error message that mentions the proper format" )
}

func TestLocConvertTooManyCoords(t * testing.T) {
	testInput := "19.45.46.75"
	_,_,_,err := LocConvert(testInput)
	require.Error( t, err, "Too many coords -- must be 3 -- should throw an error")
	assert.True( t, strings.Contains( err.Error(), "x.y.z"), "Expect an error message that mentions the proper format" )
}

func TestLocConvertNonIntCoords(t * testing.T) {
	testInput := "22.five.13"
	_,_,_,err := LocConvert(testInput)
	require.Error( t, err, "Non-int based coords should throw an error")
	assert.True( t, strings.Contains( err.Error(), "integer"), "Expect an error message that mentions the proper format" )
}

func TestLocFromCoordsID( t *testing.T) {
	loc, err := LocFromCoords( 12, 21, 0 )
	require.NoError( t, err, "Theoretically any valid ints should convert? Err: %s", err )
	assert.Equal( t, "12.21.0", loc.ID )
	loc, err = LocFromCoords( -13, 19, -27 )
	require.NoError( t, err, "Theoretically any valid ints should convert? Err: %s", err )
	assert.Equal( t, "-13.19.-27", loc.ID )
}

func TestJSONForm(t *testing.T) {
    t.Run("Positive", func(t *testing.T) {
		loc, err := LocFromString( "9.6.3" )
		if err != nil {
			t.Fatalf( "Error from Loc creation: %s", err )
		}
		expected := []byte( `{"id":"9.6.3","x":9,"y":6,"z":3,"status":"new"}` )
		assert.Equal( t, expected, loc.JSONForm() )
	})

	 t.Run("Positive", func(t *testing.T) {
		testJSON := []byte( `{"id":"19.6.13","x":19,"y":6,"z":13,"status":"new"}` )
		expected, _ := LocFromCoords( 19, 6, 13 )
		actual, err := LocFromJSON( testJSON )
		require.NoError(t, err, "Didn't want to see an error here")
		assert.Equal(t, expected, actual )
	})
	
	 t.Run("ExtraElements", func(t *testing.T) {
		testJSON := []byte( `{"id":"19.16.23","x":19,"y":16,"z":23,"status":"new","extra": 1003}` )
		expected, _ := LocFromCoords( 19, 16, 23 )
		actual, err := LocFromJSON( testJSON )
		require.NoError(t, err, "An 'extra' element should be ignored")
		assert.Equal(t, expected, actual )
	})

	 t.Run("WrongType", func(t *testing.T) {
		testJSON := []byte( `{"id":"19.16.23","x":19,"y":"16","z":23,"status":"new"}` )
		_, err := LocFromJSON( testJSON )
		require.Error(t, err, "Should complain about mismatched types")
	})

	 t.Run("MissingElemen", func(t *testing.T) {
		testJSON := []byte( `{"id":"19.6.23","x":19,"z":23,"status":"new"}` )
		_, err := LocFromJSON( testJSON )
		require.Error(t, err, "Should complain about missing y element")
		require.Contains(t, err.Error(), "missing", "Not the text we were looking for")
	})
}

// Helper function to swallow the multiple return value. Allow a newLoc call within a struct declaration.
func newLoc( x int, y int, z int ) Loc {
	result, _ := LocFromCoords( x, y, z )
	return result
}

// Table driven test to try a wide range of positive cases
func TestDistanceFrom(t *testing.T) {
	var cases = []struct {
		origin Loc
		target Loc
		expected int
	} {
		{ newLoc(  12, -7, 99 ), newLoc(  19, 10, 99 ), 17 },
		{ newLoc( 100, -7,  0 ), newLoc( 113, 10, 99 ), 99 },
		{ newLoc(   1,  2,  3 ), newLoc( -44,  2, -3 ), 45 },
		{ newLoc(   0,  0,  0 ), newLoc(   0,  0,  0 ),  0 },
	}

	for num, c := range cases {
		t.Run( fmt.Sprintf( "case#%d", num ), func( t *testing.T ) {
			//fmt.Printf( "-- case: %v.DistanceFrom( %v ) = %d\n", c.origin, c.target, c.expected )
			actual := c.origin.DistanceFrom( c.target )
			assert.Equal( t, c.expected, actual ) 
		} )
	}
}

/* func TestFindNeighbors(t *testing.T) {
	t.Run("PositiveFromOrigin", func(t *testing.T){
		t.Error(t, "Not implemented")
	})
} */
