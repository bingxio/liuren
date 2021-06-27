package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	port = ":8080"
)

var (
	earthBranch = []string{
		"子",
		"丑",
		"寅",
		"卯",
		"辰",
		"巳",
		"午",
		"未",
		"申",
		"酉",
		"戌",
		"亥",
	}
	heavenStem = []string{
		"甲",
		"乙",
		"丙",
		"丁",
		"戊",
		"己",
		"庚",
		"辛",
		"壬",
		"癸",
	}
	twelveGod = []string{
		"贵",
		"蛇",
		"朱",
		"合",
		"勾",
		"青",
		"空",
		"白",
		"常",
		"玄",
		"阴",
		"后",
	}
)

type Response struct {
	Data interface{} `json:"data"`
}

type Data struct {
	X string `json:"x"`
	Y string `json:"y"`
	Z string `json:"z"`
}

func main() {
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		r.ParseForm()
		x := r.Form.Get("x") // 日辰
		y := r.Form.Get("y") // 时辰
		z := r.Form.Get("z") // 月将

		t1, t2 := HeavenEarthPlate(y, z) // 天地盘
		log.Println(t1)
		log.Println(t2)

		t3 := TwelveGods(x, y, t2) // 十二神
		log.Println(t3)

		t4 := ShieldTiangan(x, t2) // 盾天干
		log.Println(t4)

		t5 := LessonFour(x, t2, t3, t4) // 布四课
		log.Println(t5)

		b, _ := json.Marshal(Response{
			Data: Data{
				X: BuildA(t1, t2, t3, t4),
				Y: BuildB(t5),
			},
		})
		w.Write(b)
	})
	fmt.Println("Listen and serve http on port", port)
	_ = http.ListenAndServe(port, nil)
}

// 天地盘
func HeavenEarthPlate(a, b string) ([]string, []string) {
	log.Println("时辰：", a, "月将：", b)

	heavenPlate := make([]string, 12)
	index := IndexOf(earthBranch, b)

	for k, v := range earthBranch {
		if v == a {
			for i := 0; i < 12; i++ {
				heavenPlate[k] = earthBranch[index]
				index++
				k++

				if index == 12 {
					index = 0
				}
				if k == 12 {
					k = 0
				}
			}
		}
	}
	return earthBranch, heavenPlate
}

// 十二神
func TwelveGods(a, b string, t2 []string) []string {
	x, y, z := 0, 0, 0
	switch string([]rune(a)[:1]) {
	case "甲", "戊", "庚":
		x = 1 // 牛
		y = 7 // 羊
	case "乙", "己":
		x = 0 // 鼠
		y = 8 // 猴
	case "丙", "丁":
		x = 11 // 猪
		y = 9  // 鸡
	case "壬", "癸":
		x = 5 // 蛇
		y = 3 // 兔
	case "辛":
		x = 6 // 马
		y = 2 // 虎
	}
	log.Println("日辰：", a, "时辰：", b, "贵人：", earthBranch[x], earthBranch[y])
	switch string([]rune(b)[:1]) {
	case "卯", "辰", "巳", "午", "未", "申":
		z = x
		log.Println("昼贵：", earthBranch[z])
	case "酉", "戌", "亥", "子", "丑", "寅":
		z = y
		log.Println("夜贵：", earthBranch[z])
	}
	d := IndexOf(t2, earthBranch[z])
	m := make([]string, 12)
	switch d {
	case 11, 0, 1, 2, 3, 4:
		// 顺治
		for i := 0; i < 12; i++ {
			m[d] = twelveGod[i]
			d++
			if d == 12 {
				d = 0
			}
		}
	default:
		// 逆治
		for i := 0; i < 12; i++ {
			m[d] = twelveGod[i]
			d--
			if d == -1 {
				d = 11
			}
		}
	}
	return m
}

// 盾天干
func ShieldTiangan(a string, t2 []string) []string {
	x := string([]rune(a)[:1])
	y := string([]rune(a)[1:2])
	i := IndexOf(heavenStem, x)
	j := IndexOf(earthBranch, y)
	k := j - i
	if k < 0 {
		k = 12 + k
	}
	i = IndexOf(t2, earthBranch[k])
	m := make([]string, 12)
	for p := 0; p < 12; p++ {
		if p == 10 {
			m[i] = heavenStem[0]
		} else if p == 11 {
			m[i] = heavenStem[1]
		} else {
			m[i] = heavenStem[p]
		}
		i++
		if i == 12 {
			i = 0
		}
	}
	return m
}

// 布四课
func LessonFour(a string, t2, t3, t4 []string) [4][4]string {
	m := [4][4]string{}
	m[1][3] = string([]rune(a)[:1])  // 日干
	m[1][1] = string([]rune(a)[1:2]) // 日支
	k := 0
	switch m[1][3] {
	case "甲":
		k = 2
	case "乙":
		k = 4
	case "丙", "戊":
		k = 5
	case "丁", "己":
		k = 6
	case "庚":
		k = 8
	case "辛":
		k = 10
	case "壬":
		k = 11
	case "癸":
		k = 1
	}
	m[0][3] = t2[IndexOf(earthBranch, earthBranch[k])] // 课一上
	m[1][2] = m[0][3]                                  // 课二下
	m[0][2] = t2[IndexOf(earthBranch, m[1][2])]        // 课二上
	m[0][1] = t2[IndexOf(earthBranch, m[1][1])]        // 课三上
	m[1][0] = m[0][1]                                  // 课四下
	m[0][0] = t2[IndexOf(earthBranch, m[1][0])]        // 课四上
	m[2][3] = t3[IndexOf(t2, m[0][3])]                 // 课一
	m[2][2] = t3[IndexOf(t2, m[0][2])]                 // 课二
	m[2][1] = t3[IndexOf(t2, m[0][1])]                 // 课三
	m[2][0] = t3[IndexOf(t2, m[0][0])]                 // 课四
	m[3][3] = t4[IndexOf(t2, m[0][3])]                 // 课一
	m[3][2] = t4[IndexOf(t2, m[0][2])]                 // 课二
	m[3][1] = t4[IndexOf(t2, m[0][1])]                 // 课三
	m[3][0] = t4[IndexOf(t2, m[0][0])]                 // 课四
	return m
}

