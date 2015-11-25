package gocsv

type Sample struct {
	Foo string `csv:"foo"`
	Bar int    `csv:"BAR"`
	Baz string `csv:"Baz"`
}

type EmbedSample struct {
	Qux string `csv:"first"`
	Sample
	Ignore string `csv:"-"`
	Quux   string `csv:"last"`
}

type SkipFieldSample struct {
	EmbedSample
	MoreIgnore string `csv:"-"`
	Corge      string `csv:"abc"`
}
