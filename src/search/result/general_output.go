package result

type GeneralOutput struct {
	generalOutputJSON
}

type generalOutputJSON struct {
	General

	FaviconHash string `json:"favicon_hash,omitempty"`
}
