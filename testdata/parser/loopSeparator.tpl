
{
	"rules": [
		{% for _, rule := range config.Rules sep , %}
		{
			"key": {%= rule.Id %},
			"val": {%= rule.Slug %}
		}
		{% endfor %}
	]
}