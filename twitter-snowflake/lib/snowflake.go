package lib

import (
	"errors"
	"sync"
	"time"
)

// Snowflake ID Structure:
//
// +-------------+------------+------------+------------+
// | Timestamp   | Datacenter |  Worker ID | Sequence   |
// |   41 bits   |   5 bits   |   5 bits   |  12 bits   |
// +-------------+------------+------------+------------+
//
//  63           22           17           12          0
// +-------------+------------+------------+------------+
// |   Timestamp |  Datacenter |  Worker ID | Sequence  |
// +-------------+------------+------------+------------+
//
// - Timestamp: Milliseconds since the epoch (customizable)
// - Datacenter ID: 5 bits, giving 32 datacenters
// - Worker ID: 5 bits, giving 32 workers per datacenter
// - Sequence: 12 bits, rolls over within the same millisecond

// This structure allows for:
// - 69 years of IDs from the custom epoch
// - 1024 unique worker/datacenter combinations
// - 4096 IDs per millisecond per worker/datacenter

// Constants for the Snowflake algorithm
const (
	// epochStart is the Twitter epoch (November 4, 2010).
	// You can customize this to your own epoch.
	epochStart = int64(1288834974657)

	// Bit allocation
	workerIDBits     = uint(5)  // 5 bits for worker ID
	datacenterIDBits = uint(5)  // 5 bits for datacenter ID
	sequenceBits     = uint(12) // 12 bits for sequence number

	// Maximum values for each component
	maxWorkerID     = -1 ^ (-1 << workerIDBits)     // 31
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits) // 31
	maxSequence     = -1 ^ (-1 << sequenceBits)     // 4095

	// Bit shifting
	timeShift         = workerIDBits + datacenterIDBits + sequenceBits
	datacenterIDShift = workerIDBits + sequenceBits
	workerIDShift     = sequenceBits
)

// Snowflake struct holds the state for the Snowflake ID generator
type Snowflake struct {
	mutex         sync.Mutex // Ensures thread-safety
	lastTimestamp int64      // Last timestamp used for ID generation
	workerID      int64      // Worker ID (0-31)
	datacenterID  int64      // Datacenter ID (0-31)
	sequence      int64      // Sequence number (0-4095)
}

// NewSnowflake creates a new Snowflake instance
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	// Validate workerID
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID must be between 0 and 31")
	}
	// Validate datacenterID
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("datacenter ID must be between 0 and 31")
	}
	// Return a new Snowflake instance
	return &Snowflake{
		workerID:     workerID,
		datacenterID: datacenterID,
	}, nil
}

// NextID generates the next unique ID
func (s *Snowflake) NextID() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / 1000000

	// Handle clock moving backwards
	if timestamp < s.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	// Handle same millisecond
	if timestamp == s.lastTimestamp {
		// Increment sequence
		s.sequence = (s.sequence + 1) & maxSequence
		// If sequence overflows, wait for next millisecond
		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// Reset sequence for new millisecond
		s.sequence = 0
	}

	// Update last timestamp
	s.lastTimestamp = timestamp

	// Construct the 64-bit ID
	id := ((timestamp - epochStart) << timeShift) | // Timestamp bits
		(s.datacenterID << datacenterIDShift) | // Datacenter ID bits
		(s.workerID << workerIDShift) | // Worker ID bits
		s.sequence // Sequence bits

	return id, nil
}
