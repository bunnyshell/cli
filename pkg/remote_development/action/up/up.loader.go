package up

import (
	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/remote_development/config"
	"bunnyshell.com/sdk"
)

func (up *Options) getProfile() (*sdk.ProfileItem, error) {
	if up.manager.HasProfileName() {
		return up.loadLocalProfile()
	}

	profileOrName, err := up.getProfileFromComponentConfig()
	if err != nil {
		return nil, err
	}

	if profileOrName.ProfileItem != nil {
		return profileOrName.ProfileItem, nil
	}

	if profileOrName.String != nil {
		up.manager.SetProfileName(*profileOrName.String)

		return up.loadLocalProfile()
	}

	return nil, ErrUnknownProfileType
}

func (up *Options) loadLocalProfile() (*sdk.ProfileItem, error) {
	if err := up.manager.Load(); err != nil {
		return nil, err
	}

	profile, err := up.manager.GetProfile()
	if err != nil {
		return nil, err
	}

	return up.convertProfile(profile)
}

func (up *Options) convertProfile(profile *config.Profile) (*sdk.ProfileItem, error) {
	options := component.NewRDevContextOptions(up.resourceLoader.Component.GetId())
	options.Profile = profile

	profileConfiguration, err := component.RDevContext(options)
	if err != nil {
		return nil, err
	}

	return (*sdk.ProfileItem)(profileConfiguration), nil
}

func (up *Options) selectContainerName(containers map[string]sdk.ContainerConfigItem) (*string, error) {
	if len(containers) == 0 {
		return nil, ErrEmptyList
	}

	if len(containers) == 1 && !up.ManualSelectSingleResource {
		for name := range containers {
			return &name, nil
		}
	}

	_, name, err := interactive.Choose("Choose a remote development profile", containersToProfiles(containers))
	if err != nil {
		return nil, err
	}

	return &name, nil
}

func containersToProfiles(containers map[string]sdk.ContainerConfigItem) []string {
	profiles := make([]string, 0)

	for name := range containers {
		profiles = append(profiles, name)
	}

	return profiles
}
