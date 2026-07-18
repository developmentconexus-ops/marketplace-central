package domain

import (
	"math"
	"regexp"
	"strings"
	"time"
)

type LocalProductIdentity struct {
	Codprod string
	Descr   string
	EAN     *string
	RefForn *string
	Marca   *string
}

type IdentityCandidate struct {
	CatalogProductID string
	Title            string
	Attrs            []IdentityCandidateAttribute
}

type IdentityCandidateAttribute struct {
	ID        string
	Name      string
	ValueName *string
}

type identityCandidateEvaluation struct {
	index          int
	catalogID      string
	anchors        []string
	contradictions []string
	score          int
	status         MatchStatus
}

var identityMeasurePattern = regexp.MustCompile(`\b\d+(?:[.,]\d+)?\s*(?:ml|kg|cm|mm|pol|l|m|g|\")\b`)

func ResolveIdentity(product LocalProductIdentity, candidates []IdentityCandidate, decidedAt time.Time) (MatchDecision, error) {
	if len(candidates) == 0 {
		return NewMatchDecision(MatchDecisionInput{
			ProductID: product.Codprod, Status: MatchStatusNoCandidate,
			Anchors: []string{}, Contradictions: []string{},
			ConfidenceBand: ConfidenceBandBaixa, DecidedAt: decidedAt,
		})
	}

	evaluations := make([]identityCandidateEvaluation, 0, len(candidates))
	for index, candidate := range candidates {
		evaluations = append(evaluations, evaluateIdentityCandidate(product, candidate, index))
	}

	chosen := selectIdentityCandidate(evaluations, MatchStatusAccept)
	status := MatchStatusAccept
	if chosen == nil {
		chosen = selectIdentityCandidate(evaluations, MatchStatusReview)
		status = MatchStatusReview
	}
	if chosen == nil {
		chosen = selectIdentityCandidate(evaluations, MatchStatusReject)
		status = MatchStatusReject
	}
	// EAN-absent cap (IC-01): ACCEPT requires the "ean" anchor, which requires a present,
	// agreeing local EAN — so an absent EAN can never reach ACCEPT via the per-candidate gate.
	catalogID := chosen.catalogID
	return NewMatchDecision(MatchDecisionInput{
		ProductID: product.Codprod, CatalogProductID: &catalogID, Status: status,
		Anchors: chosen.anchors, Contradictions: chosen.contradictions,
		ConfidenceBand: identityConfidenceBand(chosen.score), DecidedAt: decidedAt,
	})
}

func evaluateIdentityCandidate(product LocalProductIdentity, candidate IdentityCandidate, index int) identityCandidateEvaluation {
	localEAN := normalizedOptional(product.EAN, true)
	localBrand := normalizedOptional(product.Marca, false)
	localModel := normalizedOptional(product.RefForn, false)
	candidateEAN := normalizedOptional(candidateAttribute(candidate.Attrs, "GTIN"), true)
	candidateBrand := normalizedOptional(candidateAttribute(candidate.Attrs, "BRAND"), false)
	candidateModel := normalizedOptional(candidateAttribute(candidate.Attrs, "MODEL"), false)

	anchors := make([]string, 0, 3)
	if agrees(localEAN, candidateEAN) {
		anchors = append(anchors, "ean")
	}
	if agrees(localBrand, candidateBrand) {
		anchors = append(anchors, "marca")
	}
	if agrees(localModel, candidateModel) {
		anchors = append(anchors, "refforn")
	}

	localText := normalizeIdentity(product.Descr)
	candidateText := normalizeIdentity(candidate.Title + " " + candidateAttributeText(candidate.Attrs))
	contradictions := identityContradictions(localText, candidateText, localBrand, candidateBrand)
	score := anchorScore(anchors) + int(math.Round(titleSimilarity(localText, normalizeIdentity(candidate.Title))*10))
	status := MatchStatusReview
	if len(contradictions) > 0 {
		status = MatchStatusReject
	} else if containsIdentityValue(anchors, "ean") && (containsIdentityValue(anchors, "marca") || containsIdentityValue(anchors, "refforn")) {
		status = MatchStatusAccept
	}
	return identityCandidateEvaluation{index, candidate.CatalogProductID, anchors, contradictions, score, status}
}

func selectIdentityCandidate(evaluations []identityCandidateEvaluation, status MatchStatus) *identityCandidateEvaluation {
	var chosen *identityCandidateEvaluation
	for index := range evaluations {
		candidate := &evaluations[index]
		if candidate.status != status {
			continue
		}
		// Forward iteration + strict `>` = first-occurrence (lowest index) wins ties.
		if chosen == nil || candidate.score > chosen.score {
			chosen = candidate
		}
	}
	return chosen
}

func candidateAttribute(attrs []IdentityCandidateAttribute, id string) *string {
	for _, attr := range attrs {
		if strings.EqualFold(strings.TrimSpace(attr.ID), id) {
			return attr.ValueName
		}
	}
	return nil
}

