package main

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type PrivateTransactionCircuit struct {
	// Private inputs
	SecretKey         frontend.Variable     `gnark:"secret_key,secret"`
	OldAmount         frontend.Variable     `gnark:"old_amount,secret"`
	AmountSent        frontend.Variable     `gnark:"amount_sent,secret"`
	R_old             frontend.Variable     `gnark:"r_old,secret"`
	R_new             frontend.Variable     `gnark:"r_new,secret"`
	R_receiver        frontend.Variable     `gnark:"r_receiver,secret"`
	MerklePath        [20]frontend.Variable `gnark:"merkle_path,secret"`
	MerklePathIndices [20]frontend.Variable `gnark:"merkle_path_indices,secret"`

	// Public inputs
	MerkleRoot          frontend.Variable `gnark:"merkle_root,public"`
	Nullifier           frontend.Variable `gnark:"nullifier,public"`
	NewSenderCommitment frontend.Variable `gnark:"new_sender_commitment,public"`
	ReceiverCommitment  frontend.Variable `gnark:"receiver_commitment,public"`
}

func (circuit *PrivateTransactionCircuit) Define(api frontend.API) error {

	h, _ := mimc.NewMiMC(api)
	h.Write(circuit.SecretKey)
	publicKey := h.Sum()
	current := publicKey
	for i := 0; i < 20; i++ {
		pathElement := circuit.MerklePath[i]
		pathIndex := circuit.MerklePathIndices[i]
		left := api.Select(pathIndex, pathElement, current)
		right := api.Select(pathIndex, current, pathElement)

		hNode, _ := mimc.NewMiMC(api)
		hNode.Write(left)
		hNode.Write(right)
		current = hNode.Sum()
	}
	api.AssertIsEqual(current, circuit.MerkleRoot)

	Gx := frontend.NewConstant(ecc.BN254.GeneratorX)
	Gy := frontend.NewConstant(ecc.BN254.GeneratorY)
	G := api.NewPoint(Gx, Gy) // Generator point G

	Hx := frontend.NewConstant(ecc.BN254.Hx)
	Hy := frontend.NewConstant(ecc.BN254.Hy)
	H := api.NewPoint(Hx, Hy) // Generator point H

	oldCommitment := G.ScalarMul(api, circuit.OldAmount).Add(H.ScalarMul(api, circuit.R_old))
	newSenderCommitment := G.ScalarMul(api, api.Sub(circuit.OldAmount, circuit.AmountSent)).Add(H.ScalarMul(api, circuit.R_new))
	receiverCommitment := G.ScalarMul(api, circuit.AmountSent).Add(H.ScalarMul(api, circuit.R_receiver))

	api.AssertIsEqual(newSenderCommitment.X, circuit.NewSenderCommitment)
	api.AssertIsEqual(receiverCommitment.X, circuit.ReceiverCommitment)

	api.AssertIsEqual(api.Add(circuit.R_new, circuit.R_receiver), circuit.R_old)

	api.RangeCheck(circuit.AmountSent, 32)
	api.RangeCheck(api.Sub(circuit.OldAmount, circuit.AmountSent), 32)

	nullifier := h.Sum(circuit.SecretKey, oldCommitment.X, oldCommitment.Y)
	api.AssertIsEqual(nullifier, circuit.Nullifier)

	return nil
}

func main() {
	var circuit PrivateTransactionCircuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), frontend.Expander, &circuit)
	if err != nil {
		panic(err)
	}
}
