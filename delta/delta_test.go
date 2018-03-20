package delta

import (
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

	if *n.Ops[0].Insert != "test" {
		t.Error("failed to create Delta with test insert, got: ", *n.Ops[0].Insert)
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
	if *n.Ops[0].Insert != "test" {
		t.Error("failed to create Delta with test insert, got: ", *n.Ops[0].Insert)
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
	if *n.Ops[0].Insert != *exp.Ops[0].Insert {
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
	if *n.Ops[0].Insert != *exp.Ops[0].Insert {
		t.Logf("n.Ops and exp.Ops are not equal.\nn: %+v\nexp: %+v\n", n.Ops, exp.Ops)
		t.Errorf("n.Ops and exp.Ops are not equal.\nn: %+v\n", *n.Ops[0].Insert)
		t.Errorf("n.Ops and exp.Ops are not equal.\nn: %+v\n", *n.Ops[1].Insert)
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
	x := "test"
	n.Push(Op{Insert: &x})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != x {
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
	x := "Smith"
	n.Push(Op{Insert: &x})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != "Diego Smith" {
		t.Errorf("failed to Push to Delta, got: %+v\n", *n.Ops[0].Insert)
	}
}
func TestPushMultiInsertMathingAttrs(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Insert("Diego ", attr)
	x := "Smith"
	n.Push(Op{Insert: &x, Attributes: attr})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != "Diego Smith" {
		t.Errorf("failed to Push to Delta, got: %+v\n", *n.Ops[0].Insert)
	}
}

func TestPushMultiInsertNonMathingAttrs(t *testing.T) {
	n := New(nil)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	n.Insert("Diego ", attr1)
	x := "Smith"
	attr2 := make(map[string]interface{})
	attr2["bold"] = false
	n.Push(Op{Insert: &x, Attributes: attr2})
	if len(n.Ops) != 2 {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != "Diego " {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", *n.Ops[0].Insert)
	}
	if *n.Ops[1].Insert != "Smith" {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", *n.Ops[1].Insert)
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

	if *x[0].Insert != *expected.Ops[0].Insert {
		t.Errorf("expected %+v but got %+v\n", *expected.Ops[0].Insert, *x[0].Insert)
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

	if *x[0].Insert != *expected.Ops[0].Insert {
		t.Errorf("expected %+v but got %+v\n", *expected.Ops[0].Insert, *x[0].Insert)
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

	if *expected.Ops[0].Insert != *x[0].Insert {
		t.Errorf("expected '%+v' but got '%+v'\n", *expected.Ops[0].Insert, *x[0].Insert)
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

	if *x[0].Insert != *expected.Ops[0].Insert || *x[0].Insert != "B" {
		t.Errorf("expected %+v but got %+v\n", *expected.Ops[0].Insert, *x[0].Insert)
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
	if *x[0].Insert != "HelXlo" {
		t.Errorf("expected 'HelXlo' but got %+v\n", *x[0].Insert)
	}
}
func TestDeltaComposeInsertDeleteOrder(t *testing.T) {
	a := New(nil).Insert("Hello", nil)
	b := New(nil).Insert("Hello", nil)

	insertFirst := New(nil).Retain(3, nil).Insert("X", nil).Delete(1)
	deleteFirst := New(nil).Retain(3, nil).Delete(1).Insert("X", nil)

	xa := a.Compose(*insertFirst).Ops
	if *xa[0].Insert != "HelXo" {
		t.Errorf("expected 'HelXo' but got '%+v'\n", *xa[0].Insert)
	}
	xb := b.Compose(*deleteFirst).Ops
	if *xb[0].Insert != "HelXo" {
		t.Errorf("expected 'HelXo' but got '%+v'\n", *xb[0].Insert)
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
	if *x[0].Insert != "A" {
		t.Errorf("expected 'A' but got '%+v'\n", *x[0].Insert)
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
	if *x[0].Insert != "T" {
		t.Errorf("expected 'T' but got '%+v'\n", *x[0].Insert)
	}
	if *x[1].Insert != "t" {
		t.Errorf("expected 't' but got '%+v'\n", *x[1].Insert)
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
