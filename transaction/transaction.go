package transaction

type Transaction []byte

func (t *Transaction) ToEF() []byte {
	return []byte(*t)
}
