package gocsv

type Sample struct {
	Foo  string  `csv:"foo"`
	Bar  int     `csv:"BAR"`
	Baz  string  `csv:"Baz"`
	Frop float64 `csv:"Quux"`
	Blah *int    `csv:"Blah"`
}

type EmbedSample struct {
	Qux string `csv:"first"`
	Sample
	Ignore string  `csv:"-"`
	Grault float64 `csv:"garply"`
	Quux   string  `csv:"last"`
}

type SkipFieldSample struct {
	EmbedSample
	MoreIgnore string `csv:"-"`
	Corge      string `csv:"abc"`
}
