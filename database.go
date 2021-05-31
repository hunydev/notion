package notion

import (
	"fmt"
	"strings"
)

type Database struct {
	JSON JSON
}

func (database *Database) Object() string {
	return "database"
}

func (database *Database) ID() string {
	return database.JSON.GetString("id")
}

func (database *Database) CreatedTime() string {
	return database.JSON.GetString("created_time")
}

func (database *Database) LastEditedTime() string {
	return database.JSON.GetString("last_edited_time")
}

func (database *Database) Title() []RichText {
	list := []RichText{}

	j, ok := database.JSON.GetJSONList("title")
	if !ok {
		return nil
	}

	for _, jj := range j {
		rt := RichText{}
		if jj.Unmarshal(&rt.JSON) == nil {
			list = append(list, rt)
		}
	}

	return list
}

func (database *Database) Properties() []Configuration {
	properties := []Configuration{}

	j, ok := database.JSON.GetJSON("properties")
	if !ok {
		return nil
	}

	for k, v := range j {
		jj := JSON{}
		if j.Marshal(v) != nil {
			continue
		}

		configuration, err := AssignConfiguration(k, jj)
		if err != nil {
			continue
		}

		properties = append(properties, configuration)
	}

	return nil
}

type Configuration interface {
	Name() string
	ID() string
	Type() string

	Json() JSON
	Interface() interface{}
}

func AssignConfiguration(name string, json JSON) (Configuration, error) {
	var configuration Configuration

	t := json.GetString("type")

	switch t {
	case TypePropertyTitle:
		configuration = &ConfigurationTitle{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyRichText:
		configuration = &ConfigurationText{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyNumber:
		configuration = &ConfigurationNumber{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertySelect:
		configuration = &ConfigurationSelect{&SelectOptionsConfiguration{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}}
	case TypePropertyMultiSelect:
		configuration = &ConfigurationMultiSelect{&SelectOptionsConfiguration{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}}
	case TypePropertyDate:
		configuration = &ConfigurationDate{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyFormula:
		configuration = &ConfigurationFormula{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyPeople:
		configuration = &ConfigurationPeople{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyFiles:
		configuration = &ConfigurationFile{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyCheckbox:
		configuration = &ConfigurationCheckbox{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyURL:
		configuration = &ConfigurationURL{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyEmail:
		configuration = &ConfigurationEmail{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyPhoneNumber:
		configuration = &ConfigurationPhoneNumber{&BaseConfiguration{&BaseProperty{name: name, JSON: json}}}
	default:
		return nil, fmt.Errorf("invalid type: '%s'", t)
	}

	return configuration, nil
}

type BaseConfiguration struct {
	*BaseProperty
}

func newBaseConfiguration(Name, ID, Type string, v interface{}) *BaseConfiguration {
	base := &BaseConfiguration{
		&BaseProperty{
			name: Name,
			JSON: JSON{
				"type": Type,
				Type:   v,
			},
		},
	}

	if len(strings.TrimSpace(ID)) > 0 {
		base.JSON["id"] = ID
	}
	return base
}

type ConfigurationTitle struct {
	*BaseConfiguration
}

type ConfigurationText struct {
	*BaseConfiguration
}

type ConfigurationNumber struct {
	*BaseConfiguration
}

func (configuration *ConfigurationNumber) Format() string {
	j, ok := configuration.JSON.GetJSON("number")
	if !ok {
		return ""
	}

	return j.GetString("format")
}

type SelectOptionsConfiguration struct {
	*BaseConfiguration
}

func (configuration *SelectOptionsConfiguration) Options() []SelectOption {
	j, ok := configuration.JSON.GetJSON(configuration.Type())
	if !ok {
		return nil
	}

	jj, ok := j.GetJSONList("options")
	if !ok {
		return nil
	}
	options := []SelectOption{}
	for _, jjj := range jj {
		option := SelectOption{}
		if jjj.Unmarshal(&option) == nil {
			options = append(options, option)
		}
	}

	return options
}

type ConfigurationSelect struct {
	*SelectOptionsConfiguration
}

type ConfigurationMultiSelect struct {
	*SelectOptionsConfiguration
}

type ConfigurationDate struct {
	*BaseConfiguration
}

type ConfigurationPeople struct {
	*BaseConfiguration
}

type ConfigurationFile struct {
	*BaseConfiguration
}

type ConfigurationCheckbox struct {
	*BaseConfiguration
}

type ConfigurationURL struct {
	*BaseConfiguration
}

type ConfigurationEmail struct {
	*BaseConfiguration
}

type ConfigurationPhoneNumber struct {
	*BaseConfiguration
}

type ConfigurationFormula struct {
	*BaseConfiguration
}

func (configuration *ConfigurationFormula) Expression() string {
	j, ok := configuration.JSON.GetJSON(configuration.Type())
	if !ok {
		return ""
	}

	return j.GetString("expression")
}
