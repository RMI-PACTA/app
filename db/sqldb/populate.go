package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
)

type HasBlobs interface {
	Blobs() []*pacta.Blob
}

func (d *DB) PopulateBlobs(tx db.Tx, hbs []HasBlobs) error {
	blobs := []*pacta.Blob{}
	for _, hb := range hbs {
		blobs = append(blobs, hb.Blobs()...)
	}
	if len(blobs) == 0 {
		return nil
	}
	ids := []pacta.BlobID{}
	// TODO(grady) Dedupe ids in all multi-lookups using a standard utility
	for _, blob := range blobs {
		ids = append(ids, blob.ID)
	}
	actual, err := d.Blobs(tx, ids)
	if err != nil {
		return fmt.Errorf("looking up blobs: %w", err)
	}
	for _, blob := range blobs {
		if actualBlob, ok := actual[blob.ID]; ok {
			*blob = *actualBlob
		} else {
			return fmt.Errorf("blob %v not found, but requested", blob.ID)
		}
	}
	return nil
}

type HasAnalysisArtifacts interface {
	AnalysisArtifacts() []*pacta.AnalysisArtifact
}

func (d *DB) PopulateAnalysisArtifacts(tx db.Tx, haas []HasAnalysisArtifacts) error {
	aas := []*pacta.AnalysisArtifact{}
	for _, haa := range haas {
		aas = append(aas, haa.AnalysisArtifacts()...)
	}
	if len(aas) == 0 {
		return nil
	}
	ids := []pacta.AnalysisArtifactID{}
	for _, aa := range aas {
		ids = append(ids, aa.ID)
	}
	actual, err := d.AnalysisArtifacts(tx, ids)
	if err != nil {
		return fmt.Errorf("looking up aas: %w", err)
	}
	for _, aa := range aas {
		if actualAnalysisArtifact, ok := actual[aa.ID]; ok {
			*aa = *actualAnalysisArtifact
		} else {
			return fmt.Errorf("aa %v not found, but requested", aa.ID)
		}
	}
	return nil
}
