# Private Wallet Transactions with Polyhedra’s Expander Compiler Collection

I’m thrilled to present my hackathon project: a zero-knowledge proof (ZKP) circuit for private transactions between two wallets, built with Polyhedra’s Expander Compiler Collection (ECC). I whipped this up during the hackathon to explore how we can keep crypto transactions confidential yet verifiable?

## Project Goals

The big idea here is simple: enable private transfers between two wallets without disclosing amounts and or balances. Imagine Wallet A sending some crypto to Wallet B. I wanted to:
- Hide the transaction amount and both wallets’ balances using Pedersen commitments.
- Prove Wallet A has enough funds without showing the numbers.
- Ensure Wallet B gets credited correctly, all while keeping everything under wraps.
- Leverage Polyhedra’s ECC to turn my circuit into something the Expander prover can handle for fast, efficient proofs.

I built this as a proof-of-concept circuit using gnark, then set it up to work with ECC for proof generation.

## How It Works

Here’s the rundown:
1. **Setup**: Wallet A has a balance (`balance_A`), and I commit to it with a Pedersen commitment (`C_A_initial = g^balance_A * h^r`). Same deal for Wallet B’s starting balance and the transfer amount.
2. **Circuit Logic**: I wrote a gnark circuit that checks:
   - `balance_A >= amount` (no overdraft!).
   - New balances are computed: `balance_A - amount` and `balance_B + amount`.
   - Final commitments (`C_A_final`, `C_B_final`) match the updated balances.
3. **Polyhedra Magic**: I feed the circuit into the ECC, which spits out a layered circuit IR. Then, the Expander prover generates a ZKP that proves everything’s legit without revealing the private stuff.
4. **Verification**: Anyone can verify the proof and update the public commitments—no secrets exposed!

The `main.go` file contains the circuit code (it’s a simplified version).

## Assumptions

I made a few assumptions:
- I’m using Pedersen commitments for simplicity—they’re solid for hiding values and work well in ZKPs.
- The wallets’ blinding factors (`r`) are reused or simplified in the example. In a real system, I’d generate fresh ones for security.
- I assume a basic blockchain or off-chain setup where commitments are stored and proofs are verified—will be built as a future improvement.
- The ECC and Expander prover are running locally with default configs. I leaned on Polyhedra’s docs and GitHub for how they’d fit in.

## Limitations

This being a hackathon project, it’s not perfect (yet!). Here’s where it stands:
- **Simplified Blinding**: I didn’t fully manage the blinding factors for commitments, which could leak info if not handled right. It’s a quick-and-dirty demo fix.
- **No Full Integration**: I didn’t hook up the ECC output to the Expander prover live—time ran short! It’s more of a “here’s how it’d work” than a running system.
- **Basic Circuit**: The circuit only handles a single transfer with no extras like fees or multi-party checks.
- **Testing**: I didn’t get to test the proof generation end-to-end. 

## Future Development & Improvements

If I had more time (or after the hackathon!), here’s what I’d do:
- **Full ECC Pipeline**: Get the circuit fully compiled with ECC and generate real proofs with the Expander prover.
- **Better Security**: Fix the blinding factor management so each commitment is truly independent and secure.
- **On-Chain Integration**: Hook this into a testnet (maybe Ethereum or a ZK-rollup) to see it work with real wallets.
- **Fancy Features**: Add support for transaction fees, multiple recipients, or even a privacy mixer-style setup.
- **UI Demo**: Build a simple front-end where users can input balances and amounts, then see the proof in action.

I’m super excited to polish this up post-hackathon.

## Getting Started

Want to check it out? Here’s how:
1. Clone this repo: `git clone https://github.com/IKalonji2/ExpanderPrivateTransactions.git`
2. Install Go and gnark dependencies: `go mod tidy`
3. Peek at `main.go` for the circuit logic.
4. For the ECC part, grab Polyhedra’s tools from their GitHub (`PolyhedraZK/ExpanderCompilerCollection`) and follow their setup guide.
5. Run the Expander prover separately.


## Why This Matters

I built this because privacy in crypto is huge—people want control over their financial data. Polyhedra’s ECC and Expander stack made it possible to dream up something efficient and scalable in just a weekend. I hope the judges see the vision: a stepping stone to private, trustless transactions that don’t compromise on speed or security.
