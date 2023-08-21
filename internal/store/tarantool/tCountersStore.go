package tarantool

import (
	"github.com/pkg/errors"
)

type Resp struct {
	Count uint64
}

func (t *Tarantool) IncrCounter(userId string) (uint64, error) {
	var counts []Resp
	err := t.conn.CallTyped("incr_message", []interface{}{userId}, &counts)
	if err != nil {
		return 0, err
	}

	return counts[0].Count, nil

}

func (t *Tarantool) DecrCounter(userId string) (uint64, error) {
	var counts []Resp
	err := t.conn.CallTyped("decr_message", []interface{}{userId}, &counts)
	if err != nil {
		return 0, err
	}
	if len(counts) != 1 {
		return 0, errors.Errorf("Cannot find user with id: %s", counts)
	} else {
		return counts[0].Count, nil
	}
}

func (t *Tarantool) GetTotalMessages(userId string) (uint64, error) {
	var count []Resp
	err := t.conn.CallTyped("get_message_count", []interface{}{userId}, &count)
	if err != nil {
		return 0, err
	}
	if len(count) != 1 {
		return 0, errors.Errorf("Cannot find user with id: %s", count)
	} else {
		return count[0].Count, nil
	}
}
