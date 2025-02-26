package delta

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestEmptyDelta(t *testing.T) {
	n := New(nil)
	if n.Ops != nil {
		t.Error("failed to create Delta with nil ops")
	}
}
func TestNoopOps(t *testing.T) {
	n := New(nil)
	n.Insert("", nil).Delete(0).Retain(0, nil)
	if n.Ops != nil {
		t.Error("failed to create Delta with nil ops")
	}
}
func TestInsert1(t *testing.T) {
	n := New(nil)
	n.Insert("test", nil)
	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with insert, got: %+v\n", n.Ops)
	}

	if string(n.Ops[0].Insert) != "test" {
		t.Error("failed to create Delta with test insert, got: ", string(n.Ops[0].Insert))
	}
	if n.Ops[0].Attributes != nil {
		t.Error("failed to create Delta with only insert but no att, att was: ", n.Ops[0].Attributes)
	}
}
func TestInsertWithAttr(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Insert("test", attr)
	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with insert, got: %+v\n", n.Ops)
	}
	if string(n.Ops[0].Insert) != "test" {
		t.Error("failed to create Delta with test insert, got: ", string(n.Ops[0].Insert))
	}
	if n.Ops[0].Attributes == nil {
		t.Error("failed to create Delta with only insert with att, attr was nil")
	}
	if n.Ops[0].Attributes["bold"] != true {
		t.Errorf("failed to create Delta with only insert with att, attr was: %+v\n", n.Ops[0].Attributes)
	}
}

func TestInsertAfterDelete(t *testing.T) {
	n := New(nil)
	n.Delete(1).Insert("a", nil)
	exp := New(nil)
	exp.Insert("a", nil).Delete(1)

	if len(n.Ops) != 2 {
		t.Errorf("failed to create Delta with delete and insert, got: %+v\n", n.Ops)
	}
	if string(n.Ops[0].Insert) != string(exp.Ops[0].Insert) {
		t.Errorf("n.Ops and exp.Ops are not equal.\nn: %+v\nexp: %+v\n", n.Ops, exp.Ops)
	}
}

func TestInsertAfterDeleteWithMerge(t *testing.T) {
	n := New(nil)
	n.Insert("a", nil).Delete(1).Insert("b", nil)
	exp := New(nil)
	exp.Insert("ab", nil).Delete(1)

	if len(n.Ops) != 2 {
		t.Errorf("failed to create Delta with delete and insert merge, got: %+v\n", n.Ops)
	}
	if string(n.Ops[0].Insert) != string(exp.Ops[0].Insert) {
		t.Logf("n.Ops and exp.Ops are not equal.\nn: %+v\nexp: %+v\n", n.Ops, exp.Ops)
		t.Errorf("n.Ops and exp.Ops are not equal.\nn: %+v\n", string(n.Ops[0].Insert))
		t.Errorf("n.Ops and exp.Ops are not equal.\nn: %+v\n", string(n.Ops[1].Insert))
	}
}

func TestDelete(t *testing.T) {
	n := New(nil)
	n.Delete(0)

	if len(n.Ops) != 0 {
		t.Errorf("failed to create Delta with delete(0), got: %+v\n", n.Ops)
	}
}
func TestDeletePositive(t *testing.T) {
	n := New(nil)
	n.Delete(10)

	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with delete(10), got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Delete != 10 {
		t.Errorf("failed to create Delta with delete(10), got: %+v\n", n.Ops)
	}
}

func TestRetain(t *testing.T) {
	n := New(nil)
	n.Retain(0, nil)

	if len(n.Ops) != 0 {
		t.Errorf("failed to create Delta with retain(0), got: %+v\n", n.Ops)
	}
}

func TestRetainPositive(t *testing.T) {
	n := New(nil)
	n.Retain(2, nil)

	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with retain(2), got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Retain != 2 {
		t.Errorf("failed to create Delta with retain(2), got: %+v\n", n.Ops)
	}
}
func TestRetainPositiveAndAttrs(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Retain(2, attr)

	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with retain(2, {bold: true}), got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Retain != 2 {
		t.Errorf("failed to create Delta with retain(2, {bold: true}), got: %+v\n", n.Ops)
	}
	if n.Ops[0].Attributes["bold"] != true {
		t.Errorf("failed to create Delta with retain(2, {bold: true}), got: %+v\n", n.Ops)
	}
}

func TestPush(t *testing.T) {
	n := New(nil)
	x := []rune("test")
	n.Push(Op{Insert: x})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if string(n.Ops[0].Insert) != string(x) {
		t.Errorf("failed to Push to Delta, got: %+v\n", n.Ops)
	}
}
func TestPushMultiDelete(t *testing.T) {
	n := New(nil)
	n.Delete(2)
	x := 1
	n.Push(Op{Delete: &x})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Delete != 3 {
		t.Errorf("failed to Push to Delta, got: %+v\n", *n.Ops[0].Delete)
	}
}

