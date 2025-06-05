package app

type setsRelation int

const (
	setsRelationSame setsRelation = iota + 1
	setsRelationLeftIncludesRight
	setsRelationRightIncludesLeft
	setsRelationIntersection
)

func parseRelationshipBetweenSets(a, b map[int]bool) setsRelation {
	belongTo := func(s1, s2 map[int]bool) bool {
		for k := range s1 {
			if !s2[k] {
				return false
			}
		}

		return true
	}

	if len(a) != len(b) {
		if len(a) > len(b) {
			if belongTo(b, a) {
				return setsRelationLeftIncludesRight
			}
		} else {
			if belongTo(a, b) {
				return setsRelationRightIncludesLeft
			}
		}
	} else {
		if belongTo(a, b) {
			return setsRelationSame
		}
	}

	return setsRelationIntersection
}

func hasIntersection(a, b map[int]bool) bool {
	for k := range a {
		if b[k] {
			return true
		}
	}

	return false
}

func findRelationsBetweenCategories(newSets, oldSets []map[int]bool) (newToOld, oldToNew map[int][]int) {
	newToOld = make(map[int][]int)
	oldToNew = make(map[int][]int)

	for i, newSet := range newSets {
		for j, oldSet := range oldSets {
			if hasIntersection(newSet, oldSet) {
				newToOld[i] = append(newToOld[i], j)
				oldToNew[j] = append(oldToNew[j], i)
			}
		}
	}

	return
}
