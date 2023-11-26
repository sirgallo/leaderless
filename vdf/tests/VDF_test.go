package vdftest

import "math/big"
import "testing"

import "github.com/sirgallo/athn/vdf"


const ITERATIONS = 100
const GENESIS_INPUT = "f92d7accb0004af0565148cdb657e0199574ad70dbaed4ecdfdec345943bb24a2612cfe3e45bafb59120dce8c4550640b56664072e814a55266d0441f9708650"

func TestVDF(t *testing.T) {
	v := new(big.Int)
	v.SetString(GENESIS_INPUT, 16)

	computed, proof, computeErr := vdf.ComputeVDF(v, nil)
	if computeErr != nil { t.Errorf("error computing vdf: %s\n", computeErr.Error() )}

	isVerified := vdf.VerifyVDF(computed, v, proof, nil)
	if ! isVerified { t.Error("vdf output was not verified") }

	t.Logf("input: %d\noutput: %d\nproof: %d\nisVerified: %t\n", v, computed, proof, isVerified)
}

func TestVDFSetDifficulty(t *testing.T) {
	v := new(big.Int)
	v.SetString(GENESIS_INPUT, 16)

	for i := range make([]int, ITERATIONS) {
		computed, proof, computeErr := vdf.ComputeVDF(v, nil)
		if computeErr != nil { t.Errorf("error computing vdf: %s\n", computeErr.Error() )}
		
		isVerified := vdf.VerifyVDF(computed, v, proof, nil)
		if ! isVerified { t.Error("vdf output was not verified") }
		
		t.Logf("iteration: %d\noutput: %d\nproof: %d\nisVerified: %t\n\n", i, computed, proof, isVerified)
		v = computed
	}
}

func TestVDFDynamicDifficulty(t *testing.T) {
	v := new(big.Int)
	v.SetString(GENESIS_INPUT, 16)

	for i := range make([]int, ITERATIONS) {
		cum := uint64(i)
		computed, proof, computeErr := vdf.ComputeVDF(v, &cum)
		if computeErr != nil { t.Errorf("error computing vdf: %s\n", computeErr.Error() )}

		isVerified := vdf.VerifyVDF(computed, v, proof, &cum)
		if ! isVerified { t.Error("vdf output was not verified")}

		t.Logf("iteration: %d\noutput: %d\nproof: %d\nisVerified: %t\n\n", i, computed, proof, isVerified)
		v = computed
	}

	t.Log("Done")
}