func TestPushRetainDeleteInsert(t *testing.T) {
	a := New(nil).Retain(3, nil).Delete(2).Insert("X", nil)
	b := New(nil).Retain(3, nil).Insert("X", nil).Delete(2)

	if len(a.Ops) != 3 {
		t.Errorf("expcted 3 ops, got: %+v\n", a.Ops)
	}
	if *a.Ops[2].Delete != 2 {
		t.Errorf("expected delete action, got %+v\n", a.Ops)
	}

	if len(b.Ops) != 3 {
		t.Errorf("expcted 3 ops, got: %+v\n", b.Ops)
	}
	if *b.Ops[2].Delete != 2 {
		t.Errorf("expected delete action, got %+v\n", b.Ops)
	}

}

func TestPushMultiInsert(t *testing.T) {
	n := New(nil)
	n.Insert("Diego ", nil)
	x := []rune("Smith")
	n.Push(Op{Insert: x})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if string(n.Ops[0].Insert) != "Diego Smith" {
		t.Errorf("failed to Push to Delta, got: %+v\n", string(n.Ops[0].Insert))
	}
}
func TestPushMultiInsertMathingAttrs(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Insert("Diego ", attr)
	x := []rune("Smith")
	n.Push(Op{Insert: x, Attributes: attr})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if string(n.Ops[0].Insert) != "Diego Smith" {
		t.Errorf("failed to Push to Delta, got: %+v\n", string(n.Ops[0].Insert))
	}
}

func TestPushMultiInsertNonMathingAttrs(t *testing.T) {
	n := New(nil)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	n.Insert("Diego ", attr1)
	x := []rune("Smith")
	attr2 := make(map[string]interface{})
	attr2["bold"] = false
	n.Push(Op{Insert: x, Attributes: attr2})
	if len(n.Ops) != 2 {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", n.Ops)
	}
	if string(n.Ops[0].Insert) != "Diego " {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", string(n.Ops[0].Insert))
	}
	if string(n.Ops[1].Insert) != "Smith" {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", string(n.Ops[1].Insert))
	}
}

