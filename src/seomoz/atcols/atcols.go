package atcols

//Anchor Text Cols
const (
	// The anchor text term or phrase
	Term = 2
	// The number of internal pages linking with this term or phrase
	InternalPagesLinking = 8
	// The number of subdomains on the same root domain with at least one link with this term or phrase
	InternalSubdomainsLinking = 16
	// The number of external pages linking with this term or phrase
	ExternalPagesLinking = 32
	// The number of external subdomains with at least one link with this term or phrase
	ExternalSubdomainsLinking = 64
	// The number of (external) root domains with at least one link with this term or phrase
	ExternalRootDomainsLinking = 128
	// The amount of mozRank passed over all internal links with this term or phrase (on the 10 point scale)
	InternalMozRankPassed = 256
	// The amount of mozRank passed over all external links with this term or phrase (on the 10 point scale)
	ExternalMozRankPassed = 512
	// Currently only "1" is used to indicate the term or phrase is found in an image link
	Flags = 1024

	// This is the set of all free fields
	FreeCols = (Term |
		InternalPagesLinking |
		InternalSubdomainsLinking |
		ExternalPagesLinking |
		ExternalSubdomainsLinking |
		ExternalRootDomainsLinking |
		InternalMozRankPassed |
		ExternalMozRankPassed |
		Flags)

)