func candidateAttributeText(attrs []IdentityCandidateAttribute) string {
	var values []string
	for _, attr := range attrs {
		if attr.ValueName != nil {
			values = append(values, *attr.ValueName)
		}
	}
	return strings.Join(values, " ")
}

func normalizedOptional(value *string, digitsOnly bool) string {
	if value == nil {
		return ""
	}
	normalized := normalizeIdentity(*value)
	if !digitsOnly {
		return normalized
	}
	var digits strings.Builder
	for _, char := range normalized {
		if char >= '0' && char <= '9' {
			digits.WriteRune(char)
		}
	}
	return digits.String()
}

func normalizeIdentity(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "â", "a", "ã", "a", "ä", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i",
		"ó", "o", "ò", "o", "ô", "o", "õ", "o", "ö", "o",
		"ú", "u", "ù", "u", "û", "u", "ü", "u", "ç", "c",
	)
	return strings.Join(strings.Fields(replacer.Replace(value)), " ")
}

func agrees(left, right string) bool {
	return left != "" && right != "" && left == right
}

func anchorScore(anchors []string) int {
	score := 0
	for _, anchor := range anchors {
		switch anchor {
		case "ean":
			score += 50
		case "marca", "refforn":
			score += 20
		}
	}
	return score
}

func identityConfidenceBand(score int) ConfidenceBand {
	if score >= 85 {
		return ConfidenceBandAlta
	}
	if score >= 50 {
		return ConfidenceBandMedia
	}
	return ConfidenceBandBaixa
}

func identityContradictions(localText, candidateText, localBrand, candidateBrand string) []string {
	contradictions := make([]string, 0, 5)
	if divergentIdentityClass(identityKitValue(localText), identityKitValue(candidateText), false) {
		contradictions = append(contradictions, "kit")
	}
	if divergentIdentityClass(identityColorToken(localText), identityColorToken(candidateText), false) {
		contradictions = append(contradictions, "cor")
	}
	if divergentIdentityClass(identityMeasure(localText), identityMeasure(candidateText), false) {
		contradictions = append(contradictions, "medida")
	}
	if divergentIdentityClass(identityToken(localText, identityVoltages), identityToken(candidateText, identityVoltages), true) {
		contradictions = append(contradictions, "voltagem")
	}
	if localBrand != "" && candidateBrand != "" && localBrand != candidateBrand {
		contradictions = append(contradictions, "marca")
	}
	return contradictions
}

var identityVoltages = []string{"110v", "127v", "220v", "bivolt"}

func identityColorToken(text string) string {
	words := strings.FieldsFunc(text, func(char rune) bool { return !(char >= 'a' && char <= 'z') && !(char >= '0' && char <= '9') })
	for _, word := range words {
		switch word {
		case "preto", "preta":
			return "preto"
		case "branco", "branca":
			return "branco"
		case "vermelho", "vermelha":
			return "vermelho"
		case "amarelo", "amarela":
			return "amarelo"
		case "roxo", "roxa":
			return "roxo"
		case "azul", "verde", "cinza", "rosa", "marrom", "bege", "laranja":
			return word
		}
	}
	return ""
}

func identityKitValue(text string) string {
	if identityToken(text, []string{"kit", "combo", "conjunto"}) != "" {
		return "kit"
	}
	if identityToken(text, []string{"unit", "unidade", "unitario"}) != "" {
		return "unit"
	}
	return ""
}

func identityToken(text string, options []string) string {
	words := strings.FieldsFunc(text, func(char rune) bool { return !(char >= 'a' && char <= 'z') && !(char >= '0' && char <= '9') })
	for _, option := range options {
		for _, word := range words {
			if word == option {
				return option
			}
		}
	}
	return ""
}

func identityMeasure(text string) string {
	return strings.ReplaceAll(identityMeasurePattern.FindString(text), " ", "")
}

func divergentIdentityClass(left, right string, bivoltCompatible bool) bool {
	if left == "" || right == "" || left == right {
		return false
	}
	return !bivoltCompatible || (left != "bivolt" && right != "bivolt")
}

func containsIdentityValue(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func titleSimilarity(left, right string) float64 {
	leftTokens := uniqueIdentityTokens(left)
	rightTokens := uniqueIdentityTokens(right)
	if len(leftTokens) == 0 || len(rightTokens) == 0 {
		return 0
	}
	intersection := 0
	for _, leftToken := range leftTokens {
		if containsIdentityValue(rightTokens, leftToken) {
			intersection++
		}
	}
	union := len(leftTokens) + len(rightTokens) - intersection
	return float64(intersection) / float64(union)
}

func uniqueIdentityTokens(text string) []string {
	tokens := strings.FieldsFunc(text, func(char rune) bool { return !(char >= 'a' && char <= 'z') && !(char >= '0' && char <= '9') })
	unique := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if !containsIdentityValue(unique, token) {
			unique = append(unique, token)
		}
	}
	return unique
}
