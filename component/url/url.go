// Package url provides URL construction helpers with prefix support for route generation.
package url

// Url represents a URL with optional prefix for routing
type Url struct {
	url    string
	prefix string
}

// NewUrl creates a new URL
func NewUrl(url string) *Url {
	return &Url{
		url:    url,
		prefix: "",
	}
}

// NewUrlPrefix creates a new URL with prefix
func NewUrlPrefix(url string, prefix string) *Url {
	return &Url{
		url:    url,
		prefix: prefix,
	}
}

// Print returns the URL without prefix
func (u *Url) Print() string {
	return u.url
}

// PrintPrefix returns the full URL with prefix
func (u *Url) PrintPrefix() string {
	return u.prefix + u.url
}

// Add appends a path segment to the URL and returns the Url for method chaining.
func (u *Url) Add(url string) *Url {
	u.url += "/" + url
	return u
}

// AddPrefix prepends a prefix to the existing prefix and returns the Url for method chaining.
func (u *Url) AddPrefix(prefix string) *Url {
	u.prefix = prefix + "/" + u.prefix
	return u
}
