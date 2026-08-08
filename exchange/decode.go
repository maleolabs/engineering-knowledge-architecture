package exchange

// This file implements the standalone unit decoding helpers of the
// exchange model:
//
//	u, err := exchange.DecodeUnit(unitJSON, content)
//
// DecodeUnit is the read-side counterpart of MarshalUnit (emit.go): it
// reconstructs one exchange.Unit from its canonical unit.json bytes and
// its content payload bytes, applying the same reject-by-default
// unknown-field policy as the package loader (RSF §9.5). It exists for
// consumers that hold units as raw bytes without a surrounding package
// — the canonical store of the EKA Knowledge Runtime stores exactly
// (unit.json || content) per immutable payload and reconstructs the
// model on push; the integrity engine decodes every stored payload to
// verify it.

// DecodeUnit strictly decodes unitJSON into a Unit (unknown fields
// rejected, RSF §9.5) and attaches content as the unit's
// ContentPayload. The decoded unit carries everything the serializer
// emits: Identity, CanonicalIdentityForm, Revision, metadata,
// StateVector, ChangeLog, Relationships, Classification, Phase and the
// Content reference. It is the inverse of MarshalUnit for bytes
// produced by this package (and tolerates a trailing LF, the package
// entry normalization of RSF §9.3).
func DecodeUnit(unitJSON, content []byte) (*Unit, error) {
	var u Unit
	if err := strictDecode("unit.json", unitJSON, &u); err != nil {
		return nil, err
	}
	u.ContentPayload = content
	return &u, nil
}
