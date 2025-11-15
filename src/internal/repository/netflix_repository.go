package repository

type RawNetflixRecord struct {
	Title string
	Date  string
}

type NetflixRepository interface {
	ReadRawCSV(path string) ([]RawNetflixRecord, error)
}
