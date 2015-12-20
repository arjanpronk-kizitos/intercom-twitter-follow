package intercom

import (
	"github.com/antoinefinkelstein/intercom-twitter-follow/Godeps/_workspace/src/github.com/intercom/intercom-go/interfaces"
)

// A Client manages interacting with the Intercom API.
type Client struct {
	// Services for interacting with various resources in Intercom.
	Admins    AdminService
	Companies CompanyService
	Contacts  ContactService
	Events    EventService
	Segments  SegmentService
	Tags      TagService
	Users     UserService

	// Mappings for resources to API constructs
	AdminRepository   AdminRepository
	CompanyRepository CompanyRepository
	ContactRepository ContactRepository
	EventRepository   EventRepository
	SegmentRepository SegmentRepository
	TagRepository     TagRepository
	UserRepository    UserRepository

	// AppID For Intercom.
	AppID string

	// APIKey for Intercom's API. See http://app.intercom.io/apps/api_keys.
	APIKey string

	// HTTP Client used to interact with the API.
	HTTPClient interfaces.HTTPClient

	baseURI       string
	clientVersion string
	debug         bool
}

const (
	defaultBaseURI = "https://api.intercom.io"
	clientVersion  = "0.0.1"
)

type option func(c *Client) option

// Set Options on the Intercom Client, see TraceHTTP, BaseURI and SetHTTPClient.
func (c *Client) Option(opts ...option) (previous option) {
	for _, opt := range opts {
		previous = opt(c)
	}
	return previous
}

// NewClient returns a new Intercom API client, configured with the default HTTPClient.
func NewClient(appID, apiKey string) *Client {
	intercom := Client{AppID: appID, APIKey: apiKey, baseURI: defaultBaseURI, debug: false, clientVersion: clientVersion}
	intercom.HTTPClient = interfaces.NewIntercomHTTPClient(intercom.AppID, intercom.APIKey, &intercom.baseURI, &intercom.clientVersion, &intercom.debug)
	intercom.setup()
	return &intercom
}

// TraceHTTP turns on HTTP request/response tracing for debugging.
func TraceHTTP(trace bool) option {
	return func(c *Client) option {
		previous := c.debug
		c.debug = trace
		return TraceHTTP(previous)
	}
}

// BaseURI sets a base URI for the HTTP Client to use. Defaults to "https://api.intercom.io".
// Typically this would be used during testing to point to a stubbed service.
func BaseURI(baseURI string) option {
	return func(c *Client) option {
		previous := c.baseURI
		c.baseURI = baseURI
		return BaseURI(previous)
	}
}

// SetHTTPClient sets a HTTPClient for the Intercom Client to use.
// Useful for customising timeout behaviour etc.
func SetHTTPClient(httpClient interfaces.HTTPClient) option {
	return func(c *Client) option {
		previous := c.HTTPClient
		c.HTTPClient = httpClient
		c.setup()
		return SetHTTPClient(previous)
	}
}

func (c *Client) setup() {
	c.AdminRepository = AdminAPI{httpClient: c.HTTPClient}
	c.CompanyRepository = CompanyAPI{httpClient: c.HTTPClient}
	c.ContactRepository = ContactAPI{httpClient: c.HTTPClient}
	c.EventRepository = EventAPI{httpClient: c.HTTPClient}
	c.SegmentRepository = SegmentAPI{httpClient: c.HTTPClient}
	c.TagRepository = TagAPI{httpClient: c.HTTPClient}
	c.UserRepository = UserAPI{httpClient: c.HTTPClient}
	c.Admins = AdminService{Repository: c.AdminRepository}
	c.Companies = CompanyService{Repository: c.CompanyRepository}
	c.Contacts = ContactService{Repository: c.ContactRepository}
	c.Events = EventService{Repository: c.EventRepository}
	c.Segments = SegmentService{Repository: c.SegmentRepository}
	c.Tags = TagService{Repository: c.TagRepository}
	c.Users = UserService{Repository: c.UserRepository}
}
