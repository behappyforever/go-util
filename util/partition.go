package util

type HandleFunc func(l int, h int)

func HandleByPartition(totalLen int, partitionLen int, hf HandleFunc) {
	if partitionLen <= 0 || totalLen <= 0 {
		return
	}
	partitions := totalLen / partitionLen
	var i int
	for i = 0; i < partitions; i++ {
		hf(i*partitionLen, (i+1)*partitionLen)
	}
	if rest := totalLen % partitionLen; rest != 0 {
		hf(i*partitionLen, i*partitionLen+rest)
	}
}
