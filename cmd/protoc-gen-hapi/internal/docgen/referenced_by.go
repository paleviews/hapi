package docgen

type referencedBy struct {
	schemaID string
	by       *referencedBy
}

func (rb *referencedBy) findCycle(schemaID string) []string {
	var contains bool
	for cur := rb; cur != nil; cur = cur.by {
		if schemaID == cur.schemaID {
			contains = true
			break
		}
	}
	if !contains {
		return nil
	}
	var cycle []string
	for cur := rb; cur != nil; cur = cur.by {
		cycle = append(cycle, cur.schemaID)
		if schemaID == cur.schemaID {
			break
		}
	}
	return cycle
}

func (rb *referencedBy) add(schemaID string) *referencedBy {
	return &referencedBy{
		schemaID: schemaID,
		by:       rb,
	}
}
