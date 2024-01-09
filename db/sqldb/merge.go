package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
)

func (d *DB) RecordUserMerge(tx db.Tx, fromUserID, toUserID, actorUserID pacta.UserID) error {
	err := d.exec(tx, `
		INSERT INTO user_merges 
			(from_user_id, to_user_id, actor_user_id)
			VALUES ($1, $2, $3);`, fromUserID, toUserID, actorUserID)
	if err != nil {
		return fmt.Errorf("inserting user merge: %w", err)
	}
	return nil
}

func (d *DB) RecordOwnerMerge(tx db.Tx, fromOwnerID, toOwnerID pacta.OwnerID, actorUserID pacta.UserID) error {
	err := d.exec(tx, `
		INSERT INTO owner_merges 
			(from_owner_id, to_owner_id, actor_user_id)
			VALUES ($1, $2, $3);`, fromOwnerID, toOwnerID, actorUserID)
	if err != nil {
		return fmt.Errorf("inserting owner merge: %w", err)
	}
	return nil
}

func (d *DB) expandAuditLogQueryToAccountForMerges(tx db.Tx, q *db.AuditLogQuery) (*db.AuditLogQuery, error) {
	var err error
	for _, w := range q.Wheres {
		w.InActorID, err = d.findAllMergedUsers(tx, w.InActorID)
		if err != nil {
			return nil, fmt.Errorf("finding merged users for actor_id: %w", err)
		}

		w.InTargetID, err = d.findAllMergedUsers(tx, w.InTargetID)
		if err != nil {
			return nil, fmt.Errorf("finding merged users for target_id: %w", err)
		}

		w.InActorOwnerID, err = d.findAllMergedOwners(tx, w.InActorOwnerID)
		if err != nil {
			return nil, fmt.Errorf("finding merged owners for actor_owner_id: %w", err)
		}

		w.InTargetOwnerID, err = d.findAllMergedOwners(tx, w.InTargetOwnerID)
		if err != nil {
			return nil, fmt.Errorf("finding merged owners for actor_owner_id: %w", err)
		}
	}
	return q, nil
}

func (d *DB) findAllMergedOwners(tx db.Tx, in []pacta.OwnerID) ([]pacta.OwnerID, error) {
	relationshipFn := func(id pacta.OwnerID) ([]pacta.OwnerID, error) {
		return d.findMergedOwners(tx, pacta.OwnerID(id))
	}
	return recursivelyExpandRelationships(in, relationshipFn)
}

func (d *DB) findMergedOwners(tx db.Tx, id pacta.OwnerID) ([]pacta.OwnerID, error) {
	rows, err := d.query(tx, `
		(SELECT from_owner_id FROM owner_merges WHERE to_owner_id = $1)
		UNION
		(SELECT to_owner_id FROM owner_merges WHERE from_owner_id = $1);`, id)
	if err != nil {
		return nil, fmt.Errorf("querying owner_merges: %w", err)
	}
	ownerIDs, err := mapRowsToIDs[pacta.OwnerID]("merged_owners", rows)
	if err != nil {
		return nil, fmt.Errorf("mapping rows to owner ids: %w", err)
	}
	return ownerIDs, nil
}

func (d *DB) findAllMergedUsers(tx db.Tx, in []string) ([]string, error) {
	relationshipFn := func(id string) ([]string, error) {
		others, err := d.findMergedUsers(tx, pacta.UserID(id))
		if err != nil {
			return nil, fmt.Errorf("finding merged users for %q: %w", id, err)
		}
		return asStrs(others), nil
	}
	return recursivelyExpandRelationships(in, relationshipFn)
}

func (d *DB) findMergedUsers(tx db.Tx, id pacta.UserID) ([]pacta.UserID, error) {
	rows, err := d.query(tx, `
		(SELECT from_user_id FROM user_merges WHERE to_user_id = $1)
		UNION
		(SELECT to_user_id FROM user_merges WHERE from_user_id = $1);`, id)
	if err != nil {
		return nil, fmt.Errorf("querying user_merges: %w", err)
	}
	userIDs, err := mapRowsToIDs[pacta.UserID]("merged_users", rows)
	if err != nil {
		return nil, fmt.Errorf("mapping rows to user ids: %w", err)
	}
	return userIDs, nil
}

func recursivelyExpandRelationships[S ~string](in []S, relatedFn func(S) ([]S, error)) ([]S, error) {
	if len(in) == 0 {
		return in, nil
	}
	all := asSet(in)
	lookedUp := map[S]bool{}
	for len(lookedUp) < len(all) {
		for s := range all {
			if lookedUp[s] {
				continue
			}
			related, err := relatedFn(s)
			if err != nil {
				return nil, fmt.Errorf("finding relationships for %q: %w", s, err)
			}
			for _, r := range related {
				all[r] = true
			}
			lookedUp[s] = true
		}
	}
	return keys(all), nil
}
