package epss_test

import (
	"epssd/json"
	"fmt"
	"testing"
)

type TInfo struct {
	Type  int    `json:"type"`
	Title string `json:"title"`
	List  []int  `json:"list"`
}

func TestInsert(t *testing.T) {
	tt := `{
		"type": 212,
		"title": "通知dsww",
		"content": "所有人星期五都到会议室开会！",
		"job": {
			"id": 21123, "name": "email"
		},
		"list": [1231432213333321, "解决", 131.31, true, 1212]
	}`

	//fmt.Println(jj)
	var err error

	j, err := json.LoadString(tt)
	if err != nil {
		t.Error(err)
		return
	}
	title := j.GetBytes("job")
	s := string(title)
	fmt.Println(title, s)
	ty := j.GetInt64("type", 432)
	fmt.Println(ty)
	job := j.GetJson("job")
	if job != nil {
		id := job.GetInt("id")
		na := job.GetString("name")
		fmt.Print(id, na)
	}
	ls := j.GetJson("list")
	if ls != nil {
		for _, v := range ls.ToArray() {
			k := v.Kind()
			ii, _ := v.ToInt()
			fmt.Println(k, ii)
		}
		si := ls.GetSize()
		id := ls.GetIntAt(0)
		ff := ls.GetFloat32At(2)
		ss := ls.GetStringAt(1)
		fmt.Println(si, ff, id, ss)
	}
	//ty := GetInt(mm, "type", 123)
	//ty += 121
	return
}