func IndexOf(v []string, t string) int {
	for k, v := range v {
		if v == t {
			return k
		}
	}
	return -1
}

func GetColor(v string) string {
	switch v {
	case "甲", "乙", "寅", "卯", "合", "青":
		return "green"
	case "丙", "丁", "巳", "午", "朱", "蛇":
		return "red"
	case "戊", "己", "丑", "辰", "未", "戌", "贵", "勾", "空", "常":
		return "orange"
	case "庚", "辛", "申", "酉", "白", "阴":
		return "blue"
	case "壬", "癸", "子", "亥", "玄", "后":
		return "grey"
	default:
		return ""
	}
}

func BuildA(t1, t2, t3, t4 []string) string {
	m := strings.Builder{}
	m.WriteString(fmt.Sprintf(`<table>
	<!-- 1 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>`, GetColor(t4[5]), t4[5], GetColor(t4[6]), t4[6], GetColor(t4[7]), t4[7], GetColor(t4[8]), t4[8]))
	m.WriteString(fmt.Sprintf(`
	<!-- 2 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>`, GetColor(t3[5]), t3[5], GetColor(t3[6]), t3[6], GetColor(t3[7]), t3[7], GetColor(t3[8]), t3[8]))
	m.WriteString(fmt.Sprintf(`
	<!-- 3 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>`, GetColor(t2[5]), t2[5], GetColor(t2[6]), t2[6], GetColor(t2[7]), t2[7], GetColor(t2[8]), t2[8]))
	m.WriteString(fmt.Sprintf(`
	<!-- 4 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>`, GetColor(t1[5]), t1[5], GetColor(t1[6]), t1[6], GetColor(t1[7]), t1[7], GetColor(t1[8]), t1[8]))
	m.WriteString(fmt.Sprintf(`
	<!-- 5 -->
	<tr>
		<td style="color: %s">%s</td>
    	<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
	</tr>`,
		GetColor(t4[4]), t4[4], GetColor(t3[4]), t3[4], GetColor(t2[4]), t2[4], GetColor(t1[4]), t1[4],
		GetColor(t1[9]), t1[9], GetColor(t2[9]), t2[9], GetColor(t3[9]), t3[9], GetColor(t4[9]), t4[9]))
	m.WriteString(fmt.Sprintf(`
	<!-- 6 -->
	<tr>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
	</tr>`,
		GetColor(t4[3]), t4[3], GetColor(t3[3]), t3[3], GetColor(t2[3]), t2[3], GetColor(t1[3]), t1[3],
		GetColor(t1[10]), t1[10], GetColor(t2[10]), t2[10], GetColor(t3[10]), t3[10], GetColor(t4[10]), t4[10]))
	m.WriteString(fmt.Sprintf(`
	<!-- 7 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>`, GetColor(t1[2]), t1[2], GetColor(t1[1]), t1[1], GetColor(t1[0]), t1[0], GetColor(t1[11]), t1[11]))
	m.WriteString(fmt.Sprintf(`
	<!-- 8 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>`, GetColor(t2[2]), t2[2], GetColor(t2[1]), t2[1], GetColor(t2[0]), t2[0], GetColor(t2[11]), t2[11]))
	m.WriteString(fmt.Sprintf(`
	<!-- 9 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>`, GetColor(t3[2]), t3[2], GetColor(t3[1]), t3[1], GetColor(t3[0]), t3[0], GetColor(t3[11]), t3[11]))
	m.WriteString(fmt.Sprintf(`
	<!-- 10 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>
</table>`, GetColor(t4[2]), t4[2], GetColor(t4[1]), t4[1], GetColor(t4[0]), t4[0], GetColor(t4[11]), t4[11]))
	return m.String()
}

func BuildB(t5 [4][4]string) string {
	m := strings.Builder{}
	m.WriteString(fmt.Sprintf(`<table>
	<!-- 1 -->
	<tr><td></td></tr>
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
	</tr>`, GetColor(t5[3][0]), t5[3][0], GetColor(t5[3][1]),
		t5[3][1], GetColor(t5[3][2]), t5[3][2], GetColor(t5[3][3]), t5[3][3]))
	m.WriteString(fmt.Sprintf(`
	<!-- 2 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
	</tr>`, GetColor(t5[2][0]), t5[2][0], GetColor(t5[2][1]),
		t5[2][1], GetColor(t5[2][2]), t5[2][2], GetColor(t5[2][3]), t5[2][3]))
	m.WriteString(fmt.Sprintf(`
	<!-- 3 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
	</tr>`, GetColor(t5[0][0]), t5[0][0], GetColor(t5[0][1]),
		t5[0][1], GetColor(t5[0][2]), t5[0][2], GetColor(t5[0][3]), t5[0][3]))
	m.WriteString(fmt.Sprintf(`
	<!-- 4 -->
	<tr>
		<td></td>
		<td></td>
		<td></td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
		<td style="color: %s">%s</td>
	</tr>
</table>`, GetColor(t5[1][0]), t5[1][0], GetColor(t5[1][1]),
		t5[1][1], GetColor(t5[1][2]), t5[1][2], GetColor(t5[1][3]), t5[1][3]))
	return m.String()
}
