package umcols

//UrlMetric colums
const (
	// Flags for urlMetrics
	// Title of page if available
	Title int64= 1
	// Canonical form of the url
	Url int64= 4
	// The subdomain of the url
	Subdomain int64= 8
	// The root domain of the url
	RootDomain int64= 16
	// The number of juice-passing external links to the url
	ExternalLinks int64= 32
	// The number of juice-passing external links to the subdomain
	SubdomainExternalLinks int64= 64
	// The number of juice-passing external links
	RootDomainExternalLinks int64= 128
	// The number of juice-passing links (internal or external) to the url
	JuicePassingLinks int64= 256
	// The number of subdomains with any pages linking to the url
	SubdomainsLinking int64= 512
	// The number of root domains with any pages linking to the url
	RootDomainsLinking int64= 1024
	// The number of links (juice-passing or not, internal or external) to the url
	Links int64= 2048
	// The number of subdomains with any pages linking to the subdomain of the url
	SubdomainSubdomainsLinking int64= 4096
	// The number of root domains with any pages linking to the root domain of the url
	RootDomainRootDomainsLinking int64= 8192
	// The mozRank of the url.  Requesting this metric will provide both the
	// pretty 10-point score (in umrp) and the raw score (umrr)
	MozRank int64= 16384
	// The mozRank of the subdomain of the url. Requesting this metric will
	//provide both the pretty 10-point score (fmrp) and the raw score (fmrr)
	SubdomainMozRank int64= 32768
	// The mozRank of the Root Domain of the url. Requesting this metric will
	// provide both the pretty 10-point score (pmrp) and the raw score (pmrr)
	RootDomainMozRank int64= 65536
	// The mozTrust of the url. Requesting this metric will provide both the
	// pretty 10-point score (utrp) and the raw score (utrr).
	MozTrust int64= 131072
	// The mozTrust of the subdomain of the url.  Requesting this metric will
	// provide both the pretty 10-point score (ftrp) and the raw score (ftrr)
	SubdomainMozTrust int64= 262144
	// The mozTrust of the root domain of the url.  Requesting this metric
	// will provide both the pretty 10-point score (ptrp) and the raw score (ptrr)
	RootDomainMozTrust int64= 524288
	// The portion of the url's mozRank coming from external links.  Requesting
	// this metric will provide both the pretty 10-point score (uemrp) and the raw
	// score (uemrr)
	ExternalMozRank int64= 1048576
	// The portion of the mozRank of all pages on the subdomain coming from
	// external links.  Requesting this metric will provide both the pretty
	// 10-point score (fejp) and the raw score (fejr)
	SubdomainExternalDomainJuice int64 = 2097152
	// The portion of the mozRank of all pages on the root domain coming from
	// external links.  Requesting this metric will provide both the pretty
	// 10-point score (pejp) and the raw score (pejr)
	RootDomainExternalDomainJuice int64= 4194304
	// The mozRank of all pages on the subdomain combined.  Requesting this
	// metric will provide both the pretty 10-point score (fjp) and the raw score (fjr)
	SubdomainDomainJuice int64= 8388608
	// The mozRank of all pages on the root domain combined.  Requesting this
	// metric will provide both the pretty 10-point score (pjp) and the raw score (pjr)
	RootDomainDomainJuice int64= 16777216
	// The HTTP status code recorded by Linkscape for this URL (if available)
	HttpStatusCode int64= 536870912
	// Total links (including internal and nofollow links) to the subdomain of
	// the url in question
	LinksToSubdomain int64= 4294967296
	// Total links (including internal and nofollow links) to the root domain
	// of the url in question.
	LinksToRootDomain int64= 8589934592
	// The number of root domains with at least one link to the subdomain of
	// the url in question
	RootDomainsLinkingToSubdomain int64= 17179869184
	// A score out of 100-points representing the likelihood for arbitrary content
	// to rank on this page
	PageAuthority int64= 34359738368
	// A score out of 100-points representing the likelihood for arbitrary content
	// to rank on this dom
	DomainAuthority int64= 68719476736
	// This is the set of all free fields
	FreeCols = (Title | Url | ExternalLinks | Links | MozRank | SubdomainMozRank | HttpStatusCode | PageAuthority | DomainAuthority)

)
