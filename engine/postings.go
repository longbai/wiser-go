package engine

type PostingsBuffer struct {
	Count int

}
// expect ngram is 2
func (p *PostingsBuffer)TextToPostingsLists(documentId int, body string) {
	var lastRune rune
	ignoreLast := false
	for i, r := range body {
		if ignoreChar(r) {
			ignoreLast = true
			continue
		}
		if !ignoreLast && i != 0 {
			s := string([]rune{lastRune, r})
			TokenToPostingsList(s)
		}
		lastRune = r

	}
}



