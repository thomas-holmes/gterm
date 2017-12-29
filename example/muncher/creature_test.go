package main

import "testing"

func TestResourceRegenWholeNumbers(t *testing.T) {
	r := Resource{Current: 1, Max: 5, RegenRate: 1}
	r.Regen()
	expected := 2
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	expected = 3
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	expected = 4
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}

func TestResourceRegenCapsAtMax(t *testing.T) {
	r := Resource{Current: 4, Max: 5, RegenRate: 1}
	r.Regen()
	expected := 5
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}

func TestResourceRegenOneHalf(t *testing.T) {
	r := Resource{Current: 1, Max: 5, RegenRate: 0.5}
	r.Regen()
	expected := 1
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	expected = 2
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	expected = 2
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	expected = 3
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}

func TestResourceRegenOneQuarter(t *testing.T) {
	r := Resource{Current: 1, Max: 5, RegenRate: 0.25}
	r.Regen()
	expected := 1
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	expected = 1
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	expected = 1
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Regen()
	expected = 2
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}
