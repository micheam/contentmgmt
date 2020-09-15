package imgcontent

// ErrIllegalFilename - illgal filename.
type ErrIllegalFilename string

func (e ErrIllegalFilename) Error() string {
	return string(e)
}
