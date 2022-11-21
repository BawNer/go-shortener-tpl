package storage

type InFile struct {
	Producer *Producer
}

func NewFileStorage(filepath string) (*InFile, error) {
	producer, err := NewProducer(filepath)
	if err != nil {
		return nil, err
	}
	return &InFile{
		Producer: producer,
	}, nil
}
