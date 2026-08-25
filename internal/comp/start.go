package comp

import (
	"coldstore/internal/evap"
	"coldstore/internal/valve"
)

func (b *Bank) Start(v *valve.SupplyValve) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if v != nil {
		if err := v.Open(); err != nil {
			return err
		}
	}
	b.running = true
	return nil
}

// ClaimLead is the master/standby election entry point. It delegates the
// vacancy-check-plus-assignment to lead.Claim, which performs the whole
// decision under a single lock hold. This is what makes the election
// atomic: two units racing to become leader at the same instant cannot
// both observe an empty role and both win — the first Claim to take the
// lock writes its id, and the second sees a non-matching leader and
// loses. Do NOT inline the Leader()/SetLeader() pair here: splitting the
// check from the write across two (locked) method calls reintroduces the
// TOCTOU window that let two compressors load at once during defrost end.
func (b *Bank) ClaimLead(lead *evap.Lead) bool {
	if !lead.Claim(b.id) {
		return false
	}
	b.mu.Lock()
	b.running = true
	b.mu.Unlock()
	return true
}
