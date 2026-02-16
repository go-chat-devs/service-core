package utils

import "github.com/google/uuid"

func SortUUID(uuid1, uuid2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if uuid1.String() < uuid2.String() {
		return uuid1, uuid2
	}
	return uuid2, uuid1
}
