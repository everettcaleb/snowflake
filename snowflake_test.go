package main

import (
	"fmt"
	"testing"
	"time"
)

func fixedNow(s string) func() time.Time {
	return func() time.Time {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			panic(err)
		}
		return t
	}
}

func incrementingNow(s string) func() time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return func() time.Time {
		t = t.Add(time.Microsecond)
		return t
	}
}

func TestMSTimestamp(t *testing.T) {
	// fmt.Println(msTimestamp())
	now = fixedNow("2018-07-25T01:14:00Z")
	tsms := uint64(now().Unix() * 1000)
	v := msTimestamp()
	if tsms != msTimestamp() {
		t.Error("Expected msTimestamp() ->", tsms, "got", v)
	}

	tsns := uint64(now().UnixNano() / 1e6)
	v = msTimestamp()
	if tsns != msTimestamp() {
		t.Error("Expected msTimestamp() ->", tsns, "got", v)
	}
}

func TestSplitSnowflakeID(t *testing.T) {
	// epoch, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
	// if err != nil {
	// 	panic(err)
	// }

	/*sid := */
	splitSnowflakeID(3044420458909792707)
	// fmt.Println("Original ID:", sid.id)
	// fmt.Println("Timestamp (ms):", sid.timestamp)
	// fmt.Println("Machine ID:", sid.machineID)
	// fmt.Println("Counter:", sid.counter)
	// fmt.Println()

	// realTS := epoch.Add(time.Duration(sid.timestamp) * time.Millisecond)
	// fmt.Println("Real Timestamp (unix):", realTS.Unix())
}

func TestNextID(t *testing.T) {
	epoch, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
	if err != nil {
		panic(err)
	}

	now = fixedNow("2039-01-01T00:00:00Z")

	config := defaultConfig(epoch)
	config.UseMilliseconds = false
	sg := newSnowflakeGenerator(config, 1023)

	ids := make([]snowflakeID, 2500)
	for i := 0; i < 2500; i++ {
		ids[i] = sg.NextID()
	}

	sids := make([]*snowflakeIDParts, 2500)
	for i, id := range ids {
		sids[i] = splitSnowflakeID(id)
	}

	// check to make sure they're ordered
	for i, id := range ids {
		if i > 0 && id-1 != ids[i-1] {
			t.Error("IDs should be sequential!")
		}
	}

	// check them all
	for i, sid := range sids {
		if i > 0 && sid.counter <= sids[i-1].counter {
			t.Error("Counter skipped")
		}
		if sid.machineID != 1023 {
			t.Error("Invalid machine ID")
		}
		if i > 0 && sid.timestamp != sids[i-1].timestamp {
			t.Error("Time issue")
		}
	}

	// check the format of a few

	// TEMP: printing out the first 100
	for _, id := range ids {
		fmt.Println(id)
	}
}