func TestPushMultiRetainMathingAttrs(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Retain(5, attr)
	x := 3
	n.Push(Op{Retain: &x, Attributes: attr})
	if len(n.Ops) != 1 {
		t.Errorf("failed to multi retain with matching attr, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Retain != 8 {
		t.Errorf("failed to multi retain with matching attr, got: %+v\n", *n.Ops[0].Retain)
	}
}

func TestPushMultiRetainNonMathingAttrs(t *testing.T) {
	n := New(nil)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	n.Retain(5, attr1)
	x := 3
	attr2 := make(map[string]interface{})
	attr2["bold"] = false

	n.Push(Op{Retain: &x, Attributes: attr2})
	if len(n.Ops) != 2 {
		t.Errorf("failed to multi retain without matching attr, expected 2 ops got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Retain != 5 {
		t.Errorf("failed to multi retain without matching attr, expected 5, got: %+v\n", *n.Ops[0].Retain)
	}
	if *n.Ops[1].Retain != 3 {
		t.Errorf("failed to multi retain without matching attr, expected 3, got: %+v\n", *n.Ops[1].Retain)
	}
}

func TestDeltaComposeTwoInsert(t *testing.T) {
	a := New(nil).Insert("A", nil)
	b := New(nil).Insert("B", nil)
	expected := New(nil).Insert("B", nil).Insert("A", nil)
	x := a.Compose(*b).Ops

	if string(x[0].Insert) != string(expected.Ops[0].Insert) {
		t.Errorf("expected %+v but got %+v\n", string(expected.Ops[0].Insert), string(x[0].Insert))
	}
}
func TestDeltaComposeInsertRetain(t *testing.T) {
	a := New(nil).Insert("A", nil)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"
	attr1["font"] = nil

	attr2 := make(map[string]interface{})
	attr2["bold"] = true
	attr2["color"] = "red"

	b := New(nil).Retain(1, attr1)
	expected := New(nil).Insert("A", attr2)
	x := a.Compose(*b).Ops

	if string(x[0].Insert) != string(expected.Ops[0].Insert) {
		t.Errorf("expected %+v but got %+v\n", string(expected.Ops[0].Insert), string(x[0].Insert))
	}

	if len(x[0].Attributes) != 2 {
		t.Errorf("expected 2 entries but got %+v\n", x[0].Attributes)
	}
}
func TestDeltaComposeInsertDelete(t *testing.T) {
	a := New(nil).Insert("A", nil)
	b := New(nil).Delete(1)

	x := a.Compose(*b).Ops

	if len(x) != 0 {
		t.Errorf("expected 0 entries but got %+v\n", x)
	}
}
func TestDeltaComposeDeleteInsert(t *testing.T) {
	a := New(nil).Insert("A", nil)
	b := New(nil).Delete(1)
	expected := New(nil).Insert("A", nil).Delete(1)
	x := b.Compose(*a).Ops

	if string(expected.Ops[0].Insert) != string(x[0].Insert) {
		t.Errorf("expected '%+v' but got '%+v'\n", string(expected.Ops[0].Insert), string(x[0].Insert))
	}
	if *expected.Ops[1].Delete != *x[1].Delete {
		t.Errorf("expected '%+v' but got '%+v'\n", *expected.Ops[0].Delete, *x[0].Delete)
	}
}
func TestDeltaComposeDeleteRetain(t *testing.T) {
	a := New(nil).Delete(1)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"

	b := New(nil).Retain(1, attr1)
	expected := New(nil).Delete(1).Retain(1, attr1)
	x := a.Compose(*b).Ops

	if *x[0].Delete != *expected.Ops[0].Delete || *x[0].Delete != 1 {
		t.Errorf("expected %+v but got %+v\n", *expected.Ops[0].Delete, *x[0].Delete)
	}

	if len(x[1].Attributes) != 2 {
		t.Errorf("expected 2 entries but got %+v\n", x[1].Attributes)
	}
}
func TestDeltaComposeDeleteDelete(t *testing.T) {
	a := New(nil).Delete(1)
	b := New(nil).Delete(1)
	expected := New(nil).Delete(2)
	x := a.Compose(*b).Ops

	if *x[0].Delete != *expected.Ops[0].Delete || *x[0].Delete != 2 {
		t.Errorf("expected %+v but got %+v\n", *expected.Ops[0].Delete, *x[0].Delete)
	}

	if len(x) != 1 {
		t.Errorf("expected 1 entries but got %+v\n", x)
	}
}

func TestDeltaComposeRetainInsert(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["color"] = "blue"
	a := New(nil).Retain(1, attr1)
	b := New(nil).Insert("B", nil)

	expected := New(nil).Insert("B", nil).Retain(1, attr1)
	x := a.Compose(*b).Ops

	if string(x[0].Insert) != string(expected.Ops[0].Insert) || string(x[0].Insert) != "B" {
		t.Errorf("expected %+v but got %+v\n", string(expected.Ops[0].Insert), string(x[0].Insert))
	}

	if len(x[1].Attributes) != 1 {
		t.Errorf("expected 1 entry but got %+v\n", x[1].Attributes)
	}
}
func TestDeltaComposeRetainRetain(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["color"] = "blue"
	a := New(nil).Retain(1, attr1)
	attr2 := make(map[string]interface{})
	attr2["color"] = "red"
	attr2["font"] = nil
	attr2["bold"] = true
	b := New(nil).Retain(1, attr2)

	attr3 := make(map[string]interface{})
	attr3["color"] = "red"
	attr3["font"] = nil
	attr3["bold"] = true
	expected := New(nil).Retain(1, attr3)
	x := a.Compose(*b).Ops

	if *x[0].Retain != *expected.Ops[0].Retain || *x[0].Retain != 1 {
		t.Errorf("expected %+v but got %+v\n", *expected.Ops[0].Retain, *x[0].Retain)
	}

	if len(x[0].Attributes) != 3 {
		t.Errorf("expected 3 entry but got %+v\n", x[0].Attributes)
	}
	if !reflect.DeepEqual(x[0].Attributes, expected.Ops[0].Attributes) {
		t.Errorf("wrong attributes, got: %+v\n", x[0].Attributes)
	}
}

func TestDeltaComposeRetainDelete(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["color"] = "blue"
	a := New(nil).Retain(1, attr1)
	b := New(nil).Delete(1)

	expected := New(nil).Delete(1)
	x := a.Compose(*b).Ops

	if *x[0].Delete != *expected.Ops[0].Delete || *x[0].Delete != 1 {
		t.Errorf("expected %+v but got %+v\n", *expected.Ops[0].Delete, *x[0].Delete)
	}

	if len(x[0].Attributes) != 0 {
		t.Errorf("expected 0 entries but got %+v\n", x[0].Attributes)
	}
	if !reflect.DeepEqual(x[0].Attributes, expected.Ops[0].Attributes) {
		t.Errorf("wrong attributes, got: %+v\n", x[0].Attributes)
	}
}

func TestDeltaComposeInsertInMiddle(t *testing.T) {
	a := New(nil).Insert("Hello", nil)
	b := New(nil).Retain(3, nil).Insert("X", nil)
	x := a.Compose(*b).Ops
	if string(x[0].Insert) != "HelXlo" {
		t.Errorf("expected 'HelXlo' but got %+v\n", string(x[0].Insert))
	}
}
func TestDeltaComposeInsertDeleteOrder(t *testing.T) {
	a := New(nil).Insert("Hello", nil)
	b := New(nil).Insert("Hello", nil)

	insertFirst := New(nil).Retain(3, nil).Insert("X", nil).Delete(1)
	deleteFirst := New(nil).Retain(3, nil).Delete(1).Insert("X", nil)

	xa := a.Compose(*insertFirst).Ops
	if string(xa[0].Insert) != "HelXo" {
		t.Errorf("expected 'HelXo' but got '%+v'\n", string(xa[0].Insert))
	}
	xb := b.Compose(*deleteFirst).Ops
	if string(xb[0].Insert) != "HelXo" {
		t.Errorf("expected 'HelXo' but got '%+v'\n", string(xb[0].Insert))
	}
}
func TestDeltaComposeDeleteEntireText(t *testing.T) {
	a := New(nil).Retain(4, nil).Insert("Hello", nil)
	b := New(nil).Delete(9)

	x := a.Compose(*b).Ops
	if *x[0].Delete != 4 {
		t.Errorf("expected '4' but got '%+v'\n", *x[0].Delete)
	}

}
func TestDeltaComposeRetainExtra(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr2 := make(map[string]interface{})
	attr2["bold"] = nil
	a := New(nil).Insert("A", attr1)
	b := New(nil).Retain(1, attr2)

	x := a.Compose(*b).Ops
	if string(x[0].Insert) != "A" {
		t.Errorf("expected 'A' but got '%+v'\n", string(x[0].Insert))
	}
	if x[0].Attributes != nil {
		t.Errorf("expected 'nil' attr but got '%+v'\n", x[0].Attributes)
	}
}
func TestDeltaComposeImmutability(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr2 := make(map[string]interface{})
	attr2["color"] = "red"
	a1 := New(nil).Insert("Test", attr1)
	b1 := New(nil).Retain(1, attr2).Delete(2)

	attr3 := make(map[string]interface{})
	attr3["color"] = "red"
	attr3["bold"] = true

	x := a1.Compose(*b1).Ops
	if string(x[0].Insert) != "T" {
		t.Errorf("expected 'T' but got '%+v'\n", string(x[0].Insert))
	}
	if string(x[1].Insert) != "t" {
		t.Errorf("expected 't' but got '%+v'\n", string(x[1].Insert))
	}
	if x[0].Attributes["color"] != "red" {
		t.Errorf("expected 'color: red' attr but got '%+v'\n", x[0].Attributes)
	}
	if x[0].Attributes["bold"] != true {
		t.Errorf("expected 'bold: true' attr but got '%+v'\n", x[0].Attributes)
	}
	if len(x[0].Attributes) != 2 {
		t.Errorf("expected '2' attrs but got '%+v'\n", x[0].Attributes)
	}
}

func TestChop(t *testing.T) {
	x := New(nil).Insert("a", nil).Insert("b", nil).Insert("c", nil).Retain(1, nil)
	ret := x.Chop()
	if len(ret.Ops) != 1 {
		t.Errorf("expected 1 op but got %+v\n", ret.Ops)
	}
}
func TestChopRetainRetain(t *testing.T) {
	x := New(nil).Retain(1, nil).Retain(1, nil)
	ret := x.Chop()
	if len(ret.Ops) != 0 {
		t.Errorf("expected '0' ops but got %+v\n", ret.Ops)
	}
}

func TestConcatEmptyDelta(t *testing.T) {
	a := New(nil).Insert("Test", nil)
	b := New(nil)
	ret := a.Concat(*b)
	if len(ret.Ops) != 1 {
		t.Errorf("expected 1 op but got %+v\n", ret.Ops)
	}
	if string(ret.Ops[0].Insert) != "Test" {
		t.Errorf("expected Insert op but got %+v\n", ret.Ops)
	}
}
func TestConcatUnmergable(t *testing.T) {
	a := New(nil).Insert("Test", nil)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	c := New(nil).Insert("!", attr1)

	ret := a.Concat(*c)
	if len(ret.Ops) != 2 {
		t.Errorf("expected 2 ops but got %+v\n", ret.Ops)
	}

	if string(ret.Ops[0].Insert) != "Test" {
		t.Errorf("expected Insert op but got %+v\n", ret.Ops)
	}
	if string(ret.Ops[1].Insert) != "!" {
		t.Errorf("expected Insert op but got:\n%+v\n", ret.Ops)
	}
	if ret.Ops[1].Attributes["bold"] != true {
		t.Errorf("expeted 'bold': true but got: \n%+v\n", ret.Ops[1].Attributes)
	}
}
func TestConcatMergable(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	a := New(nil).Insert("Test", attr1)
	c := New(nil).Insert("!", attr1)

	ret := a.Concat(*c)
	if len(ret.Ops) != 1 {
		t.Errorf("expected 1 op but got %+v\n", ret.Ops)
	}

	if string(ret.Ops[0].Insert) != "Test!" {
		t.Errorf("expected Insert op but got %+v\n", ret.Ops)
	}
	if ret.Ops[0].Attributes["bold"] != true {
		t.Errorf("expeted 'bold': true but got: \n%+v\n", ret.Ops[0].Attributes)
	}
}

func TestTransformPositionInsertBeforePos(t *testing.T) {
	delta := New(nil).Insert("A", nil)
	if x := delta.TransformPosition(2, false); x != 3 {
		t.Error("expected 3 but got ", x)
	}
}
func TestTransformPositionInsertAfterPos(t *testing.T) {
	delta := New(nil).Retain(2, nil).Insert("A", nil)
	if x := delta.TransformPosition(1, false); x != 1 {
		t.Error("expected 1 but got ", x)
	}
}
func TestTransformPositionInsertAtPos(t *testing.T) {
	delta := New(nil).Retain(2, nil).Insert("A", nil)
	if x := delta.TransformPosition(2, true); x != 2 {
		t.Error("expected 2 but got ", x)
	}
	if x := delta.TransformPosition(2, false); x != 3 {
		t.Error("expected 3 but got ", x)
	}
}
func TestTransformPositionDeleteBeforePos(t *testing.T) {
	delta := New(nil).Delete(2)
	if x := delta.TransformPosition(4, false); x != 2 {
		t.Error("expected 2 but got ", x)
	}
}
func TestTransformPositionDeleteAfterPos(t *testing.T) {
	delta := New(nil).Retain(4, nil).Delete(2)
	if x := delta.TransformPosition(2, false); x != 2 {
		t.Error("expected 2 but got ", x)
	}
}
func TestTransformPositionDeleteAcrossPos(t *testing.T) {
	delta := New(nil).Retain(1, nil).Delete(4)
	if x := delta.TransformPosition(2, false); x != 1 {
		t.Error("expected 1 but got ", x)
	}
}
func TestTransformPositionInsertAndDeleteBeforePos(t *testing.T) {
	delta := New(nil).Retain(2, nil).Insert("A", nil).Delete(2)
	if x := delta.TransformPosition(4, false); x != 3 {
		t.Error("expected 3 but got ", x)
	}
}
func TestTransformPositionInsertAndDeleteAcrossPos(t *testing.T) {
	delta := New(nil).Retain(2, nil).Insert("A", nil).Delete(4)
	if x := delta.TransformPosition(4, false); x != 3 {
		t.Error("expected 3 but got ", x)
	}
}
func TestTransformPositionDeleteBeforeAndDeleteAcrossPos(t *testing.T) {
	delta := New(nil).Delete(1).Retain(1, nil).Delete(4)
	if x := delta.TransformPosition(4, false); x != 1 {
		t.Error("expected 1 but got ", x)
	}
}

func TestTransformInsertInsert(t *testing.T) {
	a1 := New(nil).Insert("A", nil)
	b1 := New(nil).Insert("B", nil)
	a2 := New(a1.Ops)
	b2 := New(b1.Ops)
	x1 := a1.Transform(*b1, true)
	if *x1.Ops[0].Retain != 1 {
		t.Errorf("expected 1 but got %+v\n", x1)
	}
	if string(x1.Ops[1].Insert) != "B" {
		t.Errorf("expected 'B' but got %+v\n", x1)
	}

	x2 := a2.Transform(*b2, false)
	if len(x2.Ops) != 1 {
		t.Errorf("expected 1 op but got %+v\n", x2)
	}
	if string(x2.Ops[0].Insert) != "B" {
		t.Errorf("expected 'B' but got %+v\n", x2)
	}
}
func TestTransformInsertInsertChinese(t *testing.T) {
	a1 := New(nil).Insert("你", nil)
	b1 := New(nil).Insert("好", nil)
	a2 := New(a1.Ops)
	b2 := New(b1.Ops)
	x1 := a1.Transform(*b1, true)
	if *x1.Ops[0].Retain != 1 {
		t.Errorf("expected 1 but got %+v\n", x1)
	}
	if string(x1.Ops[1].Insert) != "好" {
		t.Errorf("expected '好' but got %+v\n", x1)
	}

	x2 := a2.Transform(*b2, false)
	if len(x2.Ops) != 1 {
		t.Errorf("expected 1 op but got %+v\n", x2)
	}
	if string(x2.Ops[0].Insert) != "好" {
		t.Errorf("expected '好' but got %v\n", x2)
	}
}
func TestTransformInsertRetain(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"
	a := New(nil).Insert("A", nil)
	b := New(nil).Retain(1, attr1)
	x := a.Transform(*b, true)
	if *x.Ops[0].Retain != 1 {
		t.Errorf("expected '1' but got %+v\n", x)
	}
	if *x.Ops[1].Retain != 1 {
		t.Errorf("expected '1' but got %+v\n", x)
	}
	if len(x.Ops[1].Attributes) != 2 {
		t.Errorf("expected '2' but got %+v\n", x)
	}
	if x.Ops[1].Attributes["color"] != "red" {
		t.Errorf("expected 'red' but got %+v\n", x)
	}
	if x.Ops[1].Attributes["bold"] != true {
		t.Errorf("expected 'true' but got %+v\n", x)
	}
}

func TestTransformInsertDelete(t *testing.T) {
	a := New(nil).Insert("A", nil)
	b := New(nil).Delete(1)
	x := a.Transform(*b, true)
	if *x.Ops[0].Retain != 1 {
		t.Errorf("expected '1' but got %+v\n", x)
	}
	if *x.Ops[1].Delete != 1 {
		t.Errorf("expected '1' but got %+v\n", x)
	}
}
func TestTransformDeleteInsert(t *testing.T) {
	a := New(nil).Delete(1)
	b := New(nil).Insert("B", nil)
	x := a.Transform(*b, true)
	if string(x.Ops[0].Insert) != "B" {
		t.Errorf("expected 'B' but got %+v\n", x)
	}
	if len(x.Ops) > 1 {
		t.Errorf("expected '1' op but got %+v\n", x)
	}
}
func TestTransformDeleteRetain(t *testing.T) {
	a := New(nil).Delete(1)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"

	b := New(nil).Retain(1, attr1)
	x := a.Transform(*b, true)
	if len(x.Ops) > 0 {
		t.Errorf("expected '0' op but got %+v\n", x)
	}
}

func TestTransformDeleteDelete(t *testing.T) {
	a := New(nil).Delete(1)
	b := New(nil).Delete(1)

	x := a.Transform(*b, false)
	if len(x.Ops) > 0 {
		t.Errorf("expected '0' op but got %+v\n", x)
	}
}
func TestTransformRetainInsert(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["color"] = "blue"

	a := New(nil).Retain(1, attr1)
	b := New(nil).Insert("B", nil)

	x := a.Transform(*b, false)
	if len(x.Ops) > 1 {
		t.Errorf("expected '1' op but got %+v\n", x)
	}
	if string(x.Ops[0].Insert) != "B" {
		t.Errorf("expected 'B' op but got %+v\n", x)
	}
}
func TestTransformRetainRetain(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["color"] = "blue"
	attr2 := make(map[string]interface{})
	attr2["color"] = "red"
	attr2["bold"] = true

	a1 := New(nil).Retain(1, attr1)
	b1 := New(nil).Retain(1, attr2)
	a2 := New(nil).Retain(1, attr1)
	b2 := New(nil).Retain(1, attr2)

	x := a1.Transform(*b1, true)
	if len(x.Ops) > 1 {
		t.Errorf("expected '1' op but got %+v\n", x)
	}
	if *x.Ops[0].Retain != 1 {
		t.Errorf("expected 'B' op but got %+v\n", x)
	}
	if len(x.Ops[0].Attributes) > 1 {
		t.Errorf("expected '1' attr but got %+v\n", x)
	}
	x = b2.Transform(*a2, true)
	if len(x.Ops) > 0 {
		t.Errorf("expected '0' op but got %+v\n", x)
	}
}
func TestTransformRetainRetainNoPrio(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["color"] = "blue"
	attr2 := make(map[string]interface{})
	attr2["color"] = "red"
	attr2["bold"] = true

	a1 := New(nil).Retain(1, attr1)
	b1 := New(nil).Retain(1, attr2)
	a2 := New(nil).Retain(1, attr1)
	b2 := New(nil).Retain(1, attr2)

	x := a1.Transform(*b1, false)
	if len(x.Ops) > 1 {
		t.Errorf("expected '1' op but got %+v\n", x)
	}
	if *x.Ops[0].Retain != 1 {
		t.Errorf("expected 'B' op but got %+v\n", x)
	}
	if len(x.Ops[0].Attributes) != 2 {
		t.Errorf("expected '2' attr but got %+v\n", x)
	}
	x = b2.Transform(*a2, false)
	if len(x.Ops) != 1 {
		t.Errorf("expected '1' op but got %+v\n", x)
	}
}
func TestTransformAlternatingEdits(t *testing.T) {

	a := New(nil).Retain(2, nil).Insert("si", nil).Delete(5)
	b := New(nil).Retain(1, nil).Insert("e", nil).Delete(5).Retain(1, nil).Insert("ow", nil)

	x := a.Transform(*b, false)
	if len(x.Ops) != 5 {
		t.Errorf("expected '5' op but got %+v\n", x)
	}
	if *x.Ops[0].Retain != 1 {
		t.Errorf("expected 'retain 1'  but got %+v\n", x)
	}
	if string(x.Ops[1].Insert) != "e" {
		t.Errorf("expected 'insert e'  but got %+v\n", x)
	}
	if *x.Ops[2].Delete != 1 {
		t.Errorf("expected 'delete 1'  but got %+v\n", x)
	}
	if *x.Ops[3].Retain != 2 {
		t.Errorf("expected 'retain 2'  but got %+v\n", x)
	}
	if string(x.Ops[4].Insert) != "ow" {
		t.Errorf("expected 'insert ow'  but got %+v\n", x)
	}

	x = b.Transform(*a, false)
	if len(x.Ops) != 3 {
		t.Errorf("expected '3' op but got %+v\n", x)
	}
	if *x.Ops[0].Retain != 2 {
		t.Errorf("expected 'retain 2'  but got %+v\n", x)
	}
	if string(x.Ops[1].Insert) != "si" {
		t.Errorf("expected 'insert si'  but got %+v\n", x)
	}
	if *x.Ops[2].Delete != 1 {
		t.Errorf("expected 'delete 1'  but got %+v\n", x)
	}
}
func TestTransformConflictingAppends(t *testing.T) {

	a := New(nil).Retain(3, nil).Insert("aa", nil)
	b := New(nil).Retain(3, nil).Insert("bb", nil)

	x := a.Transform(*b, true)
	if len(x.Ops) != 2 {
		t.Errorf("expected '2' ops but got %+v\n", x)
	}
	if *x.Ops[0].Retain != 5 {
		t.Errorf("expected 'retain 5'  but got %+v\n", x)
	}
	if string(x.Ops[1].Insert) != "bb" {
		t.Errorf("expected 'insert bb'  but got %+v\n", x)
	}
	x = b.Transform(*a, false)
	if len(x.Ops) != 2 {
		t.Errorf("expected '2' ops but got %+v\n", x)
	}
	if *x.Ops[0].Retain != 3 {
		t.Errorf("expected 'retain 3'  but got %+v\n", x)
	}
	if string(x.Ops[1].Insert) != "aa" {
		t.Errorf("expected 'insert aa'  but got %+v\n", x)
	}
}
func TestTransformPrependAppend(t *testing.T) {

	a := New(nil).Insert("aa", nil)
	b := New(nil).Retain(3, nil).Insert("bb", nil)

	x := a.Transform(*b, false)
	if len(x.Ops) != 2 {
		t.Errorf("expected '2' ops but got %+v\n", x)
	}
	if *x.Ops[0].Retain != 5 {
		t.Errorf("expected 'retain 5' but got %+v\n", x)
	}
	if string(x.Ops[1].Insert) != "bb" {
		t.Errorf("expected 'insert bb' but got %+v\n", x)
	}
	x = b.Transform(*a, false)
	if len(x.Ops) != 1 {
		t.Errorf("expected '1' op but got %+v\n", x)
	}
	if string(x.Ops[0].Insert) != "aa" {
		t.Errorf("expected 'insert aa' but got %+v\n", x)
	}
}
func TestTransformTrailingDeletesWithDifferingLengths(t *testing.T) {

	a := New(nil).Retain(2, nil).Delete(1)
	b := New(nil).Delete(3)

	x := a.Transform(*b, false)
	if len(x.Ops) != 1 {
		t.Errorf("expected '1' op but got %+v\n", x)
	}
	if *x.Ops[0].Delete != 2 {
		t.Errorf("expected 'delete 2' but got %+v\n", x)
	}
	x = b.Transform(*a, false)
	if len(x.Ops) != 0 {
		t.Errorf("expected '0' ops but got %+v\n", x)
	}
}

func BenchmarkInsert(t *testing.B) {
	a := New(nil)
	for x := 0; x < t.N; x++ {
		a.Insert("abcdefghijkl", nil)
	}
}
func BenchmarkInsert2(t *testing.B) {
	a := New(nil)
	for x := 0; x < t.N; x++ {
		a.Insert("a", nil)
	}
}

func BenchmarkInsert3(t *testing.B) {
	a := New(nil)
	for x := 0; x < t.N; x++ {
		a.Insert("a", nil).Retain(1, nil).Insert("b", nil).Retain(4, nil).Delete(1)
	}
}

func TestFromJSON1(t *testing.T) {
	in := []byte(`{"ops":[{"retain":54},{"insert":"g"}]}`)
	d, err := FromJSON(in)
	if err != nil {
		t.Error("failed with ", err)
	}
	if *d.Ops[0].Retain != 54 {
		t.Errorf("expected 'retain 54' but got %+v\n", *d.Ops[0].Retain)
	}
	if string(d.Ops[1].Insert) != "g" {
		t.Errorf("expected 'insert g' but got %+v\n", string(d.Ops[1].Insert))
	}
}
func TestFromJSON2(t *testing.T) {
	in := []byte(`{"ops":[{"retain":35},{"retain":11,"attributes":{"bold":true}}]}`)
	d, err := FromJSON(in)
	if err != nil {
		t.Error("failed with ", err)
	}
	if *d.Ops[0].Retain != 35 {
		t.Errorf("expected 'retain 35' but got %+v\n", *d.Ops[0].Retain)
	}
	if *d.Ops[1].Retain != 11 {
		t.Errorf("expected 'retain 11' but got %+v\n", *d.Ops[1].Retain)
	}
	if d.Ops[1].Attributes["bold"] != true {
		t.Errorf("expected 'bold: true' but got %+v\n", d.Ops[1].Attributes)
	}
}
func TestFromJSON3(t *testing.T) {
	in := []byte(`{"ops":[{"retain":8},{"delete":7}]}`)
	d, err := FromJSON(in)
	if err != nil {
		t.Error("failed with ", err)
	}
	if *d.Ops[0].Retain != 8 {
		t.Errorf("expected 'retain 8' but got %+v\n", *d.Ops[0].Retain)
	}
	if *d.Ops[1].Delete != 7 {
		t.Errorf("expected 'delete 7' but got %+v\n", *d.Ops[1].Delete)
	}
}
func TestMarshalJSON1(t *testing.T) {
	in := []byte(`{"ops":[{"retain":8},{"delete":7}]}`)
	d, err := FromJSON(in)
	if err != nil {
		t.Error("failed with ", err)
	}
	out, err := json.Marshal(d)
	if err != nil {
		t.Error("failed to get json string, err: ", err)
	}
	if bytes.Compare(in, out) != 0 {
		t.Errorf("expected:\n'%+v' but got :\n%+v\n", string(in[:]), string(out[:]))
	}
}
func TestMarshalJSON2(t *testing.T) {
	in := []byte(`{"ops":[{"retain":35},{"retain":11,"attributes":{"bold":true}}]}`)
	d, err := FromJSON(in)
	if err != nil {
		t.Error("failed with ", err)
	}
	out, err := json.Marshal(d)
	if err != nil {
		t.Error("failed to get json string, err: ", err)
	}
	if bytes.Compare(in, out) != 0 {
		t.Errorf("expected:\n'%+v' but got :\n%+v\n", string(in[:]), string(out[:]))
	}
}
func TestMarshalJSON3(t *testing.T) {
	in := []byte(`{"ops":[{"insert":"diego"},{"retain":11,"attributes":{"bold":true}}]}`)
	d, err := FromJSON(in)
	if err != nil {
		t.Error("failed with ", err)
	}
	out, err := json.Marshal(d)
	if err != nil {
		t.Error("failed to get json string, err: ", err)
	}
	if bytes.Compare(in, out) != 0 {
		t.Errorf("expected:\n'%+v' but got :\n'%+v'\n", string(in[:]), string(out[:]))
	}
}
func TestMarshalJSON4(t *testing.T) {
	in := []byte(`{"ops":[{"insert":"你好"},{"retain":2,"attributes":{"bold":true}}]}`)
	d, err := FromJSON(in)
	if err != nil {
		t.Error("failed with ", err)
	}
	out, err := json.Marshal(d)
	if err != nil {
		t.Error("failed to get json string, err: ", err)
	}
	if bytes.Compare(in, out) != 0 {
		t.Errorf("expected:\n'%+v' but got :\n'%+v'\n", string(in[:]), string(out[:]))
	}
}

func BenchmarkFromJson1(t *testing.B) {
	in := []byte(`{"ops":[{"retain":35},{"retain":11,"attributes":{"bold":true}}]}`)
	var delta *Delta
	for x := 0; x < t.N; x++ {
		delta, _ = FromJSON(in)
	}
	t.StopTimer()
	t.Logf("delta is %+v\n", delta)
}

func TestDeleteChineseAll(t *testing.T) {
	hellowWorld := "你好，世界!"
	a := New(nil).Insert(hellowWorld, nil)
	b := New(nil).Delete(6)

	ret := a.Compose(*b)

	if ret.Ops != nil {
		t.Errorf("ret is %+v\n", ret)
	}
}
func TestDeleteChinese1Char(t *testing.T) {
	hellowWorld := "你好，世界!"
	expected := "你好，世界"
	a := New(nil).Insert(hellowWorld, nil)
	b := New(nil).Retain(5, nil).Delete(1)
	ret := a.Compose(*b)
	if string(ret.Ops[0].Insert) != expected {
		t.Errorf("Expected: '%s' but got %+v\n", expected, string(ret.Ops[0].Insert))
	}
}

func TestDeleteChinese2Char(t *testing.T) {
	hellowWorld := "你好，世界!"
	expected := "你好，世"
	a := New(nil).Insert(hellowWorld, nil)
	b := New(nil).Retain(4, nil).Delete(2)
	ret := a.Compose(*b)
	if string(ret.Ops[0].Insert) != expected {
		t.Errorf("Expected: '%s' but got %+v\n", expected, string(ret.Ops[0].Insert))
	}
}
