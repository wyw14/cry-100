package journal

import "github.com/wyw14/cry-100/internal/model"

func (s *Store) AppendOperation(record model.OperationRecord) error {
	_, err := AppendValue(s, "dispatch.result", string(record.OperationID), record)
	return err
}
