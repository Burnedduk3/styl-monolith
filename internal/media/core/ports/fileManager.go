package ports

type FileManager interface {
	CheckIfFolderCreated(path string) (bool, error)
	CreateFolder(path string) error
	DeleteFolder(path string) error
	DeleteImage(path string) error
	GetImage(path string) ([]byte, error)
	SaveImage(path string, data []byte) error
}
