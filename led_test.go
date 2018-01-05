package main

import "testing"

func TestLEDStringGreen(t *testing.T) {
	expected := "Toggle GREEN"
	actual := getLEDString("green")

	if actual != expected {
		t.Errorf("Test failed, expected: %s, got: %s", expected, actual)
	}
}

func TestLEDStringRed(t *testing.T) {
	expected := "Toggle RED"
	actual := getLEDString("red")

	if actual != expected {
		t.Errorf("Test failed, expected: %s, got: %s", expected, actual)
	}
}

func TestLEDStringYellow(t *testing.T) {
	expected := "Toggle YELLOW"
	actual := getLEDString("yellow")

	if actual != expected {
		t.Errorf("Test failed, expected: %s, got: %s", expected, actual)
	}
}

func TestBuildGpio(t *testing.T) {
	expected := Gpio{0, "name", 1}
	actual := buildGpio(0, "name", 1)

	if actual != expected {
		t.Errorf("Test failed")
	}
}
