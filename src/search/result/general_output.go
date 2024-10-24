package result

type GeneralOutput struct {
	generalOutputJSON
}

type generalOutputJSON struct {
	General

	FqdnHash          string `json:"fqdn_hash,omitempty"`
	FqdnHashTimestamp string `json:"fqdn_hash_timestamp,omitempty"`
}
