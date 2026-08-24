package journal

import (
	"encoding/json"
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"os"
	"path/filepath"
)

type EvidenceStore struct {
	root     string
	readOnly bool
}

func NewEvidenceStore(root string, readOnly ...bool) (*EvidenceStore, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	store := &EvidenceStore{root: root}
	if len(readOnly) > 0 {
		store.readOnly = readOnly[0]
	}
	return store, nil
}
func (e *EvidenceStore) Save(proof model.FieldProof, image []byte) error {
	if e.readOnly {
		return fmt.Errorf("evidence store is read-only")
	}
	if proof.ReportID == "" || proof.ImageName == "" {
		return fmt.Errorf("field proof requires report and image")
	}
	if err := os.WriteFile(filepath.Join(e.root, proof.ReportID+".json"), mustJSON(proof), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(e.root, proof.ImageName), image, 0o644)
}
func (e *EvidenceStore) Has(proof model.FieldProof) bool {
	_, a := os.Stat(filepath.Join(e.root, proof.ReportID+".json"))
	_, b := os.Stat(filepath.Join(e.root, proof.ImageName))
	return a == nil && b == nil
}
func mustJSON(value any) []byte { data, _ := json.MarshalIndent(value, "", "  "); return data }
