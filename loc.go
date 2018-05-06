package webstuff

import (
	"bytes"
	"math"
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
)

// Location reresents an interface for any object storeable in mongo with an explicit ID field
type Location interface {
	GetID() string
}

// Loc contains the coords and methods to handle a 3 axis location on a hex map
type Loc struct {
	ID     string  `json:"id" bson:"_id"`
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Z      int     `json:"z"`
	Status string  `json:"status"` // TODO: default to "new"
}

// GetID getter for ID field
func(l Loc) GetID() string {
	return l.ID
}

// LocFromCoords generates a Loc instance from x, y and z coordinates.
// Should enforce uniqueness at some point?
func LocFromCoords( x int, y int, z int ) (result Loc, err error) {
	id := fmt.Sprintf( "%d.%d.%d", x, y, z )
	result = Loc{ id, x, y, z, "new" }
	return result, err
}

// LocFromString generates a Loc instance from a string containing the coords in the format 'x.y.z'
func LocFromString(loc string) (result Loc, err error) {
	x,y,z,err := LocConvert(loc)
	if err == nil {
		result, err = LocFromCoords( x, y, z )
	} 
	return result, err
}

// LocFromJSON generates a Loc instance from JSON. Expected JSON form should match the struct declaration. Duh!
func LocFromJSON(jsonIn []byte) (Loc, error) {
	result := Loc{}
	if err := json.Unmarshal(jsonIn, &result); err != nil {
		return result, err
	}
	d := json.NewDecoder(bytes.NewBuffer(jsonIn))
	var dd map[string]interface{}
	if err:= d.Decode( &dd ); err != nil {
		return result, err
	}
	if dd["x"] == nil || dd["y"] == nil || dd["z"] == nil {
		return result, fmt.Errorf("missing one or more x,y,z elements")
	}
	return result, nil
}

// LocConvert parses a string of the format 'x.y.z' into the individual elements
func LocConvert(loc string) (x int, y int, z int, err error) {
	xyz := strings.Split(loc, ".")
	if len(xyz) != 3 {
		return x,y,z,fmt.Errorf("XYZ param must be of the format 'x.y.z'. Got: %s", loc )
		//return x,y,z,LocConvertError( fmt.Sprintf("XYZ param must be of the format 'x.y.z'. Got: %s", loc ) )
	}
	var n int64
	n, err = strconv.ParseInt(xyz[0], 10, 64)
	if err != nil {
		return x,y,z,fmt.Errorf("could not parse x value as integer. Got: %s", loc )
	}
	x = int(n)
	n, err = strconv.ParseInt(xyz[1], 10, 64)
	if err != nil {
		return x,y,z,fmt.Errorf("could not parse y value as integer. Got: %s", loc )
	}
	y = int(n)
	n, err = strconv.ParseInt(xyz[2], 10, 64)
	if err != nil {
		return x,y,z,fmt.Errorf("could not parse z value as integer. Got: %s", loc )
	}
	z = int(n)
	return x,y,z,nil
}

// StringForm provides the location in the "x.y.z" format
func (l Loc) StringForm() string {
	return fmt.Sprintf( "%d.%d.%d", l.X, l.Y, l.Z )
}

// JSONForm provides the location in JSON
func (l Loc) JSONForm() []byte {
	//fmt.Printf("Marshalling loc %v to JSON\n", l)
	j, err := json.Marshal(l)
	if err != nil {
		fmt.Println("bad things happened in JSON marshal" )
		panic(err)
	}
	return j
}

// DistanceFrom returns the distance from this Loc to the specified Loc
func (l Loc) DistanceFrom(target Loc) int {
	dx := math.Abs( float64(l.X) - float64(target.X) ) 
	dy := math.Abs( float64(l.Y) - float64(target.Y) ) 
	dz := math.Abs( float64(l.Z) - float64(target.Z) ) 
	max := int( math.Max( math.Max( dx, dy ), dz ) )
	return max
}