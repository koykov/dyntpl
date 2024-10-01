{
	"id":"{%= user.Id %}",
	"name":"{%= user.Name %}",
	"history":[
    {% for k, item := range user.HistoryTree sep , %}
    {%= k %}:{
      "utime":{%= item.DateUnix %},
      "cost":{%= item.Cost %},
      "desc":"{%= item.Comment %}"
    }
    {% else %}
      -1:{"error":"no history rows"}
    {% endfor %}
	]
}
