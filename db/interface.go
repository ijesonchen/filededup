package db

// FileInfo to dedup
type FileInfo struct {
	ID   string `bson:"_id"`
	Name string
	Path string
	Ext  string
	Size int64
	GB   float32 // size in GigaByte

	// hash value parameters
	Byte2Hash    int64
	FingerPrint  int64
	FingerMethod string
	HashValue    []byte
	HashType     string
}

// IClient .
type IClient interface {
	Insert(d *FileInfo) (err error)
	FindName(name string) (files []FileInfo, err error)
	LikeName(name string) (files []FileInfo, err error)
	Delete(id string) (err error)
	List(count, skip int64, order string) (files []FileInfo, err error)
}

// NewClient .
func NewClient(host, port, dbname, coll string) (client IClient, err error) {
	return
}
