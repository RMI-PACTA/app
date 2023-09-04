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

type HasOwner interface {
	Owners() []*pacta.Owner
}

func (d *DB) PopulateOwners(tx db.Tx, hos []HasOwner) error {
	os := []*pacta.Owner{}
	for _, ho := range hos {
		os = append(os, ho.Owners()...)
	}
	if len(os) == 0 {
		return nil
	}
	ids := []pacta.OwnerID{}
	for _, o := range os {
		ids = append(ids, o.ID)
	}
	actual, err := d.Owners(tx, ids)
	if err != nil {
		return fmt.Errorf("looking up owners: %w", err)
	}
	for _, o := range os {
		if actualOwner, ok := actual[o.ID]; ok {
			*o = *actualOwner
		} else {
			return fmt.Errorf("owner %v not found, but requested", o.ID)
		}
	}
	return nil
}

type HasUser interface {
	Users() []*pacta.User
}

func (d *DB) PopulateUsers(tx db.Tx, hus []HasUser) error {
	us := []*pacta.User{}
	for _, hu := range hus {
		us = append(us, hu.Users()...)
	}
	if len(us) == 0 {
		return nil
	}
	ids := []pacta.UserID{}
	for _, u := range us {
		ids = append(ids, u.ID)
	}
	actual, err := d.Users(tx, ids)
	if err != nil {
		return fmt.Errorf("looking up users: %w", err)
	}
	for _, u := range us {
		if actualUser, ok := actual[u.ID]; ok {
			*u = *actualUser
		} else {
			return fmt.Errorf("user %v not found, but requested", u.ID)
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
