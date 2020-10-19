package gitlab

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func (c *Client) generateProject(prefix, course, assignment, assignmentPath string,
	inID int) (*gitlab.Project, bool, error) {
	generated := false
	name := assignment + "-" + prefix
	description := "generated by glabs"

	if desc := viper.GetString(course + "." + assignment + ".description"); desc != "" {
		description = desc
	}

	log.Debug().Str("desciption", description).Msg("generating with description")

	containerRegistryEnabled := false

	if viper.GetBool(course + "." + assignment + ".containerRegistry") {
		containerRegistryEnabled = true
	}

	p := &gitlab.CreateProjectOptions{
		Name:                     gitlab.String(name),
		Description:              gitlab.String(description),
		NamespaceID:              gitlab.Int(inID),
		MergeRequestsAccessLevel: gitlab.AccessControl("enabled"),
		IssuesAccessLevel:        gitlab.AccessControl("enabled"),
		BuildsAccessLevel:        gitlab.AccessControl("enabled"),
		JobsEnabled:              gitlab.Bool(true),
		Visibility:               gitlab.Visibility(gitlab.PrivateVisibility),
		ContainerRegistryEnabled: gitlab.Bool(containerRegistryEnabled),
	}

	project, _, err := c.Projects.CreateProject(p)

	if err == nil {
		log.Debug().Str("name", name).Msg("generated repo")
		generated = true
	} else {
		if project == nil {
			projectname := assignmentPath + "/" + name
			log.Debug().Err(err).Str("name", projectname).Msg("searching for project")
			project, err = c.findProject(projectname)
			if err != nil {
				log.Fatal().Err(err)
				return nil, false, fmt.Errorf("%w", err)
			}
		} else {
			log.Fatal().Err(err)
		}
	}

	return project, generated, nil
}

func (c *Client) findProject(projectname string) (*gitlab.Project, error) {
	opt := &gitlab.ListProjectsOptions{
		Search:           gitlab.String(projectname),
		SearchNamespaces: gitlab.Bool(true),
	}
	projects, _, err := c.Projects.ListProjects(opt)
	if err != nil {
		log.Error().Err(err).
			Str("projectname", projectname).
			Msg("no project found")
	} else {
		switch len(projects) {
		case 1:
			return projects[0], nil
		case 0:
			log.Debug().Interface("projects", projects).Msg("more than one project found")
			return nil, errors.New("more than one project found")
		default:
			log.Debug().Msg("more than one project matching the search string found")
			for _, project := range projects {
				if project.PathWithNamespace == projectname {
					log.Debug().Str("name", projectname).Msg("found project")
					return project, nil
				}
			}
			log.Debug().Str("name", projectname).Msg("project not found")
			return nil, errors.New("project not found")
		}
	}
	return nil, nil // could not happen
}
