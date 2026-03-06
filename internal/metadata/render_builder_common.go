package metadata

type Example struct {
	Header      *string
	Description *string
	HCL         []byte
}

type ImportId struct {
	Format  string
	Example string
}
