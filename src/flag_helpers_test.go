package main

import "testing"

func TestGlobalAndCharFlagHelpersRoundTripBits(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.specialFlag = 0
	sys.setGSF(GSF_intro)
	if !sys.gsf(GSF_intro) {
		t.Fatal("gsf() should report a set global flag")
	}
	sys.unsetGSF(GSF_intro)
	if sys.gsf(GSF_intro) {
		t.Fatal("unsetGSF() should clear a global flag")
	}

	c := &Char{}
	if c.scf(SCF_disabled) || c.csf(CSF_destroy) || c.asf(ASF_noinput) {
		t.Fatal("fresh Char should not report any tested flags")
	}

	c.setSCF(SCF_disabled)
	if !c.scf(SCF_disabled) {
		t.Fatal("setSCF() should set a system flag")
	}
	c.unsetSCF(SCF_disabled)
	if c.scf(SCF_disabled) {
		t.Fatal("unsetSCF() should clear a system flag")
	}

	c.setCSF(CSF_destroy)
	if !c.csf(CSF_destroy) {
		t.Fatal("setCSF() should set a char special flag")
	}
	c.unsetCSF(CSF_destroy)
	if c.csf(CSF_destroy) {
		t.Fatal("unsetCSF() should clear a char special flag")
	}

	c.setASF(ASF_noinput)
	if !c.asf(ASF_noinput) {
		t.Fatal("setASF() should set an assert flag")
	}
	c.unsetASF(ASF_noinput)
	if c.asf(ASF_noinput) {
		t.Fatal("unsetASF() should clear an assert flag")
	}
}
