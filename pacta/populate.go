package sqldb

func (a *AnalysisArtifact) Blobs() []*Blob {
	if a == nil {
		return nil
	}
	result := []*Blob{}
	result = append(result, a.Blob.Blobs()...)
	return result
}

func (a *Analysis) Blobs() []*Blob {
	if a == nil {
		return nil
	}
	result := []*Blob{}
	for _, aa := range a.Artifacts {
		result = append(result, aa.Blobs()...)
	}
	return result
}

func (b *Blob) Blobs() []*Blob {
	if b == nil {
		return nil
	}
	return []*Blob{b}
}

func (ii *IncompleteUpload) Blobs() []*Blob {
	if ii == nil {
		return nil
	}
	result := []*Blob{}
	result = append(result, ii.Blob.Blobs()...)
	return result
}

func (p *Portfolio) Blobs() []*Blob {
	if p == nil {
		return nil
	}
	result := []*Blob{}
	result = append(result, p.Blob.Blobs()...)
	return result
}
