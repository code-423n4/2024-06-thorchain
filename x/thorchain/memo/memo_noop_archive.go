package thorchain

func ParseNoOpMemoV1(parts []string) (NoOpMemo, error) {
	return NewNoOpMemo(GetPart(parts, 1)), nil
}
