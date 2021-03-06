package apps

import (
	"github.com/jenkins-x/jx-apps/pkg/jxapps"
	"github.com/jenkins-x/jx-promote/pkg/rules"
	"github.com/pkg/errors"
)

// AppsRule uses a jx-apps.yml file
func AppsRule(r *rules.PromoteRule) error {
	config := r.Config
	if config.Spec.AppsRule == nil {
		return errors.Errorf("no appsRule configured")
	}
	rule := config.Spec.AppsRule
	err := modifyAppsFile(r, r.Dir, rule.Namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to modify chart files in dir %s", r.Dir)
	}
	return nil
}

// ModifyAppsFile modifies the 'jx-apps.yml' file to add/update/remove apps
func modifyAppsFile(r *rules.PromoteRule, dir string, promoteNS string) error {
	appsConfig, fileName, err := jxapps.LoadAppConfig(dir)
	if fileName == "" {
		// if we don't have a `jx-apps.yml` then just return immediately
		return nil
	}
	if err != nil {
		return err
	}
	err = modifyApps(r, appsConfig, promoteNS)
	if err != nil {
		return err
	}

	err = appsConfig.SaveConfig(fileName)
	if err != nil {
		return err
	}
	return nil
}

func modifyApps(r *rules.PromoteRule, appsConfig *jxapps.AppConfig, promoteNS string) error {
	if r.DevEnvContext == nil {
		return errors.Errorf("no devEnvContext")
	}
	app := r.AppName
	version := r.Version
	details, err := r.DevEnvContext.ChartDetails(app, r.HelmRepositoryURL)
	if err != nil {
		return errors.Wrapf(err, "failed to get chart details for %s repo %s", app, r.HelmRepositoryURL)
	}
	details.DefaultPrefix(appsConfig, "dev")

	for i := range appsConfig.Apps {
		appConfig := &appsConfig.Apps[i]
		if appConfig.Name == app || appConfig.Name == details.Name {
			appConfig.Version = version
			if promoteNS != "" {
				appConfig.Namespace = promoteNS
			}
			return nil
		}
	}
	appsConfig.Apps = append(appsConfig.Apps, jxapps.App{
		Name:      details.Name,
		Version:   version,
		Namespace: promoteNS,
	})
	return nil
}
