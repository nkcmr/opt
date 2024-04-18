package opt

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJoin(t *testing.T) {
	t.Run("BothNone", func(t *testing.T) {
		var x, y Option[int]
		result := Join(x, y, func(x, y int) int {
			panic("should not be called")
		})
		require.True(t, result.None())
	})
	t.Run("XorNone", func(t *testing.T) {
		var y Option[int]
		x := Some(int(5))
		result := Join(x, y, func(x, y int) int {
			panic("should not be called")
		})
		require.True(t, result.None())
		result = Join(y, x, func(y, x int) int {
			panic("should not be called")
		})
		require.True(t, result.None())
	})
	t.Run("BothSome", func(t *testing.T) {
		x := Some(int(4))
		y := Some(int(8))
		result := Join(x, y, func(x, y int) int {
			require.Equal(t, int(4), x)
			require.Equal(t, int(8), y)
			return x + y
		})
		require.True(t, result.Some())
		require.Equal(t, int(12), result.Unwrap())
	})
}

func TestMaybeUnwrap(t *testing.T) {
	xOpt := Some(int(5))
	yOpt := None[int]()

	x, xok := xOpt.MaybeUnwrap()
	require.True(t, xok)
	require.Equal(t, int(5), x)

	y, yok := yOpt.MaybeUnwrap()
	require.False(t, yok)
	require.Equal(t, int(0), y)
}

func TestFromMaybe(t *testing.T) {
	x := FromMaybe(int(3), false)
	require.True(t, x.None())
	require.Equal(t, int(0), x.UnwrapOrZero())

	y := FromMaybe(int(15), true)
	require.True(t, y.Some())
	require.Equal(t, int(15), y.Unwrap())
}

func TestUnwrapOrZero(t *testing.T) {
	xOpt := Some(int(5))
	yOpt := None[int]()

	x := xOpt.UnwrapOrZero()
	require.Equal(t, int(5), x)

	y := yOpt.UnwrapOrZero()
	require.Equal(t, int(0), y)
}

func TestMap(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		var x Option[int]
		y := Map(x, func(in int) Option[int] {
			return Some(in * 2)
		})
		require.True(t, y.None())
	})
	t.Run("Some", func(t *testing.T) {
		x := Some[int](5)
		y := Map(x, func(in int) Option[int] {
			return Some(in * 2)
		})
		require.True(t, y.Some())
		require.Equal(t, int(10), y.Unwrap())
	})
}

func TestEqual(t *testing.T) {
	t.Run(`Some(1) == Some(1)`, func(t *testing.T) {
		a := Some(1)
		b := Some(1)
		require.True(t, Equal(a, b))
	})
	t.Run(`Some(1) != Some(2)`, func(t *testing.T) {
		a := Some(1)
		b := Some(2)
		require.False(t, Equal(a, b))
	})
	t.Run(`None() != Some(2)`, func(t *testing.T) {
		a := None[int]()
		b := Some(2)
		require.False(t, Equal(a, b))
	})
	t.Run(`Some(2) != None()`, func(t *testing.T) {
		a := Some(2)
		b := None[int]()
		require.False(t, Equal(a, b))
	})
	t.Run("None() == None()", func(t *testing.T) {
		a := None[int]()
		b := None[int]()
		require.True(t, Equal(a, b))
	})
}

func TestFromPointer(t *testing.T) {
	x := new(int64)
	*x = 5
	xopt := FromPointer(x)
	require.True(t, xopt.Some())
	require.Equal(t, int64(5), xopt.Unwrap())
	require.Equal(t, int64(5), xopt.UnwrapOr(1))

	var y *string
	yopt := FromPointer(y)
	require.True(t, yopt.None())
	require.Panics(t, func() {
		_ = yopt.Unwrap()
	})
}

func TestOption(t *testing.T) {
	t.Run("zero value is valid", func(t *testing.T) {
		var ov Option[int]
		require.True(t, ov.None())
		require.False(t, ov.Some())
		require.Panics(t, func() {
			_ = ov.Unwrap()
		})
		require.Equal(t, int(5), ov.UnwrapOr(5))
	})
	t.Run("normal stuff", func(t *testing.T) {
		ov := Some(5)
		require.False(t, ov.None())
		require.True(t, ov.Some())
		require.Equal(t, int(5), ov.Unwrap())
		require.Equal(t, int(5), ov.UnwrapOr(1))
	})
}

func TestCoalesce(t *testing.T) {
	a := Coalesce(
		Some(int(0)),
		Some(int(1)),
	)
	require.True(t, a.Some())
	require.False(t, a.None())
	require.Equal(t, int(0), a.Unwrap())

	b := Coalesce(
		None[int](),
		Some(int(89)),
		Some(int(90)),
	)
	require.True(t, b.Some())
	require.False(t, b.None())
	require.Equal(t, int(89), b.Unwrap())

	c := Coalesce[int]()
	require.True(t, c.None())
	require.False(t, c.Some())
	require.Panics(t, func() {
		_ = c.Unwrap()
	})

	d := Coalesce(
		None[int](),
		None[int](),
	)
	require.True(t, d.None())
	require.False(t, d.Some())
	require.Panics(t, func() {
		_ = d.Unwrap()
	})
	require.Equal(t, int(5), d.UnwrapOr(5))
}

func TestJSON(t *testing.T) {
	type TestStruct struct {
		Foo string
		Bar Option[int]
	}

	t.Run("encode", func(t *testing.T) {
		a := TestStruct{}
		adata, err := json.Marshal(a)
		require.NoError(t, err)
		require.Equal(t, `{"Foo":"","Bar":null}`, string(adata))

		b := TestStruct{
			Foo: "beep",
			Bar: Some(5),
		}
		bdata, err := json.Marshal(b)
		require.NoError(t, err)
		require.Equal(t, `{"Foo":"beep","Bar":5}`, string(bdata))
	})
	t.Run("decode", func(t *testing.T) {
		var v TestStruct
		err := json.Unmarshal([]byte(`{"Foo":"bar","Bar":null}`), &v)
		require.NoError(t, err)
		require.Equal(t, "bar", v.Foo)
		require.True(t, v.Bar.None())

		var v2 TestStruct
		err = json.Unmarshal([]byte(`{"Foo":"bap","Bar":8}`), &v2)
		require.NoError(t, err)
		require.Equal(t, "bap", v2.Foo)
		require.True(t, v2.Bar.Some())
		require.Equal(t, int(8), v2.Bar.Unwrap())
	})
}
