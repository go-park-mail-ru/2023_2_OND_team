package check

type CheckSize func(width, height float64) bool

func BothSidesFallIntoRange(minSize, maxSize int) CheckSize {
	return func(width, height float64) bool {
		if width < float64(minSize) || width > float64(maxSize) ||
			height < float64(minSize) || height > float64(minSize) {

			return false
		}
		return true
	}
}

func AnySize(width, height float64) bool {
	return true
}
