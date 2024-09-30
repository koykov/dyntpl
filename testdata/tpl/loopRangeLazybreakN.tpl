{
	"id":"{%= user.Id %}",
	"name":"{%= user.Name %}",
	"fin_history":[
    {% for k, item := range user.Finance.History sep , %}
    {%= k %}:{
      "utime":{%= item.DateUnix %},
      "cost":{%= item.Cost %},
      "desc":"{%= item.Comment %}"
      {% if k == 2 %}{% lazybreak 2 %}{% endif %}
    }
    {% endfor %}
	]
}
