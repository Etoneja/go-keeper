package ctl

import "github.com/etoneja/go-keeper/internal/ctl/types"

func diffSecrets(local []*types.Secret, remote []types.Secreter) *types.SecretDiff {
	localMap := make(map[string]*types.Secret)
	remoteMap := make(map[string]types.Secreter)

	for _, secret := range local {
		localMap[secret.UUID] = secret
	}
	for _, secret := range remote {
		remoteMap[secret.GetUUID()] = secret
	}

	diff := &types.SecretDiff{}

	for id, secret := range localMap {
		if _, exists := remoteMap[id]; !exists {
			diff.LocalOnly = append(diff.LocalOnly, secret)
		}
	}

	for id, secret := range remoteMap {
		if _, exists := localMap[id]; !exists {
			diff.RemoteOnly = append(diff.RemoteOnly, secret)
		}
	}

	for id, localSecret := range localMap {
		if remoteSecret, exists := remoteMap[id]; exists {
			diff.Both = append(diff.Both, &types.CheckPair{
				Local:  localSecret,
				Remote: remoteSecret,
			})
		}
	}

	return diff
}
