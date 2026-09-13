package server

type TLSConfig struct {
	CertPath string
	KeyPath  string
}

func (c *TLSConfig) Empty() bool {
	return c == nil || c.CertPath == "" || c.KeyPath == ""
}
