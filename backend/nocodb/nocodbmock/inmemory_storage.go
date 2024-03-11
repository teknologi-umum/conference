package nocodbmock

import (
	"errors"
	"sync"
)

var errNotFound = errors.New("not found")
var errInvalidType = errors.New("invalid type")

type storage struct {
	m *sync.Map
}

func newInMemoryStorage() *storage {
	return &storage{m: &sync.Map{}}
}

func (s *storage) GetByTableId(tableId string) (records []map[string]any, err error) {
	value, ok := s.m.Load(tableId)
	if !ok {
		return nil, errNotFound
	}

	v, ok := value.([]map[string]any)
	if !ok {
		return nil, errInvalidType
	}

	return v, nil
}

func (s *storage) GetByRecordId(tableId string, recordId int64) (record map[string]any, err error) {
	records, err := s.GetByTableId(tableId)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record["Id"] == recordId {
			return record, nil
		}
	}

	return nil, errNotFound
}

func (s *storage) Insert(tableId string, records []map[string]any) (ids []int64, err error) {
	var lastRecordId int64 = 0
	if oldRecords, err := s.GetByTableId(tableId); err == nil {
		lastIndex := len(oldRecords) - 1
		if v, ok := oldRecords[lastIndex]["Id"]; ok {
			if i, ok := v.(int64); ok {
				lastRecordId = i
			}
		}
	}

	var recordIds []int64
	for i := 0; i < len(records); i++ {
		id := lastRecordId + 1
		records[i]["Id"] = id
		recordIds = append(recordIds, id)
		lastRecordId++
	}

	s.m.Store(tableId, records)
	return recordIds, nil
}

func (s *storage) Update(tableId string, records []map[string]any) (ids []int64, err error) {
	oldRecords, err := s.GetByTableId(tableId)
	if err != nil {
		return nil, err
	}

	var recordIds []int64

	for i := 0; i < len(oldRecords); i++ {
		// I know that this is O(n^2) but because this is a mock, I don't really care
		for _, record := range records {
			if oldRecords[i]["Id"] == record["Id"] {
				// Found one
				for key, value := range record {
					oldRecords[i][key] = value
				}

				id, ok := record["Id"].(int64)
				if ok {
					recordIds = append(recordIds, id)
				}

				break
			}
		}
	}

	s.m.Store(tableId, oldRecords)
	return recordIds, nil
}
