package dataStorageHeap

type DataStorage struct {
	Ctime int64
	Data  string
}

type DataStorageHeap []*DataStorage

func (dsh DataStorageHeap) Len() int { return len(dsh) }

func (dsh DataStorageHeap) Less(i, j int) bool { return dsh[i].Ctime < dsh[j].Ctime }

func (dsh DataStorageHeap) Swap(i, j int) { dsh[i], dsh[j] = dsh[j], dsh[i] }

func (dsh *DataStorageHeap) Push(value any) {
	*dsh = append(*dsh, value.(*DataStorage))
}

func (dsh *DataStorageHeap) Pop() any {
	current := *dsh
	n := len(current)
	if n == 0 {
		return nil
	}
	value := current[n-1]
	*dsh = current[0 : n-1]
	return value
}
