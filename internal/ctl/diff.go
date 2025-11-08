package ctl

import "github.com/etoneja/go-keeper/internal/ctl/types"

func diffSecrets(local []*types.LocalSecret, remote []*types.RemoteSecret) *types.SecretsDiff {
	localMap := make(map[string]*types.LocalSecret)
	remoteMap := make(map[string]*types.RemoteSecret)

	for _, secret := range local {
		localMap[secret.UUID] = secret
	}
	for _, secret := range remote {
		remoteMap[secret.UUID] = secret
	}

	diff := &types.SecretsDiff{}

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
