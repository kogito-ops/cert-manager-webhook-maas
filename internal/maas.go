package internal

type Config struct {
	ApiKey, ZoneName, ApiUrl, ApiVersion string
}

type DNSResourceRecord struct {
	ID       int    `json:"id"`
	Type     string `json:"rrtype"`
	Name     string `json:"name"`
	Data     string `json:"rrdata"`
	TTL      int    `json:"ttl"`
	FQDN     string `json:"fqdn"`
	DomainID int    `json:"domain"`
}

type DNSResourceRecordResponse struct {
	Records []DNSResourceRecord `json:"dnsresourcerecords"`
}

type Domain struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DomainResponse struct {
	Domains []Domain `json:"domains"`
}
